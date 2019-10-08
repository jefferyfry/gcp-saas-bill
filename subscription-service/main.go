package main

import (
	"errors"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/config"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/dbinterface"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/web"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)
// @contact.name CloudBees Support
// @contact.url http://support.cloudbees.com
// @contact.email support@cloudbees.com
// @host localhost:8085
// @BasePath /api/v1
// @termsOfService https://www.cloudbees.com/products/terms-service
func main() {
	sentry.Init(sentry.ClientOptions{
		Dsn: "https://1ae3b5e486aa45e6b23689a42bada322@sentry.io/1770543",
	})
	sentry.CaptureException(errors.New("Sentry initialized for Cloud Bill SaaS Subscription Service."))
	// Since sentry emits events in the background we need to make sure
	// they are sent before we shut down
	sentry.Flush(time.Second * 5)

	log.Println("Starting Cloud Bill SaaS Subscription Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	datastoreClient := dbinterface.NewPersistenceLayer(dbinterface.DATASTOREDB,config.GcpProjectId)

	//start service
	log.Fatal(web.SetUpService(datastoreClient,config.SubscriptionServiceEndpoint))
}
