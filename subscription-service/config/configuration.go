package config

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
)

var (
	SubscriptionServiceEndpoint 		= ":8085"
	GcpProjectId				        = "cloud-bill-saas"
)

type ServiceConfig struct {
	SubscriptionServiceEndpoint    	string	`json:"subscriptionServiceEndpoint"`
	GcpProjectId    				string	`json:"gcpProjectId"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		SubscriptionServiceEndpoint,
		GcpProjectId,
	}

	if dir, err := os.Getwd(); err != nil {
		log.Println("Unable to determine working directory.")
		return conf, err
	} else {
		log.Printf("Running service with working directory %s \n", dir)
	}

	//parse commandline arguments
	configFile := flag.String("configFile", "", "set the path to the configuration json file")
	subscriptionServiceEndpoint := flag.String("subscriptionServiceEndpoint", "", "set the value of this service endpoint")
	gcpProjectId := flag.String("gcpProjectId", "", "set the GCP Project Id")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CLOUD_BILL_SUBSCRIPTION_CONFIG_FILE")
	}
	if *subscriptionServiceEndpoint == "" {
		*subscriptionServiceEndpoint = os.Getenv("CLOUD_BILL_SUBSCRIPTION_SERVICE_ENDPOINT")
	}
	if *gcpProjectId == "" {
		*gcpProjectId = os.Getenv("CLOUD_BILL_SUBSCRIPTION_GCP_PROJECT_ID")
	}


	if *configFile == "" {
		//try other flags
		conf.SubscriptionServiceEndpoint = *subscriptionServiceEndpoint
		conf.GcpProjectId = *gcpProjectId
	} else {
		if file, err := os.Open(*configFile); err != nil {
			log.Printf("Error reading confile file %s %s", *configFile, err)
			return conf, err
		} else {
			if err = json.NewDecoder(file).Decode(&conf); err != nil {
				return conf, errors.New("Configuration file not found.")
			}
			log.Printf("Using confile file %s \n", *configFile)
		}
	}

	valid := true

	if conf.SubscriptionServiceEndpoint == "" {
		log.Println("SubscriptionServiceEndpoint was not set.")
		valid = false
	}

	if conf.GcpProjectId == "" {
		log.Println("GcpProjectId was not set.")
		valid = false
	}

	if credPath,envExists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); !envExists {
		log.Println("GOOGLE_APPLICATION_CREDENTIALS was not set. This is fine with an emulator but will fail in production. ")
	} else {
		if _, errPath := os.Stat(credPath); os.IsNotExist(errPath) {
			log.Println("GOOGLE_APPLICATION_CREDENTIALS file does not exist: ", credPath)
			valid = false
		} else {
			log.Println("Using GOOGLE_APPLICATION_CREDENTIALS file: ", credPath)
		}
	}

	if !valid {
		return conf, errors.New("Subscription service configuration is not valid!")
	} else {
		return conf, nil
	}
}
