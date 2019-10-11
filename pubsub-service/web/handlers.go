package web

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/jefferyfry/funclog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

type PubSubServiceHandler struct {
	PubSubSubscription    			string
	SubscriptionServiceUrl 			string
	CloudCommerceProcurementUrl    	string
	PartnerId 						string
	GcpProjectId                    string
}

var (
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

func GetPubSubServiceHandler(pubSubSubscription string, subscriptionServiceUrl string, cloudCommerceProcurementUrl string, partnerId string, gcpProjectId string) *PubSubServiceHandler {
	return &PubSubServiceHandler{
		pubSubSubscription,
		subscriptionServiceUrl,
		cloudCommerceProcurementUrl,
		partnerId,
		gcpProjectId,
	}
}

func (hdlr *PubSubServiceHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	subscriptionServiceUrl := hdlr.SubscriptionServiceUrl+"/healthz"
	if subResp, err := http.Get(subscriptionServiceUrl); err == nil {
		if subResp.StatusCode == http.StatusOK {
			procurementUrl := hdlr.CloudCommerceProcurementUrl +  "/providers/" +  hdlr.PartnerId + "/accounts/"

			if client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform"); clientErr != nil {
				LogE.Printf("Healthz failed. Failed to create oath2 client for the procurement API %#v \n", clientErr)
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				if procResp, err := client.Get(procurementUrl); nil != err {
					LogE.Printf("Healthz failed. Cloud Commerce API check failed: %s %#v \n", procurementUrl, err)
					http.Error(w,err.Error(),procResp.StatusCode)
				} else {
					ctx := context.Background()
					client, err := pubsub.NewClient(ctx, hdlr.GcpProjectId)
					if err != nil {
						LogE.Fatalf("Healthz failed. Error creating pubsub client %s: %#v", hdlr.PubSubSubscription, err)
					}

					subscription := client.Subscription(hdlr.PubSubSubscription)

					if exists, errSub := subscription.Exists(ctx); !exists && errSub == nil {
						LogE.Printf("Healthz failed. Marketplace subscription %s does not exist \n", subscription.String())
						w.WriteHeader(http.StatusNotFound)
					} else if errSub != nil{
						LogE.Printf("Healthz failed. Error checking for subscription: %#v", errSub)
						w.WriteHeader(http.StatusInternalServerError)
					}
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	} else {
		LogE.Printf("Healthz failed. Subscription Service check failed: %s %s %#v \n", subscriptionServiceUrl,subResp.StatusCode,err)
		http.Error(w,err.Error(),subResp.StatusCode)
	}
}








