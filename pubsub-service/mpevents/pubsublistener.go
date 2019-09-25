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
	exists, errTp := topic.Exists(ctx)
	if errTp != nil {
		log.Printf("Error checking for topic: %#v", errTp)
		return errTp
	}
	if !exists {
		if _, err := client.CreateTopic(ctx, topicId); err != nil {
			log.Printf("Failed to create topic: %#v", err)
			return err
		}
	}

	subscription := client.Subscription(lstnr.PubSubSubscription)

	exists, errSub := subscription.Exists(ctx)
	if errSub != nil {
		log.Printf("Error checking for subscription: %#v", errSub)
		return errSub
	}

	if !exists {
		_, err = client.CreateSubscription(ctx, lstnr.PubSubSubscription, pubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			log.Printf("Failed to create subscription: %#v", err)
			return err
		}
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
		if processMsg(pubSubMsg) {
			msg.Ack()
			log.Println("Message acked.")
		} else {
			msg.Nack()
			log.Println("Message nacked.")
		}
	})
	if errRcv != nil {
		return errRcv
	}
	return nil
}

func processMsg(pubSubMsg PubSubMsg) bool {
	switch pubSubMsg.EventType {
	case "ACCOUNT_ACTIVE":
		err := syncAccount(pubSubMsg.Account.Id)
		if err != nil {
			log.Printf("Unable to update account %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
		//check to see if we have any entitlements we need to approve
	case "ENTITLEMENT_CREATION_REQUESTED":
		//query to get the entitlement
		//get the account name
		//query for the account, contact
		//if it exists approve
		//if not store it for later
		err := syncEntitlement(pubSubMsg.Entitlement.Id)
		if err != nil {
			log.Printf("Unable to update entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
		return false
		/*err := postEntitlementApproval(pubSubMsg.Entitlement.Id, false)
		if err != nil {
			log.Printf("Unable to approve entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		} else {
			err := syncEntitlement(pubSubMsg.Entitlement.Id)
			if err != nil {
				log.Printf("Unable to update entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
				return false
			}
		}*/
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
		err := syncEntitlement(pubSubMsg.Entitlement.Id)
		if err != nil {
			log.Printf("Unable to update entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ENTITLEMENT_DELETED":
		err := deleteEntitlement(pubSubMsg.Entitlement.Id)
		if err != nil {
			log.Printf("Unable to delete entitlement %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ACCOUNT_DELETED":
		err := deleteAccount(pubSubMsg.Account.Id)
		if err != nil {
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
	err := getEntitlement(entitlementName, &rcvdEntitlement)
	if err == nil {
		err := updateEntitlement(&rcvdEntitlement)
		if err != nil {
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

	err = json.NewDecoder(resp.Body).Decode(&entitlement)

	defer resp.Body.Close()

	if err != nil {
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
	if nil != err {
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
	err := getAccount(accountName, &rcvdAcct)
	if err == nil {
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

	errJson := json.NewDecoder(resp.Body).Decode(&account)

	defer resp.Body.Close()

	if errJson != nil {
		log.Printf("Failed to decode account %s %#v %#v \n", procurementUrl, resp.Body, err)
		return errJson
	}
	responseDump, err := httputil.DumpResponse(resp, true)
	log.Println(string(responseDump))
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
