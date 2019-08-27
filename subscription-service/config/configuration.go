package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	SubscriptionServiceEndpoint 		= ":8085"
	CloudCommerceProcurementUrl       	= "https://cloudcommerceprocurement.googleapis.com/"
	PartnerId							= "000"
	GoogleProjectId				        = "jenkins-support-saas"
)

type ServiceConfig struct {
	SubscriptionServiceEndpoint    	string	`json:"subscriptionServiceEndpoint"`
	CloudCommerceProcurementUrl    	string	`json:"cloudCommerceProcurementUrl"`
	PartnerId    					string	`json:"partnerId"`
	GoogleProjectId    				string	`json:"googleProjectId"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		SubscriptionServiceEndpoint,
		CloudCommerceProcurementUrl,
		PartnerId,
		GoogleProjectId,
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to determine working directory.")
		return conf, err
	}
	fmt.Printf("Running subscription service with working directory %s \n",dir)

	//parse commandline arguments
	configFile := flag.String("configFile", "", "set the path to the configuration json file")
	subscriptionServiceEndpoint := flag.String("subscriptionServiceEndpoint", "", "set the value of this service endpoint")
	cloudCommerceProcurementUrl := flag.String("cloudCommerceProcurementUrl", "", "set root url for the cloud commerce procurement API")
	partnerId := flag.String("partnerId", "", "set the CloudBees Partner Id")
	googleProjectId := flag.String("googleProjectId", "", "set the Google Project Id")
	flag.Parse()

	//try environment variables if necessary
	if configFile == nil {
		*configFile = os.Getenv("JENKINS_SUPPORT_SAAS_CONFIG_FILE")
	}
	if subscriptionServiceEndpoint == nil {
		*subscriptionServiceEndpoint = os.Getenv("JENKINS_SUPPORT_SAAS_SUBSCRIPTION_SERVICE_ENDPOINT")
	}
	if cloudCommerceProcurementUrl == nil {
		*cloudCommerceProcurementUrl = os.Getenv("JENKINS_SUPPORT_SAAS_CLOUD_COMMERCE_PROCUREMENT_URL")
	}
	if partnerId == nil {
		*partnerId = os.Getenv("JENKINS_SUPPORT_SAAS_PARTNER_ID")
	}
	if googleProjectId == nil {
		*googleProjectId = os.Getenv("JENKINS_SUPPORT_SAAS_GOOGLE_PROJECT_ID")
	}


	if *configFile == "" {
		//try other flags
		conf.SubscriptionServiceEndpoint = *subscriptionServiceEndpoint
		conf.CloudCommerceProcurementUrl = *cloudCommerceProcurementUrl
		conf.PartnerId = *partnerId
		conf.GoogleProjectId = *googleProjectId
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

	if conf.SubscriptionServiceEndpoint == "" {
		fmt.Println("SubscriptionServiceEndpoint was not set.")
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

	_,pathExists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !pathExists {
		fmt.Println("GOOGLE_APPLICATION_CREDENTIALS was not set.")
		valid = false
	}

	if conf.GoogleProjectId == "" {
		fmt.Println("GoogleProjectId was not set.")
		valid = false
	}

	if !valid {
		err = errors.New("Subscription service configuration is not valid!")
	}

	return conf, err
}
