package mpevents

import (
	"bytes"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/http/httputil"
	"golang.org/x/oauth2/google"
)

type PubSubMsg struct {
	EventId     	string	`json:"eventId"`
	EventType   	string	`json:"eventType"`
	Entitlement		EntitlementMeta `json:entitlement,omitempty`
	Account			AccountMeta `json:account,omitempty`
}

type EntitlementMeta struct {
	Id     			string	`json:"id,omitempty"`
	UpdateTime    	string	`json:"updateTime,omitempty"`
}

type AccountMeta struct {
	Id  			string     	`json:"id"`
	UpdateTime   	string    	`json:"updateTime"`
}

type Account struct {
	Name  			string     	`json:"name"`
	UpdateTime   	string    	`json:"updateTime,omitempty"`
	CreateTime      string    	`json:"createTime,omitempty"`
	Provider     	string		`json:"provider,omitempty"`
	State 	 		string      `json:"state,omitempty"`
	Approvals    	[]Approval     `json:"approvals,omitempty"`
}

type Approval struct {
	Name  			string     	`json:"name"`
	State  			string     	`json:"state"`
	Reason  		string     	`json:"reason"`
	UpdateTime  	string     	`json:"updateTime"`
}

type Entitlement struct {
	Name     			string	`json:"name"`
	Account   			string	`json:"account"`
	Provider    		string	`json:"provider"`
	Product  			string	`json:"product"`
	Plan     	  		string	`json:"plan"`
	NewPendingPlan 	  	string	`json:"newPendingPlan"`
	State    	  		string	`json:"state"`
	UpdateTime    	  	string	`json:"updateTime"`
	CreateTime    	  	string	`json:"createTime"`
	UsageReportingId    string	`json:"usageReportingId"`
	MessageToUser    	string	`json:"messageToUser"`
}

type PubSubListener struct {
	PubSubSubscription    			string
	PubSubTopicPrefix				string
	SubscriptionServiceUrl 			string
	CloudCommerceProcurementUrl    	string
	PartnerId    					string
	GcpProjectId                    string
}

var (
	cloudCommerceProcurementBaseUrl string
	subscriptionServiceBaseUrl string
)

func GetPubSubListener(pubSubSubscription string, pubSubTopicPrefix string, subscriptionServiceUrl string, cloudCommerceProcurementUrl string, partnerId string, gcpProjectId string) *PubSubListener {
	cloudCommerceProcurementBaseUrl = cloudCommerceProcurementUrl + "/providers/" + partnerId
	subscriptionServiceBaseUrl = subscriptionServiceUrl
	return &PubSubListener{
		pubSubSubscription,
		pubSubTopicPrefix,
		subscriptionServiceUrl,
		cloudCommerceProcurementUrl,
		partnerId,
		gcpProjectId,
	}
}

func (lstnr *PubSubListener) Listen() error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, lstnr.GcpProjectId)
	if err != nil {
		log.Printf("Error creating pubsub client %s: %#v", lstnr.PubSubSubscription, err)
		return err
	}

	topicId := lstnr.PubSubTopicPrefix+lstnr.GcpProjectId
	topic := client.Topic(topicId)
	if exists, errTp := topic.Exists(ctx); !exists && errTp == nil {
		if _, err := client.CreateTopic(ctx, topicId); err != nil {
			log.Printf("Failed to create topic: %#v", err)
			return err
		}
	} else if errTp != nil {
		log.Printf("Error checking for topic: %#v", errTp)
		return errTp
	}

	subscription := client.Subscription(lstnr.PubSubSubscription)

	if exists, errSub := subscription.Exists(ctx); !exists && errSub == nil {
		if _, err = client.CreateSubscription(ctx, lstnr.PubSubSubscription, pubsub.SubscriptionConfig{Topic: topic}); err != nil {
			log.Printf("Failed to create subscription: %#v", err)
			return err
		}
	} else if errSub != nil{
		log.Printf("Error checking for subscription: %#v", errSub)
		return errSub
	}

	log.Printf("Begin receiving messages from %s \n", lstnr.PubSubSubscription)
	errRcv := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		pubSubMsg := PubSubMsg{}
		if err := json.Unmarshal(msg.Data, &pubSubMsg); err != nil {
			log.Printf("could not decode message data: %#v \n", msg)
			msg.Nack()
			return
		}

		log.Printf("Received msg %#v", pubSubMsg)
		if processPubSubMsg(pubSubMsg) {
			msg.Ack()
			log.Printf("Message %s acked.",pubSubMsg.EventId)
		} else {
			msg.Nack()
			log.Printf("Message %s nacked.",pubSubMsg.EventId)
		}
	})
	if errRcv != nil {
		return errRcv
	}
	return nil
}

