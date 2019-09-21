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
	PartnerId							= "000"
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
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Println("Unable to determine working directory.")
		return conf, err
	}
	log.Printf("Running subscription service with working directory %s \n",dir)

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
		log.Printf("Using confile file %s to launch subscription frontend service \n", *configFile)
	}

	valid := true

	if conf.FrontendServiceEndpoint == "" {
		log.Println("FrontendServiceEndpoint was not set.")
		valid = false
	}

	if conf.SubscriptionServiceUrl == "" {
		log.Println("Subscription Service URL was not set.")
		valid = false
	} else if strings.HasSuffix(conf.SubscriptionServiceUrl,"/"){
		conf.SubscriptionServiceUrl = conf.SubscriptionServiceUrl[:len(conf.SubscriptionServiceUrl)-1]
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
	} else if strings.HasSuffix(conf.CloudCommerceProcurementUrl,"/"){
		conf.CloudCommerceProcurementUrl = conf.CloudCommerceProcurementUrl[:len(conf.CloudCommerceProcurementUrl)-1]
	}

	if conf.PartnerId == "" {
		log.Println("PartnerId was not set.")
		valid = false
	}

	if !valid {
		err = errors.New("Subscription frontend service configuration is not valid!")
	}

	return conf, err
}
