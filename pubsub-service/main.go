package main

import (
	"errors"
	"github.com/cloudbees/cloud-bill-saas/pubsub-service/config"
	"github.com/cloudbees/cloud-bill-saas/pubsub-service/mpevents"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

func main() {
	sentry.Init(sentry.ClientOptions{
		Dsn: "https://1ae3b5e486aa45e6b23689a42bada322@sentry.io/1770543",
	})
	sentry.CaptureException(errors.New("Sentry initialized for Cloud Bill SaaS PubSub Service."))
	// Since sentry emits events in the background we need to make sure
	// they are sent before we shut down
	sentry.Flush(time.Second * 5)

	log.Println("Starting Cloud Bill SaaS PubSub Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %#v", err)
	}

	//start service
	pubSubListener := mpevents.GetPubSubListener(config.PubSubSubscription,config.SubscriptionServiceUrl,config.CloudCommerceProcurementUrl,config.PartnerId,config.GcpProjectId)
	log.Fatal(pubSubListener.Listen())
}


