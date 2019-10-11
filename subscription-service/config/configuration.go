package config

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/jefferyfry/funclog"
	"os"
)

var (
	SubscriptionServiceEndpoint 		= "8085"
	HealthCheckEndpoint 				= "8095"
	GcpProjectId				        = "cloud-bill-saas"
	SentryDsn							= ""

	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

type ServiceConfig struct {
	SubscriptionServiceEndpoint    	string	`json:"subscriptionServiceEndpoint"`
	HealthCheckEndpoint string `json:"healthCheckEndpoint"`
	GcpProjectId    				string	`json:"gcpProjectId"`
	SentryDsn						string	`json:"sentryDsn"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		SubscriptionServiceEndpoint,
		HealthCheckEndpoint,
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
	subscriptionServiceEndpoint := flag.String("subscriptionServiceEndpoint", "", "set the value of this service endpoint")
	healthCheckEndpoint := flag.String("healthCheckEndpoint", "", "set the value of the health check endpoint port")
	gcpProjectId := flag.String("gcpProjectId", "", "set the GCP Project Id")
	sentryDsn := flag.String("sentryDsn", "", "set the Sentry DSN")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CLOUD_BILL_SUBSCRIPTION_CONFIG_FILE")
	}
	if *subscriptionServiceEndpoint == "" {
		*subscriptionServiceEndpoint = os.Getenv("CLOUD_BILL_SUBSCRIPTION_SERVICE_ENDPOINT")
	}
	if *healthCheckEndpoint == "" {
		*healthCheckEndpoint = os.Getenv("CLOUD_BILL_SUBSCRIPTION_HEALTH_CHECK_ENDPOINT")
	}
	if *gcpProjectId == "" {
		*gcpProjectId = os.Getenv("CLOUD_BILL_SUBSCRIPTION_GCP_PROJECT_ID")
	}

	if *sentryDsn == "" {
		*sentryDsn = os.Getenv("CLOUD_BILL_SUBSCRIPTION_SENTRY_DSN")
	}


	if *configFile == "" {
		//try other flags
		conf.SubscriptionServiceEndpoint = *subscriptionServiceEndpoint
		conf.HealthCheckEndpoint = *healthCheckEndpoint
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

	if conf.SubscriptionServiceEndpoint == "" {
		LogE.Println("SubscriptionServiceEndpoint was not set.")
		valid = false
	}
	if conf.HealthCheckEndpoint == "" {
		LogE.Println("HealthCheckEndpoint was not set.")
		valid = false
	}

	if conf.GcpProjectId == "" {
		LogE.Println("GcpProjectId was not set.")
		valid = false
	}

	if conf.SentryDsn == "" {
		LogE.Println("SentryDsn was not set. Will run without Sentry.")
	}

	if credPath,envExists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); !envExists {
		LogE.Println("GOOGLE_APPLICATION_CREDENTIALS was not set. This is fine with an emulator but will fail in production. ")
	} else {
		if _, errPath := os.Stat(credPath); os.IsNotExist(errPath) {
			LogE.Println("GOOGLE_APPLICATION_CREDENTIALS file does not exist: ", credPath)
			valid = false
		} else {
			LogI.Println("Using GOOGLE_APPLICATION_CREDENTIALS file: ", credPath)
		}
	}

	if !valid {
		return conf, errors.New("Subscription service configuration is not valid!")
	} else {
		return conf, nil
	}
}
