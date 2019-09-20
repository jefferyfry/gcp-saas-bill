package main

import (
	"fmt"
	"github.com/cloudbees/cloud-bill-saas/frontend-service/config"
	"github.com/cloudbees/cloud-bill-saas/frontend-service/web"
	"log"
)

func main() {
	fmt.Println("Starting Cloud Bill SaaS Frontend Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	//start service
	log.Fatal(web.SetUpService(config.FrontendServiceEndpoint,config.SubscriptionServiceUrl,config.ClientId,config.ClientSecret,config.CallbackUrl,config.Issuer,config.SessionKey,config.CloudCommerceProcurementUrl,config.PartnerId))
}
