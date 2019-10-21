package web

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

//SetUpService sets up the subscription service.
func SetUpService(webServiceEndpoint string,healthCheckEndpoint string,subscriptionServiceUrl string,clientId string, clientSecret string, callbackUrl string, issuer string, sessionKey string, cloudCommerceProcurementUrl string, partnerId string, finishUrl string, finishUrlTitle string, testMode string) error {
	handler := GetSubscriptionFrontendHandler(subscriptionServiceUrl,clientId, clientSecret, callbackUrl, issuer, sessionKey, cloudCommerceProcurementUrl, partnerId, finishUrl, finishUrlTitle)

	healthCheck := mux.NewRouter()
	healthCheck.Methods(http.MethodGet).Path("/healthz").HandlerFunc(handler.Healthz)
	go http.ListenAndServe(":"+healthCheckEndpoint, healthCheck)

	webService := mux.NewRouter()
	if testModeBool,err := strconv.ParseBool(testMode); err==nil && testModeBool {
		webService.Methods(http.MethodGet).Path("/resetsaas").HandlerFunc(handler.ResetSaas)
		webService.Methods(http.MethodGet).Path("/signupsaastest").HandlerFunc(handler.SignupSaasTest)
	}
	webService.Methods(http.MethodGet).Path("/signupprod/{accountId}").HandlerFunc(handler.SignupProd)
	webService.Methods(http.MethodPost).Path("/signupsaas").HandlerFunc(handler.SignupSaas)
	webService.Methods(http.MethodGet).Path("/login").HandlerFunc(handler.Auth0Login)
	webService.Methods(http.MethodGet).Path("/callback").HandlerFunc(handler.Auth0Callback)
	webService.Methods(http.MethodPost).Path("/finishSaas").HandlerFunc(handler.FinishSaas)
	webService.Methods(http.MethodPost).Path("/finishProd").HandlerFunc(handler.FinishProd)

	webService.Methods(http.MethodGet).Path("/healthz").HandlerFunc(handler.Healthz)

	webService.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.cloudbees.com", http.StatusFound)
	})

	return http.ListenAndServe(":"+webServiceEndpoint, webService)
}