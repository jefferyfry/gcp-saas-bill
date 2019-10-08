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
	FrontendServiceEndpoint = ":8086"
	SubscriptionServiceUrl = "https://subscription-service.cloudbees-jenkins-support.svc.cluster.local"
	ClientId 		= "123456"
	ClientSecret    = "abcdef"
	CallbackUrl		= "http://localhost/callback"
	Issuer			= "http://localhost"
	SessionKey		= "cloudbeesjenkinssupportsessionkey1cl0udb33s1"
	CloudCommerceProcurementUrl       	= "https://cloudcommerceprocurement.googleapis.com/"
	PartnerId							= ""
	FinishUrl							= ""
	FinishUrlTitle						= ""
	TestMode							= "false"
)

type ServiceConfig struct {
	FrontendServiceEndpoint string `json:"frontendServiceEndpoint"`
	SubscriptionServiceUrl 	string `json:"subscriptionServiceUrl"`
	ClientId    			string	`json:"clientId"`
	ClientSecret    		string	`json:"clientSecret"`
	CallbackUrl    			string	`json:"callbackUrl"`
	Issuer    				string	`json:"issuer"`
	SessionKey    			string	`json:"sessionKey"`
	CloudCommerceProcurementUrl    	string	`json:"cloudCommerceProcurementUrl"`
	PartnerId    					string	`json:"partnerId"`
	FinishUrl    					string	`json:"finishUrl"`
	FinishUrlTitle    				string	`json:"finishUrlTitle"`
	TestMode    					string	`json:"testMode"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		FrontendServiceEndpoint,
		SubscriptionServiceUrl,
		ClientId,
		ClientSecret,
		CallbackUrl,
		Issuer,
		SessionKey,
		CloudCommerceProcurementUrl,
		PartnerId,
		FinishUrl,
		FinishUrlTitle,
		TestMode,
	}

	if dir, err := os.Getwd(); err != nil {
		log.Println("Unable to determine working directory.")
		return conf, err
	} else {
		log.Printf("Running service with working directory %s \n", dir)
	}

	//parse commandline arguments
	configFile := flag.String("configFile", "", "set the path to the configuration json file")
	frontendServiceEndpoint := flag.String("frontendServiceEndpoint", "", "set the value of the frontend service endpoint port")
	subscriptionServiceUrl := flag.String("subscriptionServiceUrl", "", "set the value of subscription service url")
	clientId := flag.String("clientId", "", "set the value of the Auth0 client ID")
	clientSecret := flag.String("clientSecret", "", "set the value of the Auth0 client secret")
	callbackUrl := flag.String("callbackUrl", "", "set the value for the Auth0 callback URL")
	issuer := flag.String("issuer", "", "set the value of the Auth0 issuer")
	sessionKey := flag.String("sessionKey", "", "set the value of http session key")
	cloudCommerceProcurementUrl := flag.String("cloudCommerceProcurementUrl", "", "set root url for the cloud commerce procurement API")
	partnerId := flag.String("partnerId", "", "set the CloudBees Partner Id")
	finishUrl := flag.String("finishUrl", "", "set the finish url")
	finishUrlTitle := flag.String("finishUrlTitle", "", "set the finish url title")
	testMode := flag.String("testMode", "", "set whether this runs in test mode")

	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CLOUD_BILL_FRONTEND_CONFIG_FILE")
	}
	if *frontendServiceEndpoint == "" {
		*frontendServiceEndpoint = os.Getenv("CLOUD_BILL_FRONTEND_SERVICE_ENDPOINT")
	}
	if *subscriptionServiceUrl == "" {
		*subscriptionServiceUrl = os.Getenv("CLOUD_BILL_SUBSCRIPTION_SERVICE_URL")
	}
	if *clientId == "" {
		*clientId = os.Getenv("CLOUD_BILL_FRONTEND_CLIENT_ID")
	}
	if *clientSecret == "" {
		*clientSecret = os.Getenv("CLOUD_BILL_FRONTEND_CLIENT_SECRET")
	}
	if *callbackUrl == "" {
		*callbackUrl = os.Getenv("CLOUD_BILL_FRONTEND_CALLBACK_URL")
	}
	if *issuer == "" {
		*issuer = os.Getenv("CLOUD_BILL_FRONTEND_ISSUER")
	}
	if *sessionKey == "" {
		*sessionKey = os.Getenv("CLOUD_BILL_FRONTEND_SESSION_KEY")
	}
	if *cloudCommerceProcurementUrl == "" {
		*cloudCommerceProcurementUrl = os.Getenv("CLOUD_BILL_FRONTEND_CLOUD_COMMERCE_PROCUREMENT_URL")
	}
	if *partnerId == "" {
		*partnerId = os.Getenv("CLOUD_BILL_FRONTEND_PARTNER_ID")
	}
	if *finishUrl == "" {
		*finishUrl = os.Getenv("CLOUD_BILL_FRONTEND_FINISH_URL")
	}
	if *finishUrlTitle == "" {
		*finishUrlTitle = os.Getenv("CLOUD_BILL_FRONTEND_FINISH_URL_TITLE")
	}
	if *testMode == "" {
		*testMode = os.Getenv("CLOUD_BILL_FRONTEND_TEST_MODE")
	}


	if *configFile == "" {
		//try other flags
		conf.FrontendServiceEndpoint = *frontendServiceEndpoint
		conf.SubscriptionServiceUrl = *subscriptionServiceUrl
		conf.ClientId = *clientId
		conf.ClientSecret = *clientSecret
		conf.CallbackUrl = *callbackUrl
		conf.Issuer = *issuer
		conf.SessionKey = *sessionKey
		conf.CloudCommerceProcurementUrl = *cloudCommerceProcurementUrl
		conf.PartnerId = *partnerId
		conf.FinishUrl = *finishUrl
		conf.FinishUrlTitle = *finishUrlTitle
		conf.TestMode = *testMode
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

	if conf.FrontendServiceEndpoint == "" {
		log.Println("FrontendServiceEndpoint was not set.")
		valid = false
	}

	if conf.SubscriptionServiceUrl == "" {
		log.Println("Subscription Service URL was not set.")
		valid = false
	} else {
		conf.SubscriptionServiceUrl = strings.TrimSuffix(conf.SubscriptionServiceUrl,"/")
	}

	if conf.ClientId == "" {
		log.Println("Client ID was not set.")
		valid = false
	}

	if conf.ClientSecret == "" {
		log.Println("ClientSecret was not set.")
		valid = false
	}

	if conf.CallbackUrl == "" {
		log.Println("Callback URL was not set.")
		valid = false
	}

	if conf.Issuer == "" {
		log.Println("Issuer was not set.")
		valid = false
	}

	if conf.SessionKey == "" {
		log.Println("SessionKey was not set.")
		valid = false
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

	if conf.FinishUrl == "" {
		log.Println("FinishUrl was not set.")
		valid = false
	}

	if conf.FinishUrlTitle == "" {
		log.Println("FinishUrlTitle was not set.")
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
		return conf, errors.New("Subscription frontend service configuration is not valid!")
	} else {
		return conf, nil
	}
}
