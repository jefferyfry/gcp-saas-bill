package main

import (
	"github.com/cloudbees/cloud-bill-saas/datastore-backup/backup"
	"github.com/cloudbees/cloud-bill-saas/datastore-backup/config"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

func main() {
	log.Println("Starting Cloud Bill SaaS Datastore Backup Job...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %#v", err)
	}

	//wait for istio
	time.Sleep(10 * time.Second)

	if config.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryDsn,
			Environment: config.GcpProjectId,
			ServerName: "datastore-backup",
		})

		sentry.CaptureMessage("Sentry initialized for Cloud Bill SaaS Datastore Backup Job.")
		// Since sentry emits events in the background we need to make sure
		// they are sent before we shut down
		sentry.Flush(time.Second * 5)
	}

	//start service
	datastoreBackup := backup.GetDatastoreBackupHandler(config.GcpProjectId,config.GcsBucket)

	if err := datastoreBackup.Run(); err != nil {
		log.Printf("Datastore Backup Job encountered err %s",err)
	} else {
		log.Println("Datastore Backup Job completed successfully.")
	}
}


