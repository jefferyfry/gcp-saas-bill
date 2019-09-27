package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

//SetUpService sets up the subscription service.
func SetUpService(serviceEndpoint string,subscriptionServiceUrl string,clientId string, clientSecret string, callbackUrl string, issuer string, sessionKey string, cloudCommerceProcurementUrl string, partnerId string) error {
	handler := GetSubscriptionFrontendHandler(subscriptionServiceUrl,clientId, clientSecret, callbackUrl, issuer, sessionKey, cloudCommerceProcurementUrl, partnerId)
	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/resetsaas").HandlerFunc(handler.ResetSaas)
	r.Methods(http.MethodGet).Path("/signupsaastest").HandlerFunc(handler.SignupSaasTest)
	r.Methods(http.MethodGet).Path("/signupprod/{accountId}").HandlerFunc(handler.SignupProd)
	r.Methods(http.MethodPost).Path("/signupsaas").HandlerFunc(handler.SignupSaas)
	r.Methods(http.MethodGet).Path("/login").HandlerFunc(handler.Auth0Login)
	r.Methods(http.MethodGet).Path("/").HandlerFunc(handler.EmailConfirm)
	r.Methods(http.MethodGet).Path("/callback").HandlerFunc(handler.Auth0Callback)
	r.Methods(http.MethodPost).Path("/finishSaas").HandlerFunc(handler.FinishSaas)
	r.Methods(http.MethodPost).Path("/finishProd").HandlerFunc(handler.FinishProd)

	return http.ListenAndServe(":"+serviceEndpoint, r)
}