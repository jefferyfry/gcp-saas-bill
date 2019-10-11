package config

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/jefferyfry/funclog"
	"os"
)

var (
	GcpProjectId	= "cloud-billing-saas"
	GcsBucket		= "gs://bucket"
	SentryDsn		= ""

	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

type ServiceConfig struct {
	GcpProjectId    string	`json:"gcpProjectId"`
	GcsBucket    	string	`json:"gcsBucket"`
	SentryDsn		string	`json:"sentryDsn"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		GcpProjectId,
		GcsBucket,
		SentryDsn,
	}

	if dir, err := os.Getwd(); err != nil {
		LogE.Println("Unable to determine working directory.")
		return conf, err
	} else {
		LogI.Printf("Running service with working directory %s \n", dir)
	}

	//parse commandline arguments
	configFile := flag.String("configFile", "", "set the path to the configuration json file")
	gcpProjectId := flag.String("gcpProjectId", "", "set the GCP Project Id")
	gcsBucket := flag.String("gcsBucket", "", "set the GCS bucket")
	sentryDsn := flag.String("sentryDsn", "", "set the Sentry DSN")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CLOUD_BILL_DATASTORE_BACKUP_CONFIG_FILE")
	}
	if *gcpProjectId == "" {
		*gcpProjectId = os.Getenv("CLOUD_BILL_DATASTORE_BACKUP_GCP_PROJECT_ID")
	}

	if *gcsBucket == "" {
		*gcsBucket = os.Getenv("CLOUD_BILL_DATASTORE_BACKUP_GCS_BUCKET")
	}

	if *sentryDsn == "" {
		*sentryDsn = os.Getenv("CLOUD_BILL_DATASTORE_BACKUP_SENTRY_DSN")
	}

	if *configFile == "" {
		//try other flags
		conf.GcpProjectId = *gcpProjectId
		conf.GcsBucket = *gcsBucket
		conf.SentryDsn = *sentryDsn
	} else {
		if file, err := os.Open(*configFile); err != nil {
			LogE.Printf("Error reading confile file %s %s", *configFile, err)
			return conf, err
		} else {
			if err = json.NewDecoder(file).Decode(&conf); err != nil {
				return conf, errors.New("Configuration file not found.")
			}
			LogI.Printf("Using confile file %s \n", *configFile)
		}
	}

	valid := true

	if conf.GcpProjectId == "" {
		LogE.Println("GcpProjectId was not set.")
		valid = false
	}

	if conf.GcsBucket == "" {
		LogE.Println("GcsBucket was not set.")
		valid = false
	}

	if conf.SentryDsn == "" {
		LogE.Println("SentryDsn was not set. Will run without Sentry.")
	}

	if gAppCredPath,gAppCredExists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); !gAppCredExists {
		LogE.Println("GOOGLE_APPLICATION_CREDENTIALS was not set. ")
		valid = false
	} else {
		if _, gAppCredPathErr := os.Stat(gAppCredPath); os.IsNotExist(gAppCredPathErr) {
			LogE.Println("GOOGLE_APPLICATION_CREDENTIALS file does not exist: ", gAppCredPath)
			valid = false
		} else {
			LogI.Println("Using GOOGLE_APPLICATION_CREDENTIALS file: ", gAppCredPath)
		}
	}

	if !valid {
		return conf, errors.New("Datastore backup configuration is not valid!")
	} else {
		return conf, nil
	}
}
