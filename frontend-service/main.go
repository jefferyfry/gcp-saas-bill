package main

import (
	"github.com/cloudbees/cloud-bill-saas/frontend-service/config"
	"github.com/cloudbees/cloud-bill-saas/frontend-service/web"
	"github.com/getsentry/sentry-go"
	"github.com/jefferyfry/funclog"
	"time"
)

var (
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

func main() {
	LogI.Println("Starting Cloud Bill SaaS Frontend Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		LogE.Fatalf("Invalid configuration: %v", err)
	}

	if config.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryDsn,
			Environment: config.GcpProjectId,
			ServerName: "frontend-service",
			Debug:true,
		})

		//wait for istio
		time.Sleep(10 * time.Second)

		sentry.CaptureMessage("Sentry initialized for Cloud Bill SaaS Datastore Backup Job.")
		sentry.Flush(time.Second * 5)
	}

	//start web service
	LogE.Fatal(web.SetUpService(config.FrontendServiceEndpoint,config.HealthCheckEndpoint,config.SubscriptionServiceUrl,config.GoogleSubscriptionsUrl,config.ClientId,config.ClientSecret,config.CallbackUrl,config.Issuer,config.SessionKey,config.CloudCommerceProcurementUrl,config.PartnerId,config.FinishUrl,config.FinishUrlTitle,config.TestMode))
}
