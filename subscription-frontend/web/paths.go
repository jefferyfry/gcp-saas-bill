package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

//SetUpService sets up the subscription service.
func SetUpService(serviceEndpoint string) error {
	handler := GetSubscriptionFrontendHandler()
	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/signup").HandlerFunc(handler.Signup)
	r.Methods(http.MethodPost).Path("/login").HandlerFunc(handler.Auth0Login) //this shouldn't be called unless istio fails. this redirects to auth0
	r.Methods(http.MethodPost).Path("/callback").HandlerFunc(handler.Auth0Callback)
	r.Methods(http.MethodPost).Path("/confirm").HandlerFunc(handler.Confirm)

	return http.ListenAndServe(serviceEndpoint, r)
}