package config

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"strings"
)

var (
	PubSubSubscription       			= "codelab"
	PubSubTopicPrefix       			= "DEMO-"
	SubscriptionServiceUrl 				= "https://subscription-service.cloudbees-jenkins-support.svc.cluster.local"
	CloudCommerceProcurementUrl       	= "https://cloudcommerceprocurement.googleapis.com/"
	PartnerId							= "000"
	GcpProjectId				        = "cloud-billing-saas"
)

type ServiceConfig struct {
	PubSubSubscription    			string	`json:"pubSubSubscription"`
	PubSubTopicPrefix				string 	`json:"pubSubTopicPrefix"`
	SubscriptionServiceUrl 			string `json:"subscriptionServiceUrl"`
	CloudCommerceProcurementUrl    	string	`json:"cloudCommerceProcurementUrl"`
	PartnerId    					string	`json:"partnerId"`
	GcpProjectId    				string	`json:"gcpProjectId"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		PubSubSubscription,
		PubSubTopicPrefix,
		SubscriptionServiceUrl,
		CloudCommerceProcurementUrl,
		PartnerId,
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
	pubSubSubscription := flag.String("pubSubSubscription", "", "set the value of the pubsub subscription")
	pubSubTopicPrefix := flag.String("pubSubTopicPrefix", "", "set the value of the pubsub topic prefix")
	subscriptionServiceUrl := flag.String("subscriptionServiceUrl", "", "set the value of subscription service url")
	cloudCommerceProcurementUrl := flag.String("cloudCommerceProcurementUrl", "", "set root url for the cloud commerce procurement API")
	partnerId := flag.String("partnerId", "", "set the CloudBees Partner Id")
	gcpProjectId := flag.String("gcpProjectId", "", "set the GCP Project Id")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CLOUD_BILL_PUBSUB_CONFIG_FILE")
	}
	if *pubSubSubscription == "" {
		*pubSubSubscription = os.Getenv("CLOUD_BILL_PUBSUB_SUBSCRIPTION")
	}
	if *pubSubTopicPrefix == "" {
		*pubSubTopicPrefix = os.Getenv("CLOUD_BILL_PUBSUB_TOPIC_PREFIX")
	}
	if *subscriptionServiceUrl == "" {
		*subscriptionServiceUrl = os.Getenv("CLOUD_BILL_SUBSCRIPTION_SERVICE_URL")
	}
	if *cloudCommerceProcurementUrl == "" {
		*cloudCommerceProcurementUrl = os.Getenv("CLOUD_BILL_PUBSUB_CLOUD_COMMERCE_PROCUREMENT_URL")
	}
	if *partnerId == "" {
		*partnerId = os.Getenv("CLOUD_BILL_PUBSUB_PARTNER_ID")
	}
	if *gcpProjectId == "" {
		*gcpProjectId = os.Getenv("CLOUD_BILL_PUBSUB_GCP_PROJECT_ID")
	}

	if *configFile == "" {
		//try other flags
		conf.PubSubSubscription = *pubSubSubscription
		conf.PubSubTopicPrefix = *pubSubTopicPrefix
		conf.SubscriptionServiceUrl = *subscriptionServiceUrl
		conf.CloudCommerceProcurementUrl = *cloudCommerceProcurementUrl
		conf.PartnerId = *partnerId
		conf.GcpProjectId = *gcpProjectId
	} else {
		if file, err := os.Open(*configFile); err != nil {
			log.Printf("Error reading confile file %s %s", *configFile, err)
			return conf, err
		} else {
			if err = json.NewDecoder(file).Decode(&conf); err != nil {
				return conf, errors.New("Configuration file not found.")
			}
			log.Printf("Using confile file %s to launch subscription frontend service \n", *configFile)
		}
	}

	valid := true

	if conf.PubSubSubscription == "" {
		log.Println("PubSubSubscription was not set.")
		valid = false
	}

	if conf.PubSubTopicPrefix == "" {
		log.Println("PubSubTopicPrefix was not set.")
		valid = false
	}

	if conf.SubscriptionServiceUrl == "" {
		log.Println("Subscription Service URL was not set.")
		valid = false
	} else {
		conf.SubscriptionServiceUrl = strings.TrimSuffix(conf.SubscriptionServiceUrl,"/")
	}

	if conf.CloudCommerceProcurementUrl == "" {
		log.Println("CloudCommerceProcurementUrl was not set.")
		valid = false
	} else {
		conf.CloudCommerceProcurementUrl = strings.TrimSuffix(conf.CloudCommerceProcurementUrl,"/")
	}

	if conf.PartnerId == "" {
		log.Println("PartnerId was not set.")
		valid = false
	}

	if conf.GcpProjectId == "" {
		log.Println("GcpProjectId was not set.")
		valid = false
	}

	if gAppCredPath,gAppCredExists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); !gAppCredExists {
		log.Println("GOOGLE_APPLICATION_CREDENTIALS was not set. ")
		valid = false
	} else {
		if _, gAppCredPathErr := os.Stat(gAppCredPath); os.IsNotExist(gAppCredPathErr) {
			log.Println("GOOGLE_APPLICATION_CREDENTIALS file does not exist: ", gAppCredPath)
			valid = false
		} else {
			log.Println("Using GOOGLE_APPLICATION_CREDENTIALS file: ", gAppCredPath)
		}
	}

	if !valid {
		return conf, errors.New("Subscription service configuration is not valid!")
	} else {
		return conf, nil
	}
}
