package main

import (
	"github.com/cloudbees/cloud-bill-saas/entitlement-check/check"
	"github.com/cloudbees/cloud-bill-saas/entitlement-check/config"
	"github.com/getsentry/sentry-go"
	"github.com/jefferyfry/funclog"
	"time"
)

var (
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

func main() {
	LogI.Println("Starting Cloud Bill SaaS Entitlement Check Job...")
	config, err := config.GetConfiguration()

	if err != nil {
		LogE.Fatalf("Invalid configuration: %#v", err)
	}

	//wait for istio
	time.Sleep(10 * time.Second)

	if config.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryDsn,
			Environment: config.GcpProjectId,
			ServerName: "entitlement-check",
		})

		sentry.CaptureMessage("Sentry initialized for Cloud Bill SaaS Entitlement Check Job.")
		// Since sentry emits events in the background we need to make sure
		// they are sent before we shut down
		sentry.Flush(time.Second * 5)
	}

	//start service
	entitlementCheck := check.GetEntitlementCheckHandler(config.Products,config.SubscriptionServiceUrl,config.GoogleSubscriptionsUrl)

	if err := entitlementCheck.Run(); err != nil {
		LogE.Printf("Entitlement Check Job encountered err %s",err)
	} else {
		LogI.Println("Entitlement Check Job completed successfully.")
	}
}


