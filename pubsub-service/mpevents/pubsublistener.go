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
	"path/filepath"
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
	Id  			string     	`json:"id"`
	Name  			string     	`json:"name"`
	UpdateTime   	string    	`json:"updateTime,omitempty"`
	CreateTime      string    	`json:"createTime,omitempty"`
	Provider     	string		`json:"provider,omitempty"`
	State 	 		string      `json:"state,omitempty"`
	Approvals    	[]Approval  `json:"approvals,omitempty"`
}

type Approval struct {
	Name  			string     	`json:"name"`
	State  			string     	`json:"state"`
	Reason  		string     	`json:"reason"`
	UpdateTime  	string     	`json:"updateTime"`
}

type Entitlement struct {
	Id     				string	`json:"id"`
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
	SubscriptionServiceUrl 			string
	CloudCommerceProcurementUrl    	string
	PartnerId    					string
	GcpProjectId                    string
}

var (
	cloudCommerceProcurementBaseUrl string
	subscriptionServiceBaseUrl string
)

func GetPubSubListener(pubSubSubscription string, subscriptionServiceUrl string, cloudCommerceProcurementUrl string, partnerId string, gcpProjectId string) *PubSubListener {
	cloudCommerceProcurementBaseUrl = cloudCommerceProcurementUrl + "/providers/" + partnerId
	subscriptionServiceBaseUrl = subscriptionServiceUrl
	return &PubSubListener{
		pubSubSubscription,
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
		log.Fatalf("Error creating pubsub client %s: %#v", lstnr.PubSubSubscription, err)
	}

	subscription := client.Subscription(lstnr.PubSubSubscription)

	if exists, errSub := subscription.Exists(ctx); !exists && errSub == nil {
		log.Fatalf("Marketplace subscription %s does not exist \n", subscription.String())
	} else if errSub != nil{
		log.Fatalf("Error checking for subscription: %#v", errSub)
	}

	log.Printf("Begin receiving messages from subscription %s \n", subscription.String())
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
		if _,err := syncAccount(pubSubMsg.Account.Id); err != nil {
			log.Printf("Unable to update account %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		} else {
			if entitlements, err := getUnapprovedEntitlementsFromDb(pubSubMsg.Account.Id); err != nil {
				log.Printf("Unable to check for unapproved entitlements %#v due to error %#v \n", pubSubMsg.Entitlement, err)
				return false
			} else if entitlements!=nil && len(entitlements) > 0 {
				for _, ent := range entitlements {
					postEntitlementApprovalToCommerceApi(ent.Id)
				}
			} else {
				log.Printf("No unapproved entitlments were found for account %s",pubSubMsg.Account.Id)
			}
		}
	case "ENTITLEMENT_CREATION_REQUESTED":
		if entitlement,err := syncEntitlement(pubSubMsg.Entitlement.Id); err != nil {
			if accountExists, acctErr := accountExistsInDb(entitlement.Account); acctErr != nil {
				log.Printf("Unable to determine if account %#v exists due to error %#v \n", entitlement.Account, acctErr)
				return false
			} else if accountExists {
				if postErr := postEntitlementApprovalToCommerceApi(pubSubMsg.Entitlement.Id); postErr != nil {
					log.Printf("Unable to approve entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, postErr)
					return false
				}
			}
		}
	case "ENTITLEMENT_ACTIVE":
		if _,err := syncEntitlement(pubSubMsg.Entitlement.Id); err != nil {
			log.Printf("Unable to update entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ENTITLEMENT_PLAN_CHANGE_REQUESTED":
		if err := postEntitlementChangeApprovalToCommerceApi(pubSubMsg.Entitlement.Id); err == nil {
			if _,err := syncEntitlement(pubSubMsg.Entitlement.Id); err != nil {
				log.Printf("Unable to update entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
				return false
			}
		} else {
			log.Printf("Unable to approve entitlement plan change %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ENTITLEMENT_PLAN_CHANGED":
		if _,err := syncEntitlement(pubSubMsg.Entitlement.Id); err != nil {
			log.Printf("Unable to update entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ENTITLEMENT_PENDING_CANCELLATION":
		log.Println("ENTITLEMENT_PENDING_CANCELLATION ignored.")
	case "ENTITLEMENT_CANCELLED":
		if _,err := syncEntitlement(pubSubMsg.Entitlement.Id); err != nil {
			log.Printf("Unable to update entitlement plan %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ENTITLEMENT_DELETED":
		if err := deleteEntitlementFromDb(pubSubMsg.Entitlement.Id); err != nil {
			log.Printf("Unable to delete entitlement %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "ACCOUNT_DELETED":
		if err := deleteAccountFromDb(pubSubMsg.Account.Id); err != nil {
			log.Printf("Unable to delete account %#v due to error %#v \n", pubSubMsg.Entitlement, err)
			return false
		}
	case "TEST":
		log.Printf("Test message %s was received \n", pubSubMsg.EventId)
	default:
		log.Printf("Unknown pubsub event type %#v \n", pubSubMsg)
	}
	return true
}

func postEntitlementApprovalToCommerceApi(entitlementId string) error {
	procurementUrl := cloudCommerceProcurementBaseUrl + "/entitlements/" + entitlementId + ":approve"

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

func postEntitlementChangeApprovalToCommerceApi(entitlementId string) error {
	entitlement := Entitlement{}
	if err := getEntitlementFromCommerceApi(entitlementId,&entitlement); err != nil {
		log.Printf("Unable to determine entitlement to approve from procurement API %#v \n", err)
	}
	jsonApproval := []byte(`{
				"pendingPlanName": "`+entitlement.NewPendingPlan+`"
			}`)
	procurementUrl := cloudCommerceProcurementBaseUrl + "/entitlements/" + entitlementId + ":approvePlanChange"

	client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform")

	if clientErr != nil {
		log.Printf("Failed to create oath2 client for the procurement API %#v \n", clientErr)
		return clientErr
	}

	log.Printf("Sending entitlement change approval: %s %s \n", procurementUrl, jsonApproval)
	procResp, err := client.Post(procurementUrl,"",bytes.NewBuffer(jsonApproval))
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

func syncEntitlement(entitlementId string) (*Entitlement,error) {
	entitlement := Entitlement{}
	if err := getEntitlementFromCommerceApi(entitlementId, &entitlement); err == nil {
		entitlement.Account = filepath.Base(entitlement.Account)
		if err := saveEntitlementToDb(&entitlement); err != nil {
			log.Printf("Unable to update entitlement %#v due to error %#v \n", entitlement, err)
		}
	} else {
		log.Printf("Unable to retrieve entitlement %s due to error %#v \n", entitlement, err)
		return nil, err
	}
	return &entitlement,nil
}

func getEntitlementFromCommerceApi(entitlementId string,entitlement *Entitlement) error {
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
	entitlement.Id = entitlementId

	return nil
}

func getUnapprovedEntitlementsFromDb(accountId string) ([]Entitlement, error) {
	procurementUrl := subscriptionServiceBaseUrl + "/accounts/"+accountId+"/entitlements?filter?state=ENTITLEMENT_CREATION_REQUESTED"

	log.Printf("Getting unapproved entitlements: %s \n", procurementUrl)
	resp, err := http.Get(procurementUrl)
	if err != nil {
		log.Printf("Failed to get entitlement %s %#v \n",procurementUrl, err)
		return nil,err
	}
	if resp.StatusCode != 200 {
		log.Println("Get entitlement received error response: ",resp.StatusCode)
		responseDump, _ := httputil.DumpResponse(resp, true)
		log.Println(string(responseDump))
		return nil,errors.New(resp.Status)
	} else {
		log.Printf("%s %s",procurementUrl,resp.Status)
	}
	defer resp.Body.Close()

	entitlements := make([]Entitlement,0)
	if err = json.NewDecoder(resp.Body).Decode(&entitlements); err != nil {
		log.Printf("Error decoding entitlements %s %#v %#v \n", procurementUrl, resp.Body, err)
		return nil,err
	}

	return entitlements,nil
}

func saveEntitlementToDb(entitlement *Entitlement) error {
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

func deleteEntitlementFromDb(entitlementId string) error {
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

func syncAccount(accountId string) (*Account, error) {
	account := Account{}
	if err := getAccountFromCommerceApi(accountId, &account); err == nil {
		err := saveAccountToDb(&account)
		if err != nil {
			log.Printf("Unable to update account %#v due to error %#v \n", account, err)
		}
	} else {
		log.Printf("Unable to retrieve account %s due to error %#v \n", accountId, err)
		return nil,err
	}
	return &account,nil
}

func accountExistsInDb(accountId string) (bool, error){
	subscriptionServiceUrl := subscriptionServiceBaseUrl+"/accounts/"+accountId

	log.Printf("Getting account: %s \n", subscriptionServiceUrl)
	resp, err := http.Get(subscriptionServiceUrl)
	if err != nil {
		log.Printf("Failed to get account %s %#v \n",subscriptionServiceUrl, err)
		return false, err
	}

	if resp.StatusCode != 200 {
		return true, nil
	} else {
		return false, nil
	}
}

func getAccountFromCommerceApi(accountId string,account *Account) error {
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
	account.Id = accountId

	return nil
}

func saveAccountToDb(account *Account) error {
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

func deleteAccountFromDb(accountId string) error {
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
