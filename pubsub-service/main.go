package main

import (
	"github.com/cloudbees/cloud-bill-saas/pubsub-service/config"
	"github.com/cloudbees/cloud-bill-saas/pubsub-service/mpevents"
	"log"
)

func main() {
	log.Println("Starting Cloud Bill SaaS PubSub Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %#v", err)
	}

	//start service
	pubSubListener := mpevents.GetPubSubListener(config.PubSubSubscription,config.SubscriptionServiceUrl,config.CloudCommerceProcurementUrl,config.PartnerId,config.GcpProjectId)
	log.Fatal(pubSubListener.Listen())
}


