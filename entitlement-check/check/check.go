package check

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jefferyfry/funclog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"net/http/httputil"
	"strings"
)

var (
	subscriptionServiceBaseUrl string
	googleSubscriptionsBaseUrl string

	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

type EntitlementCheckHandler struct {
	Products    			string
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

type Subscription struct {
	Name 				string     	`json:"name"`
	ExternalAccountId 	string     	`json:"externalAccountId"`
	Version				string     	`json:"version,omitempty"`
	Status				string     	`json:"status"`
	SubscribedResources	[]SubscribedResource     	`json:"subscribedResources"`
	RequiredApprovals	string     	`json:"requiredApprovals,omitempty"`
	StartDate			json.RawMessage     	`json:"startDate,omitempty"`
	EndDate				json.RawMessage     	`json:"endDate,omitempty"`
	CreateTime			string     	`json:"createTime,omitempty"`
	UpdateTime			string     	`json:"updateTime,omitempty"`
}

type SubscribedResource struct {
	SubscriptionProvider 	string     	`json:"subscriptionProvider"`
	Resource 				string     	`json:"resource"`
	Labels					json.RawMessage     	`json:"labels,omitempty"`
}

func GetEntitlementCheckHandler(products string, subscriptionServiceUrl string, googleSubscriptionsUrl string) *EntitlementCheckHandler {
	subscriptionServiceBaseUrl = subscriptionServiceUrl
	googleSubscriptionsBaseUrl = googleSubscriptionsUrl
	return &EntitlementCheckHandler{
		products,
	}
}

func (hdlr *EntitlementCheckHandler) Run() error {
	//query subscription service for entitlements
	products := strings.Split(hdlr.Products, ",")

	for _, product := range products {
		LogI.Printf("Checking entitlements for product %s", product)
		if entitlements, err := getActiveEntitlementsForProduct(product); err == nil {
			for _, entitlement := range entitlements {
				LogI.Printf("Checking entitlement %s", entitlement.Id)
				if entitlementStatus, err := getProdEntitlementStatus(entitlement.Id); err == nil {
					status := "ENTITLEMENT_"+entitlementStatus
					if status != entitlement.State {
						entitlement.State = status
						saveEntitlementToDb(&entitlement)
						LogI.Printf("Updated entitlement %s with status %s", entitlement.Id, entitlement.State)
					} else {
						LogI.Printf("Entitlement %s status with status %s is unchanged.", entitlement.Id, status)
					}
				} else {
					LogE.Printf("Failed to get entitlement status %s %#v \n",googleSubscriptionsBaseUrl, err)
				}
			}
		} else {
			LogE.Printf("Failed to get entitlement %s %#v \n",subscriptionServiceBaseUrl, err)
			return err
		}
	}

	return nil
}

func getActiveEntitlementsForProduct(product string) ([]Entitlement, error) {
	subscriptionServiceUrl := subscriptionServiceBaseUrl + "/entitlements?filters=state%3DENTITLEMENT_ACTIVE%2Cproduct%3D"+product

	LogI.Printf("Getting entitlements: %s \n", subscriptionServiceUrl)
	resp, err := http.Get(subscriptionServiceUrl)
	if err != nil {
		LogE.Printf("Failed to get entitlements %s %#v \n",subscriptionServiceUrl, err)
		return nil,err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		LogI.Printf("No active entitlements for %s found.",product)
	} else if resp.StatusCode != 200 {
		LogE.Println("Get entitlement received error response: ",resp.StatusCode)
		responseDump, _ := httputil.DumpResponse(resp, true)
		LogE.Println(string(responseDump))
		return nil,errors.New(resp.Status)
	} else {
		LogI.Printf("Got entitlements %s %s",subscriptionServiceUrl,resp.Status)
	}

	entitlements := make([]Entitlement,0)
	if err = json.NewDecoder(resp.Body).Decode(&entitlements); err != nil {
		LogE.Printf("Error decoding entitlements %s %#v %#v \n", subscriptionServiceUrl, resp.Body, err)
		return nil,err
	}

	return entitlements,nil
}

func getProdEntitlementStatus(entitlementId string) (string, error) {
	subscriptionsUrl := googleSubscriptionsBaseUrl+"/partnerSubscriptions/"+entitlementId
	if client, clientErr := google.DefaultClient(oauth2.NoContext, "https://www.googleapis.com/auth/cloud-platform"); clientErr != nil {
		LogE.Printf("Failed to create oath2 client for the procurement API %#v \n", clientErr)
		return "", clientErr
	} else {
		LogI.Printf("Getting subscription entitlement : %s \n", subscriptionsUrl)
		if subResp, err := client.Get(subscriptionsUrl); nil != err {
			LogE.Printf("Failed subscription entitlement request %s %#v \n", subscriptionsUrl, err)
			return "", err
		} else {
			defer subResp.Body.Close()
			if subResp.StatusCode != 200 {
				LogE.Println("Getting subscription entitlement received error response: ", subResp.StatusCode)
				responseDump, _ := httputil.DumpResponse(subResp, true)
				LogE.Println(string(responseDump))
				return "",errors.New("Getting subscription entitlement received error response: " + subResp.Status)
			}
			subscription := Subscription{}
			if err = json.NewDecoder(subResp.Body).Decode(&subscription); err != nil {
				LogE.Printf("Error decoding subscription %s %#v %#v \n", subscriptionsUrl, subResp.Body, err)
				responseDump, _ := httputil.DumpResponse(subResp, true)
				LogE.Println(string(responseDump))
				return "",err
			}
			return subscription.Status, nil
		}
	}
}

func saveEntitlementToDb(entitlement *Entitlement) error {
	entitlementBytes, err := json.Marshal(entitlement)
	if err != nil {
		LogE.Printf("Error marshalling entitlement %#v \n", err)
		return err
	}
	url := subscriptionServiceBaseUrl+"/entitlements";
	entitlementReq, err := http.NewRequest(http.MethodPut, url , bytes.NewBuffer(entitlementBytes))
	if err != nil {
		LogE.Printf("Failed creating entitlement update request %s %#v \n",subscriptionServiceBaseUrl, err)
		return err
	}
	entitlementResp, err := http.DefaultClient.Do(entitlementReq)
	if err != nil {
		LogE.Printf("Failed sending entitlement update request %s %#v \n",subscriptionServiceBaseUrl, err)
		return err
	}
	defer entitlementResp.Body.Close()
	if entitlementResp.StatusCode != 204 {
		LogE.Println("Update entitlement received error response: ",entitlementResp.StatusCode)
		responseDump, _ := httputil.DumpResponse(entitlementResp, true)
		LogE.Println(string(responseDump))
		return errors.New(entitlementResp.Status)
	} else {
		LogI.Printf("Updated entitlement %s %s",url,entitlementResp.Status)
	}
	return nil
}

