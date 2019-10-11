package main

import (
	"github.com/cloudbees/cloud-bill-saas/pubsub-service/config"
	"github.com/cloudbees/cloud-bill-saas/pubsub-service/mpevents"
	"github.com/cloudbees/cloud-bill-saas/pubsub-service/web"
	"github.com/getsentry/sentry-go"
	"github.com/jefferyfry/funclog"
	"time"
)

var (
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

func main() {
	LogI.Println("Starting Cloud Bill SaaS PubSub Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		LogE.Fatalf("Invalid configuration: %#v", err)
	}

	if config.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryDsn,
			Environment: config.GcpProjectId,
			ServerName: "pubsub-service",
		})

		//wait for istio
		time.Sleep(10 * time.Second)

		sentry.CaptureMessage("Sentry initialized for Cloud Bill PubSub Service.")
		sentry.Flush(time.Second * 5)
	}

	//start the web service
	go web.SetUpService(config.HealthCheckEndpoint,config.PubSubSubscription,config.SubscriptionServiceUrl,config.CloudCommerceProcurementUrl,config.PartnerId,config.GcpProjectId)

	//start the pub sub listener
	pubSubListener := mpevents.GetPubSubListener(config.PubSubSubscription,config.SubscriptionServiceUrl,config.CloudCommerceProcurementUrl,config.PartnerId,config.GcpProjectId)
	LogE.Fatal(pubSubListener.Listen())
}


