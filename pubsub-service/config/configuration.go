package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	PubSubSubscription       			= "codelab"
	PubSubTopicPrefix       			= "DEMO-"
	CloudCommerceProcurementUrl       	= "https://cloudcommerceprocurement.googleapis.com/"
	PartnerId							= "000"
	GcpProjectId				        = "cloud-billing-saas"
)

type ServiceConfig struct {
	PubSubSubscription    			string	`json:"pubSubSubscription"`
	PubSubTopicPrefix				string 	`json:"pubSubTopicPrefix"`
	CloudCommerceProcurementUrl    	string	`json:"cloudCommerceProcurementUrl"`
	PartnerId    					string	`json:"partnerId"`
	GcpProjectId    				string	`json:"gcpProjectId"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		PubSubSubscription,
		PubSubTopicPrefix,
		CloudCommerceProcurementUrl,
		PartnerId,
		GcpProjectId,
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to determine working directory.")
		return conf, err
	}
	fmt.Printf("Running subscription service with working directory %s \n",dir)

	//parse commandline arguments
	configFile := flag.String("configFile", "", "set the path to the configuration json file")
	pubSubSubscription := flag.String("pubSubSubscription", "", "set the value of the pubsub subscription")
	pubSubTopicPrefix := flag.String("pubSubTopicPrefix", "", "set the value of the pubsub topic prefix")
	cloudCommerceProcurementUrl := flag.String("cloudCommerceProcurementUrl", "", "set root url for the cloud commerce procurement API")
	partnerId := flag.String("partnerId", "", "set the CloudBees Partner Id")
	gcpProjectId := flag.String("gcpProjectId", "", "set the GCP Project Id")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CLOUD_BILLING_PUBSUB_CONFIG_FILE")
	}
	if *pubSubSubscription == "" {
		*pubSubSubscription = os.Getenv("CLOUD_BILLING_PUBSUB_SUBSCRIPTION")
	}
	if *pubSubTopicPrefix == "" {
		*pubSubTopicPrefix = os.Getenv("CLOUD_BILLING_PUBSUB_TOPIC_PREFIX")
	}
	if *cloudCommerceProcurementUrl == "" {
		*cloudCommerceProcurementUrl = os.Getenv("CLOUD_BILLING_PUBSUB_CLOUD_COMMERCE_PROCUREMENT_URL")
	}
	if *partnerId == "" {
		*partnerId = os.Getenv("CLOUD_BILLING_PUBSUB_PARTNER_ID")
	}
	if *gcpProjectId == "" {
		*gcpProjectId = os.Getenv("CLOUD_BILLING_PUBSUB_GCP_PROJECT_ID")
	}

	if *configFile == "" {
		//try other flags
		conf.PubSubSubscription = *pubSubSubscription
		conf.PubSubTopicPrefix = *pubSubTopicPrefix
		conf.CloudCommerceProcurementUrl = *cloudCommerceProcurementUrl
		conf.PartnerId = *partnerId
		conf.GcpProjectId = *gcpProjectId
	} else {
		file, err := os.Open(*configFile)
		if err != nil {
			fmt.Printf("Error reading confile file %s %s", *configFile, err)
			return conf, err
		}

		err = json.NewDecoder(file).Decode(&conf)
		if err != nil {
			fmt.Println("Configuration file not found. Continuing with default values.")
			return conf, err
		}
		fmt.Printf("Using confile file %s to launch subscription service \n", *configFile)
	}

	valid := true

	if conf.PubSubSubscription == "" {
		fmt.Println("PubSubSubscription was not set.")
		valid = false
	}

	if conf.PubSubTopicPrefix == "" {
		fmt.Println("PubSubTopicPrefix was not set.")
		valid = false
	}

	if conf.CloudCommerceProcurementUrl == "" {
		fmt.Println("CloudCommerceProcurementUrl was not set.")
		valid = false
	}

	if conf.PartnerId == "" {
		fmt.Println("PartnerId was not set.")
		valid = false
	}

	if conf.GcpProjectId == "" {
		fmt.Println("GcpProjectId was not set.")
		valid = false
	}

	credPath,envExists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !envExists {
		fmt.Println("GOOGLE_APPLICATION_CREDENTIALS was not set or path does not exist. This is fine with an emulator but will fail in production. ")
	}

	_, errPath := os.Stat(credPath)
	if os.IsNotExist(errPath) {
		fmt.Println("GOOGLE_APPLICATION_CREDENTIALS file does not exist: %s.",credPath)
		valid = false
	}

	if !valid {
		err = errors.New("Subscription service configuration is not valid!")
	}

	return conf, err
}
