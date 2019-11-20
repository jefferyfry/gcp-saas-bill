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
	GcpProjectId	= "cloud-billing-saas"
	Products	= ""
	SubscriptionServiceUrl = "https://subscription-service.cloudbees-jenkins-support.svc.cluster.local"
	GoogleSubscriptionsUrl = "https://cloudbilling.googleapis.com/v1"
	SentryDsn		= ""

	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

type ServiceConfig struct {
	GcpProjectId    string	`json:"gcpProjectId"`
	Products    string	`json:"products"`
	SubscriptionServiceUrl 	string `json:"subscriptionServiceUrl"`
	GoogleSubscriptionsUrl 	string `json:"googleSubscriptionsUrl"`
	SentryDsn		string	`json:"sentryDsn"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		GcpProjectId,
		Products,
		SubscriptionServiceUrl,
		GoogleSubscriptionsUrl,
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
	products := flag.String("products", "", "a comma separated list of products to check")
	subscriptionServiceUrl := flag.String("subscriptionServiceUrl", "", "set the subscription service url")
	googleSubscriptionsUrl := flag.String("googleSubscriptionsUrl", "", "set the Google subscription url")
	sentryDsn := flag.String("sentryDsn", "", "set the Sentry DSN")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CLOUD_BILL_ENTITLEMENT_CHECK_CONFIG_FILE")
	}
	if *gcpProjectId == "" {
		*gcpProjectId = os.Getenv("CLOUD_BILL_ENTITLEMENT_CHECK_GCP_PROJECT_ID")
	}

	if *products == "" {
		*products = os.Getenv("CLOUD_BILL_ENTITLEMENT_CHECK_PRODUCTS")
	}

	if *subscriptionServiceUrl == "" {
		*subscriptionServiceUrl = os.Getenv("CLOUD_BILL_ENTITLEMENT_CHECK_SUBSCRIPTION_SERVICE_URL")
	}

	if *googleSubscriptionsUrl == "" {
		*googleSubscriptionsUrl = os.Getenv("CLOUD_BILL_ENTITLEMENT_CHECK_GOOGLE_SUBSCRIPTIONS_URL")
	}

	if *sentryDsn == "" {
		*sentryDsn = os.Getenv("CLOUD_BILL_ENTITLEMENT_CHECK_SENTRY_DSN")
	}

	if *configFile == "" {
		//try other flags
		conf.GcpProjectId = *gcpProjectId
		conf.Products = *products
		conf.SubscriptionServiceUrl = *subscriptionServiceUrl
		conf.GoogleSubscriptionsUrl = *googleSubscriptionsUrl
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

	if conf.Products == "" {
		LogE.Println("Products was not set.")
		valid = false
	} else {
		LogI.Printf("Running entitlement checks for products: %s",conf.Products)
	}

	if conf.SubscriptionServiceUrl == "" {
		LogE.Println("SubscriptionServiceUrl was not set.")
		valid = false
	} else {
		conf.SubscriptionServiceUrl = strings.TrimSuffix(conf.SubscriptionServiceUrl,"/")
	}

	if conf.GoogleSubscriptionsUrl == "" {
		LogE.Println("GoogleSubscriptionsUrl was not set.")
		valid = false
	} else {
		conf.GoogleSubscriptionsUrl = strings.TrimSuffix(conf.GoogleSubscriptionsUrl,"/")
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
		return conf, errors.New("Entitlement check configuration is not valid!")
	} else {
		return conf, nil
	}
}
