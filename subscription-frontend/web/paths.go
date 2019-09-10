package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

//SetUpService sets up the subscription service.
func SetUpService(serviceEndpoint string,subscriptionServiceUrl string,clientId string, clientSecret string, callbackUrl string, issuer string, sessionKey string) error {
	handler := GetSubscriptionFrontendHandler(subscriptionServiceUrl,clientId, clientSecret, callbackUrl, issuer, sessionKey)
	r := mux.NewRouter()

	r.Methods(http.MethodPost).Path("/signuptest").HandlerFunc(handler.SignupTest)
	r.Methods(http.MethodPost).Path("/signup").HandlerFunc(handler.Signup)
	r.Methods(http.MethodGet).Path("/login").HandlerFunc(handler.Auth0Login) //this shouldn't be called unless istio fails. this redirects to auth0
	r.Methods(http.MethodPost).Path("/callback").HandlerFunc(handler.Auth0Callback)
	r.Methods(http.MethodPost).Path("/finish").HandlerFunc(handler.Finish)

	return http.ListenAndServe(":"+serviceEndpoint, r)
}