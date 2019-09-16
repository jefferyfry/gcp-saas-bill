package main

import (
	"fmt"
	"github.com/cloudbees/jenkins-support-saas/frontend-service/config"
	"github.com/cloudbees/jenkins-support-saas/frontend-service/web"
	"log"
)

func main() {
	fmt.Println("Starting Jenkins Support SaaS Subscription Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	//start service
	log.Fatal(web.SetUpService(config.FrontendServiceEndpoint,config.SubscriptionServiceUrl,config.ClientId,config.ClientSecret,config.CallbackUrl,config.Issuer,config.SessionKey))
}
