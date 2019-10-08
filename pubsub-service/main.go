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
	log.Println("Starting Cloud Bill SaaS PubSub Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %#v", err)
	}

	if config.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryDsn,
		})

		sentry.CaptureException(errors.New("Sentry initialized for Cloud Bill SaaS Datastore Backup Job."))
		// Since sentry emits events in the background we need to make sure
		// they are sent before we shut down
		sentry.Flush(time.Second * 5)
	}

	//start service
	pubSubListener := mpevents.GetPubSubListener(config.PubSubSubscription,config.SubscriptionServiceUrl,config.CloudCommerceProcurementUrl,config.PartnerId,config.GcpProjectId)
	log.Fatal(pubSubListener.Listen())
}


