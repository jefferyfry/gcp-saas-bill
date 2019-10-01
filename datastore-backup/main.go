package main

import (
	"github.com/cloudbees/cloud-bill-saas/datastore-backup/config"
	"github.com/cloudbees/cloud-bill-saas/datastore-backup/backup"
	"log"
)

func main() {
	log.Println("Starting Cloud Bill SaaS Datastore Backup Job...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %#v", err)
	}

	//start service
	datastoreBackup := backup.GetDatastoreBackupHandler(config.GcpProjectId,config.GcsBucket)
	if err := datastoreBackup.Run(); err != nil {
		log.Printf("Datastore Backup Job encounter err %s",err)
	} else {
		log.Println("Datastore Backup Job completed succesfully.")
	}
}


