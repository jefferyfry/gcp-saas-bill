package config

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/jefferyfry/funclog"
	"os"
	"strings"
)

var (
	HealthCheckEndpoint 				= "8097"
	PubSubSubscription       			= "codelab"
	SubscriptionServiceUrl 				= "https://subscription-service.cloudbees-jenkins-support.svc.cluster.local"
	CloudCommerceProcurementUrl       	= "https://cloudcommerceprocurement.googleapis.com/"
	PartnerId							= "000"
	GcpProjectId				        = "cloud-billing-saas"
	SentryDsn							= ""

	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

type ServiceConfig struct {
	HealthCheckEndpoint 			string `json:"healthCheckEndpoint"`
	PubSubSubscription    			string	`json:"pubSubSubscription"`
	SubscriptionServiceUrl 			string `json:"subscriptionServiceUrl"`
	CloudCommerceProcurementUrl    	string	`json:"cloudCommerceProcurementUrl"`
	PartnerId    					string	`json:"partnerId"`
	GcpProjectId    				string	`json:"gcpProjectId"`
	SentryDsn						string	`json:"sentryDsn"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		HealthCheckEndpoint,
		PubSubSubscription,
		SubscriptionServiceUrl,
		CloudCommerceProcurementUrl,
		PartnerId,
		GcpProjectId,
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
	healthCheckEndpoint := flag.String("healthCheckEndpoint", "", "set the value of the health check endpoint port")
	pubSubSubscription := flag.String("pubSubSubscription", "", "set the value of the pubsub subscription")
	subscriptionServiceUrl := flag.String("subscriptionServiceUrl", "", "set the value of subscription service url")
	cloudCommerceProcurementUrl := flag.String("cloudCommerceProcurementUrl", "", "set root url for the cloud commerce procurement API")
	partnerId := flag.String("partnerId", "", "set the CloudBees Partner Id")
	gcpProjectId := flag.String("gcpProjectId", "", "set the GCP Project Id")
	sentryDsn := flag.String("sentryDsn", "", "set the Sentry DSN")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CLOUD_BILL_PUBSUB_CONFIG_FILE")
	}
	if *healthCheckEndpoint == "" {
		*healthCheckEndpoint = os.Getenv("CLOUD_BILL_PUBSUB_HEALTH_CHECK_ENDPOINT")
	}
	if *pubSubSubscription == "" {
		*pubSubSubscription = os.Getenv("CLOUD_BILL_PUBSUB_SUBSCRIPTION")
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
	if *sentryDsn == "" {
		*sentryDsn = os.Getenv("CLOUD_BILL_PUBSUB_SENTRY_DSN")
	}

	if *configFile == "" {
		//try other flags
		conf.HealthCheckEndpoint = *healthCheckEndpoint
		conf.PubSubSubscription = *pubSubSubscription
		conf.SubscriptionServiceUrl = *subscriptionServiceUrl
		conf.CloudCommerceProcurementUrl = *cloudCommerceProcurementUrl
		conf.PartnerId = *partnerId
		conf.GcpProjectId = *gcpProjectId
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

	if conf.HealthCheckEndpoint == "" {
		LogE.Println("HealthCheckEndpoint was not set.")
		valid = false
	}

	if conf.PubSubSubscription == "" {
		LogE.Println("PubSubSubscription was not set.")
		valid = false
	}

	if conf.SubscriptionServiceUrl == "" {
		LogE.Println("Subscription Service URL was not set.")
		valid = false
	} else {
		conf.SubscriptionServiceUrl = strings.TrimSuffix(conf.SubscriptionServiceUrl,"/")
	}

	if conf.CloudCommerceProcurementUrl == "" {
		LogE.Println("CloudCommerceProcurementUrl was not set.")
		valid = false
	} else {
		conf.CloudCommerceProcurementUrl = strings.TrimSuffix(conf.CloudCommerceProcurementUrl,"/")
	}

	if conf.PartnerId == "" {
		LogE.Println("PartnerId was not set.")
		valid = false
	}

	if conf.GcpProjectId == "" {
		LogE.Println("GcpProjectId was not set.")
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
		return conf, errors.New("Subscription service configuration is not valid!")
	} else {
		return conf, nil
	}
}
