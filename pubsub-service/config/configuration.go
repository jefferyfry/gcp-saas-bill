package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	FrontendServiceEndpoint = ":8086"
	SubscriptionServiceUrl = "https://subscription-service.cloudbees-jenkins-support.svc.cluster.local"
	ClientId 		= "123456"
	ClientSecret    = "abcdef"
	CallbackUrl		= "http://localhost/callback"
	Issuer			= "http://localhost"
	SessionKey		= "cloudbeesjenkinssupportsessionkey1cl0udb33s1"
)

type ServiceConfig struct {
	FrontendServiceEndpoint string `json:"frontendServiceEndpoint"`
	SubscriptionServiceUrl string `json:"subscriptionServiceUrl"`
	ClientId    	string	`json:"clientId"`
	ClientSecret    string	`json:"clientSecret"`
	CallbackUrl    	string	`json:"callbackUrl"`
	Issuer    		string	`json:"issuer"`
	SessionKey    	string	`json:"sessionKey"`
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
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to determine working directory.")
		return conf, err
	}
	fmt.Printf("Running subscription service with working directory %s \n",dir)

	//parse commandline arguments
	configFile := flag.String("configFile", "", "set the path to the configuration json file")
	frontendServiceEndpoint := flag.String("frontendServiceEndpoint", "", "set the value of the frontend service endpoint port")
	subscriptionServiceUrl := flag.String("subscriptionServiceUrl", "", "set the value of subscription service url")
	clientId := flag.String("clientId", "", "set the value of the Auth0 client ID")
	clientSecret := flag.String("clientSecret", "", "set the value of the Auth0 client secret")
	callbackUrl := flag.String("callbackUrl", "", "set the value for the Auth0 callback URL")
	issuer := flag.String("issuer", "", "set the value of the Auth0 issuer")
	sessionKey := flag.String("sessionKey", "", "set the value of http session key")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("JENKINS_SUPPORT_SUB_FRONTEND_CONFIG_FILE")
	}
	if *frontendServiceEndpoint == "" {
		*frontendServiceEndpoint = os.Getenv("JENKINS_SUPPORT_SUB_FRONTEND_SERVICE_ENDPOINT")
	}
	if *subscriptionServiceUrl == "" {
		*subscriptionServiceUrl = os.Getenv("JENKINS_SUPPORT_SUB_SERVICE_URL")
	}
	if *clientId == "" {
		*clientId = os.Getenv("JENKINS_SUPPORT_SUB_FRONTEND_CLIENT_ID")
	}
	if *clientSecret == "" {
		*clientSecret = os.Getenv("JENKINS_SUPPORT_SUB_FRONTEND_CLIENT_SECRET")
	}
	if *callbackUrl == "" {
		*callbackUrl = os.Getenv("JENKINS_SUPPORT_SUB_FRONTEND_CALLBACK_URL")
	}
	if *issuer == "" {
		*issuer = os.Getenv("JENKINS_SUPPORT_SUB_FRONTEND_ISSUER")
	}
	if *sessionKey == "" {
		*sessionKey = os.Getenv("JENKINS_SUPPORT_SUB_FRONTEND_SESSION_KEY")
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
		fmt.Printf("Using confile file %s to launch subscription frontend service \n", *configFile)
	}

	valid := true

	if conf.FrontendServiceEndpoint == "" {
		fmt.Println("FrontendServiceEndpoint was not set.")
		valid = false
	}

	if conf.SubscriptionServiceUrl == "" {
		fmt.Println("Subscription Service URL was not set.")
		valid = false
	}

	if conf.ClientId == "" {
		fmt.Println("Client ID was not set.")
		valid = false
	}

	if conf.ClientSecret == "" {
		fmt.Println("ClientSecret was not set.")
		valid = false
	}

	if conf.CallbackUrl == "" {
		fmt.Println("Callback URL was not set.")
		valid = false
	}

	if conf.Issuer == "" {
		fmt.Println("Issuer was not set.")
		valid = false
	}

	if conf.SessionKey == "" {
		fmt.Println("SessionKey was not set.")
		valid = false
	}

	if !valid {
		err = errors.New("Subscription frontend service configuration is not valid!")
	}

	return conf, err
}
