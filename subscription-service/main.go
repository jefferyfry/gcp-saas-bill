package main

import (
	"github.com/cloudbees/cloud-bill-saas/subscription-service/config"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/dbinterface"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/web"
	"log"
)
// @contact.name CloudBees Support
// @contact.url http://support.cloudbees.com
// @contact.email support@cloudbees.com
// @host localhost:8085
// @BasePath /api/v1
// @termsOfService https://www.cloudbees.com/products/terms-service
func main() {
	log.Println("Starting Cloud Bill SaaS Subscription Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	datastoreClient := dbinterface.NewPersistenceLayer(dbinterface.DATASTOREDB,config.GcpProjectId)

	//start service
	log.Fatal(web.SetUpService(datastoreClient,config.SubscriptionServiceEndpoint))
}
