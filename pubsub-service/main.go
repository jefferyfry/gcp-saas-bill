package main

import (
	"github.com/cloudbees/cloud-bill-saas/pubsub-service/config"
	"github.com/cloudbees/cloud-bill-saas/pubsub-service/mpevents"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

func main() {
	log.Println("Starting Cloud Bill SaaS PubSub Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %#v", err)
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

	//start service
	pubSubListener := mpevents.GetPubSubListener(config.PubSubSubscription,config.SubscriptionServiceUrl,config.CloudCommerceProcurementUrl,config.PartnerId,config.GcpProjectId)
	log.Fatal(pubSubListener.Listen())
}


