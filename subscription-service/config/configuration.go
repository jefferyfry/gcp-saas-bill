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
	GcpProjectId				        = "cloud-billing-saas"
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

	dir, err := os.Getwd()
	if err != nil {
		log.Println("Unable to determine working directory.")
		return conf, err
	}
	log.Printf("Running subscription service with working directory %s \n",dir)

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
		file, err := os.Open(*configFile)
		if err != nil {
			log.Printf("Error reading confile file %s %s", *configFile, err)
			return conf, err
		}

		err = json.NewDecoder(file).Decode(&conf)
		if err != nil {
			log.Println("Configuration file not found. Continuing with default values.")
			return conf, err
		}
		log.Printf("Using confile file %s to launch subscription service \n", *configFile)
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

	credPath,envExists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !envExists {
		log.Println("GOOGLE_APPLICATION_CREDENTIALS was not set or path does not exist. This is fine with an emulator but will fail in production. ")
	}

	_, errPath := os.Stat(credPath)
	if os.IsNotExist(errPath) {
		log.Println("GOOGLE_APPLICATION_CREDENTIALS file does not exist: %s.",credPath)
		valid = false
	}

	if !valid {
		err = errors.New("Subscription service configuration is not valid!")
	}

	return conf, err
}
