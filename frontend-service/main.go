package main

import (
	"errors"
	"github.com/cloudbees/cloud-bill-saas/frontend-service/config"
	"github.com/cloudbees/cloud-bill-saas/frontend-service/web"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

func main() {
	log.Println("Starting Cloud Bill SaaS Frontend Service...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
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
	log.Fatal(web.SetUpService(config.FrontendServiceEndpoint,config.SubscriptionServiceUrl,config.ClientId,config.ClientSecret,config.CallbackUrl,config.Issuer,config.SessionKey,config.CloudCommerceProcurementUrl,config.PartnerId,config.FinishUrl,config.FinishUrlTitle,config.TestMode))
}