func processPubSubMsg(pubSubMsg PubSubMsg) bool {
	switch pubSubMsg.EventType {
	case "ACCOUNT_ACTIVE":
		if err := syncAccount(pubSubMsg.Account.Id); err != nil {
			log.Printf("Unable to update account %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		} else {

		}
		//check to see if we have any entitlements we need to approve
	case "ENTITLEMENT_CREATION_REQUESTED":
		rcvdEntitlement := Entitlement{}
		if err := getEntitlement(pubSubMsg.Entitlement.Id, &rcvdEntitlement); err != nil {
			if accountExists, acctErr := accountExists(rcvdEntitlement.Account); accountExists && acctErr == nil {
				if postErr := postEntitlementApproval(pubSubMsg.Entitlement.Id,false); postErr != nil {
					log.Printf("Unable to approve entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, postErr)
					return false
				}
			} else if !accountExists && acctErr == nil {
				if err := updateEntitlement(&rcvdEntitlement); err != nil {
					log.Printf("Unable to update entitlement %#v due to error %#v \n", rcvdEntitlement, err)
				}
			} else {
				log.Printf("Unable to determine if account %#v exists due to error %#v \n", rcvdEntitlement.Account, acctErr)
				return false
			}
		}
	case "ENTITLEMENT_ACTIVE":
		err := syncEntitlement(pubSubMsg.Entitlement.Id)
		if err != nil {
			log.Printf("Unable to update entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ENTITLEMENT_PLAN_CHANGE_REQUESTED":
		err := postEntitlementApproval(pubSubMsg.Entitlement.Id, true)
		if err == nil {
			err := syncEntitlement(pubSubMsg.Entitlement.Id)
			if err != nil {
				log.Printf("Unable to update entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
				return false
			}
		} else {
			log.Printf("Unable to approve entitlement plan change %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ENTITLEMENT_PLAN_CHANGED":
		err := syncEntitlement(pubSubMsg.Entitlement.Id)
		if err != nil {
			log.Printf("Unable to update entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ENTITLEMENT_PENDING_CANCELLATION":
		log.Println("ENTITLEMENT_PENDING_CANCELLATION ignored.")
	case "ENTITLEMENT_CANCELLED":
		if err := syncEntitlement(pubSubMsg.Entitlement.Id); err != nil {
			log.Printf("Unable to update entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ENTITLEMENT_DELETED":
		if err := deleteEntitlement(pubSubMsg.Entitlement.Id); err != nil {
			log.Printf("Unable to delete entitlement %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ACCOUNT_DELETED":
		if err := deleteAccount(pubSubMsg.Account.Id); err != nil {
			log.Printf("Unable to delete account %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	default:
		log.Printf("Unknown pubsub event type %#v \n", pubSubMsg)
	}
	return true
}

func postEntitlementApproval(entitlementId string, planChange bool) error {
	procurementUrl := cloudCommerceProcurementBaseUrl + "/entitlements/" + entitlementId + ":approve"
	if planChange {
		procurementUrl = cloudCommerceProcurementBaseUrl + "/entitlements/" + entitlementId + ":approveChange"
	}

	client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform")

	if clientErr != nil {
		log.Printf("Failed to create oath2 client for the procurement API %#v \n", clientErr)
		return clientErr
	}

	log.Printf("Sending entitlement approval: %s \n", procurementUrl)
	procResp, err := client.Post(procurementUrl,"",nil)
	if nil != err {
		log.Printf("Failed sending entitlement approval request %s %#v \n",procurementUrl, err)
		return err
	}
	defer procResp.Body.Close()
	if procResp.StatusCode != 200 {
		log.Println("Entitlement approval received error response: ",procResp.StatusCode)
		responseDump, _ := httputil.DumpResponse(procResp, true)
		log.Println(string(responseDump))
		return errors.New("Entitlement approval received error response: "+procResp.Status)
	} else {
		log.Printf("%s %s",procurementUrl,procResp.Status)
	}
	return nil
}

func syncEntitlement(entitlementName string) error {
	rcvdEntitlement := Entitlement{}
	if err := getEntitlement(entitlementName, &rcvdEntitlement); err == nil {
		if err := updateEntitlement(&rcvdEntitlement); err != nil {
			log.Printf("Unable to update entitlement %#v due to error %#v \n", rcvdEntitlement, err)
		}
	} else {
		log.Printf("Unable to retrieve entitlement %s due to error %#v \n", rcvdEntitlement, err)
		return err
	}
	return nil
}

func getEntitlement(entitlementId string,entitlement *Entitlement) error {
	procurementUrl := cloudCommerceProcurementBaseUrl + "/entitlements/" + entitlementId

	client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform")

	if clientErr != nil {
		log.Printf("Failed to create oath2 client for the procurement API %#v \n", clientErr)
		return clientErr
	}

	log.Printf("Getting entitlement: %s \n", procurementUrl)
	resp, err := client.Get(procurementUrl)
	if err != nil {
		log.Printf("Failed to get entitlement %s %#v \n",procurementUrl, err)
		return err
	}
	if resp.StatusCode != 200 {
		log.Println("Get entitlement received error response: ",resp.StatusCode)
		responseDump, _ := httputil.DumpResponse(resp, true)
		log.Println(string(responseDump))
		return errors.New(resp.Status)
	} else {
		log.Printf("%s %s",procurementUrl,resp.Status)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&entitlement); err != nil {
		log.Printf("Error decoding entitlement %s %#v %#v \n", procurementUrl, resp.Body, err)
		return err
	}

	return nil
}

func updateEntitlement(entitlement *Entitlement) error {
	entitlementBytes, err := json.Marshal(entitlement)
	if err != nil {
		log.Printf("Error marshalling entitlement %#v \n", err)
		return err
	}
	url := subscriptionServiceBaseUrl+"/entitlements";
	entitlementReq, err := http.NewRequest(http.MethodPut, url , bytes.NewBuffer(entitlementBytes))
	if err != nil {
		log.Printf("Failed creating entitlement update request %s %#v \n",subscriptionServiceBaseUrl, err)
		return err
	}
	entitlementResp, err := http.DefaultClient.Do(entitlementReq)
	if err != nil {
		log.Printf("Failed sending entitlement update request %s %#v \n",subscriptionServiceBaseUrl, err)
		return err
	}
	if entitlementResp.StatusCode != 204 {
		log.Println("Update entitlement received error response: ",entitlementResp.StatusCode)
		responseDump, _ := httputil.DumpResponse(entitlementResp, true)
		log.Println(string(responseDump))
		return errors.New(entitlementResp.Status)
	} else {
		log.Printf("%s %s",url,entitlementResp.Status)
	}
	defer entitlementResp.Body.Close()
	return nil
}

func deleteEntitlement(entitlementId string) error {
	url := subscriptionServiceBaseUrl+"/entitlements/"+entitlementId
	entitlementReq, err := http.NewRequest(http.MethodDelete, url,nil)
	if nil != err {
		log.Printf("Failed creating entitlement delete request %s %#v \n",subscriptionServiceBaseUrl, err)
		return err
	}
	entitlementResp, err := http.DefaultClient.Do(entitlementReq)
	if err != nil {
		log.Printf("Failed sending entitlement delete request %s %#v \n",subscriptionServiceBaseUrl, err)
		return err
	}
	if entitlementResp.StatusCode != 204 {
		log.Println("Delete entitlement received error response: ",entitlementResp.StatusCode)
		responseDump, _ := httputil.DumpResponse(entitlementResp, true)
		log.Println(string(responseDump))
		return errors.New(entitlementResp.Status)
	} else {
		log.Printf("%s %s",url,entitlementResp.Status)
	}
	defer entitlementResp.Body.Close()
	return nil
}

func syncAccount(accountName string) error {
	rcvdAcct := Account{}
	if err := getAccount(accountName, &rcvdAcct); err == nil {
		err := updateAccount(&rcvdAcct)
		if err != nil {
			log.Printf("Unable to update account %#v due to error %#v \n", rcvdAcct, err)
		}
	} else {
		log.Printf("Unable to retrieve account %s due to error %#v \n", accountName, err)
		return err
	}
	return nil
}

func accountExists(accountId string) (bool, error){
	procurementUrl := cloudCommerceProcurementBaseUrl + "/accounts/" + accountId
	client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform")

	if clientErr != nil {
		log.Printf("Failed to create oath2 client for the procurement API %#v \n", clientErr)
		return false, clientErr
	}

	log.Printf("Getting account: %s \n", procurementUrl)
	resp, err := client.Get(procurementUrl)
	if err != nil {
		log.Printf("Failed to get account %s %#v \n",procurementUrl, err)
		return false, err
	}

	if resp.StatusCode != 200 {
		return true, nil
	} else {
		return false, nil
	}
}

func getAccount(accountId string,account *Account) error {
	procurementUrl := cloudCommerceProcurementBaseUrl + "/accounts/" + accountId
	client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform")

	if clientErr != nil {
		log.Printf("Failed to create oath2 client for the procurement API %#v \n", clientErr)
		return clientErr
	}

	log.Printf("Getting account: %s \n", procurementUrl)
	resp, err := client.Get(procurementUrl)
	if err != nil {
		log.Printf("Failed to get account %s %#v \n",procurementUrl, err)
		return err
	}

	if resp.StatusCode != 200 {
		log.Println("Get account received error response: ",resp.StatusCode)
		responseDump, _ := httputil.DumpResponse(resp, true)
		log.Println(string(responseDump))
		return errors.New(resp.Status)
	} else {
		log.Printf("%s %s",procurementUrl,resp.Status)
	}

	defer resp.Body.Close()

	if errJson := json.NewDecoder(resp.Body).Decode(&account); errJson != nil {
		log.Printf("Failed to decode account %s %#v %#v \n", procurementUrl, resp.Body, err)
		return errJson
	}

	return nil
}

func updateAccount(account *Account) error {
	accountBytes, err := json.Marshal(account)
	if err != nil {
		log.Printf("Error marshalling account %#v \n", err)
		return err
	}
	url := subscriptionServiceBaseUrl+"/accounts"
	accountReq, err := http.NewRequest(http.MethodPut, url , bytes.NewBuffer(accountBytes))
	if nil != err {
		log.Printf("Failed creating account update request %s %#v \n",subscriptionServiceBaseUrl, err)
		return err
	}
	accountResp, err := http.DefaultClient.Do(accountReq)
	if err != nil {
		log.Printf("Failed sending account update request %s %#v \n",subscriptionServiceBaseUrl, err)
		return err
	}
	if accountResp.StatusCode != 204 {
		log.Println("Update account received error response: ",accountResp.StatusCode)
		responseDump, _ := httputil.DumpResponse(accountResp, true)
		log.Println(string(responseDump))
		return errors.New(accountResp.Status)
	} else {
		log.Printf("%s %s",url,accountResp.Status)
	}
	defer accountResp.Body.Close()

	return nil
}

func deleteAccount(accountId string) error {
	url := subscriptionServiceBaseUrl+"/accounts/"+accountId
	accountReq, err := http.NewRequest(http.MethodDelete, url,nil)
	if nil != err {
		log.Printf("Failed creating account delete request %s %#v \n",subscriptionServiceBaseUrl, err)
		return err
	}
	accountResp, err := http.DefaultClient.Do(accountReq)
	if err != nil {
		log.Printf("Failed sending account delete request %s %#v \n",subscriptionServiceBaseUrl, err)
		return err
	}
	if accountResp.StatusCode != 204 {
		log.Println("Delete entitlement received error response: ",accountResp.StatusCode)
		responseDump, _ := httputil.DumpResponse(accountResp, true)
		log.Println(string(responseDump))
		return errors.New(accountResp.Status)
	} else {
		log.Printf("%s %s",url,accountResp.Status)
	}
	defer accountResp.Body.Close()
	return nil
}
