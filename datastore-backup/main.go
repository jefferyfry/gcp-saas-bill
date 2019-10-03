package main

import (
	"fmt"
	"github.com/cloudbees/cloud-bill-saas/datastore-backup/config"
	"github.com/cloudbees/cloud-bill-saas/datastore-backup/backup"
	"log"
	"time"
)

func main() {
	log.Println("Starting Cloud Bill SaaS Datastore Backup Job...")
	config, err := config.GetConfiguration()

	if err != nil {
		log.Fatalf("Invalid configuration: %#v", err)
	}

	//start service
	datastoreBackup := backup.GetDatastoreBackupHandler(config.GcpProjectId,config.GcsBucket)

	fmt.Printf("Current Unix Time: %v waiting to start backup...\n", time.Now().Unix())

	time.Sleep(60 * time.Second)

	fmt.Printf("Current Unix Time: %v starting backup now\n", time.Now().Unix())

	if err := datastoreBackup.Run(); err != nil {
		log.Printf("Datastore Backup Job encountered err %s",err)
	} else {
		log.Println("Datastore Backup Job completed successfully.")
	}
}


