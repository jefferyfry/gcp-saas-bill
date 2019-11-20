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
	FrontendServiceEndpoint = "8086"
	HealthCheckEndpoint = "8096"
	SubscriptionServiceUrl = "https://subscription-service.cloudbees-jenkins-support.svc.cluster.local"
	GoogleSubscriptionsUrl = "https://cloudbilling.googleapis.com/v1"
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
	SentryDsn							= ""
	GcpProjectId				        = "cloud-billing-saas"
	
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

type ServiceConfig struct {
	FrontendServiceEndpoint string `json:"frontendServiceEndpoint"`
	HealthCheckEndpoint string `json:"healthCheckEndpoint"`
	SubscriptionServiceUrl 	string `json:"subscriptionServiceUrl"`
	GoogleSubscriptionsUrl 	string `json:"googleSubscriptionsUrl"`
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
	SentryDsn						string	`json:"sentryDsn"`
	GcpProjectId    				string	`json:"gcpProjectId"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		FrontendServiceEndpoint,
		HealthCheckEndpoint,
		SubscriptionServiceUrl,
		GoogleSubscriptionsUrl,
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
		SentryDsn,
		GcpProjectId,
	}

	if dir, err := os.Getwd(); err != nil {
		LogE.Println("Unable to determine working directory.")
		return conf, err
	} else {
		LogI.Printf("Running service with working directory %s \n", dir)
	}

	//parse commandline arguments
	configFile := flag.String("configFile", "", "set the path to the configuration json file")
	frontendServiceEndpoint := flag.String("frontendServiceEndpoint", "", "set the value of the frontend service endpoint port")
	healthCheckEndpoint := flag.String("healthCheckEndpoint", "", "set the value of the health check endpoint port")
	subscriptionServiceUrl := flag.String("subscriptionServiceUrl", "", "set the value of subscription service url")
	googleSubscriptionsUrl := flag.String("googleSubscriptionsUrl", "", "set the Google subscription url")
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
	sentryDsn := flag.String("sentryDsn", "", "set the Sentry DSN")
	gcpProjectId := flag.String("gcpProjectId", "", "set the GCP Project Id")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CLOUD_BILL_FRONTEND_CONFIG_FILE")
	}
	if *frontendServiceEndpoint == "" {
		*frontendServiceEndpoint = os.Getenv("CLOUD_BILL_FRONTEND_SERVICE_ENDPOINT")
	}
	if *healthCheckEndpoint == "" {
		*healthCheckEndpoint = os.Getenv("CLOUD_BILL_FRONTEND_HEALTH_CHECK_ENDPOINT")
	}
	if *subscriptionServiceUrl == "" {
		*subscriptionServiceUrl = os.Getenv("CLOUD_BILL_FRONTEND_SUBSCRIPTION_SERVICE_URL")
	}
	if *googleSubscriptionsUrl == "" {
		*googleSubscriptionsUrl = os.Getenv("CLOUD_BILL_FRONTEND_GOOGLE_SUBSCRIPTIONS_URL")
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
	if *gcpProjectId == "" {
		*gcpProjectId = os.Getenv("CLOUD_BILL_FRONTEND_GCP_PROJECT_ID")
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
	if *sentryDsn == "" {
		*sentryDsn = os.Getenv("CLOUD_BILL_FRONTEND_SENTRY_DSN")
	}

	if *configFile == "" {
		//try other flags
		conf.FrontendServiceEndpoint = *frontendServiceEndpoint
		conf.HealthCheckEndpoint = *healthCheckEndpoint
		conf.SubscriptionServiceUrl = *subscriptionServiceUrl
		conf.GoogleSubscriptionsUrl = *googleSubscriptionsUrl
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
		conf.SentryDsn = *sentryDsn
		conf.GcpProjectId = *gcpProjectId
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

	if conf.FrontendServiceEndpoint == "" {
		LogE.Println("FrontendServiceEndpoint was not set.")
		valid = false
	}

	if conf.HealthCheckEndpoint == "" {
		LogE.Println("HealthCheckEndpoint was not set.")
		valid = false
	}

	if conf.SubscriptionServiceUrl == "" {
		LogE.Println("Subscription Service URL was not set.")
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

	if conf.ClientId == "" {
		LogE.Println("Client ID was not set.")
		valid = false
	}

	if conf.ClientSecret == "" {
		LogE.Println("ClientSecret was not set.")
		valid = false
	}

	if conf.CallbackUrl == "" {
		LogE.Println("Callback URL was not set.")
		valid = false
	}

	if conf.Issuer == "" {
		LogE.Println("Issuer was not set.")
		valid = false
	}

	if conf.SessionKey == "" {
		LogE.Println("SessionKey was not set.")
		valid = false
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

	if conf.FinishUrl == "" {
		LogE.Println("FinishUrl was not set.")
		valid = false
	}

	if conf.FinishUrlTitle == "" {
		LogE.Println("FinishUrlTitle was not set.")
		valid = false
	}

	if conf.TestMode == "" {
		LogE.Println("TestMode was not set. Setting to false.")
		conf.TestMode = "false"
	}

	if conf.SentryDsn == "" {
		LogE.Println("SentryDsn was not set. Will run without Sentry.")
	}

	if conf.GcpProjectId == "" {
		LogE.Println("GcpProjectId was not set.")
		valid = false
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
		return conf, errors.New("Subscription frontend service configuration is not valid!")
	} else {
		return conf, nil
	}
}
