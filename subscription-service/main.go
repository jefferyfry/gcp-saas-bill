package main

import (
	"github.com/cloudbees/cloud-bill-saas/subscription-service/config"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/dbinterface"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/web"
	"github.com/getsentry/sentry-go"
	"github.com/jefferyfry/funclog"
	"time"
)

var (
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

// @contact.name CloudBees Support
// @contact.url http://support.cloudbees.com
// @contact.email support@cloudbees.com
// @host localhost:8085
// @BasePath /api/v1
// @termsOfService https://www.cloudbees.com/products/terms-service
func main() {
	LogI.Println("Starting Cloud Bill SaaS Subscription Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		LogE.Fatalf("Invalid configuration: %v", err)
	}

	if config.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryDsn,
			Environment: config.GcpProjectId,
			ServerName: "subscription-service",
		})

		//wait for istio
		time.Sleep(10 * time.Second)

		sentry.CaptureMessage("Sentry initialized for Cloud Bill Subscription Service.")
		sentry.Flush(time.Second * 5)
	}

	datastoreClient := dbinterface.NewPersistenceLayer(dbinterface.DATASTOREDB,config.GcpProjectId)

	//start web service
	LogE.Fatal(web.SetUpService(datastoreClient,config.SubscriptionServiceEndpoint,config.HealthCheckEndpoint))
}
