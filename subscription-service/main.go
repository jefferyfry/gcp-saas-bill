package main

import (
	"fmt"
	"github.com/cloudbees/jenkins-support-saas/subscription-service/config"
	"github.com/cloudbees/jenkins-support-saas/subscription-service/persistence/dbinterface"
	"github.com/cloudbees/jenkins-support-saas/subscription-service/rest"
	"log"
)
// @contact.name CloudBees Support
// @contact.url http://support.cloudbees.com
// @contact.email support@cloudbees.com
// @host localhost:8085
// @BasePath /api/v1
// @termsOfService https://www.cloudbees.com/products/terms-service
func main() {
	fmt.Println("Starting Jenkins Support SaaS Subscription Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	datastoreClient := dbinterface.NewPersistenceLayer(dbinterface.DATASTOREDB,config.GoogleProjectId)

	//start service
	log.Fatal(rest.SetUpService(datastoreClient,config.SubscriptionServiceEndpoint,config.CloudCommerceProcurementUrl,config.PartnerId))
}
