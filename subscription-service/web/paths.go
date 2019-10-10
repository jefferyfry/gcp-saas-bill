package web

import (
	_ "github.com/cloudbees/cloud-bill-saas/subscription-service/docs"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/persistence"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"net/http"
)

//SetUpService sets up the subscription service.
func SetUpService(dbHandler persistence.DatabaseHandler,webServiceEndpoint string, healthCheckEndpoint string) error {
	handler := GetSubscriptionServiceHandler(dbHandler)
	healthCheck := mux.NewRouter()
	healthCheck.Methods(http.MethodGet).Path("/healthz").HandlerFunc(handler.Healthz)
	go http.ListenAndServe(":"+healthCheckEndpoint, healthCheck)

	webService := mux.NewRouter()
	apiV1 := webService.PathPrefix("/api/v1").Subrouter()

	//accounts
	apiV1.Methods(http.MethodPost).Path("/accounts").HandlerFunc(handler.UpsertAccount)
	apiV1.Methods(http.MethodGet).Path("/accounts/{accountId}").HandlerFunc(handler.GetAccount)
	apiV1.Methods(http.MethodPut).Path("/accounts").HandlerFunc(handler.UpsertAccount)
	apiV1.Methods(http.MethodDelete).Path("/accounts/{accountId}").HandlerFunc(handler.DeleteAccount)
	apiV1.Methods(http.MethodGet).Path("/accounts").HandlerFunc(handler.GetAccounts)

	//contacts
	apiV1.Methods(http.MethodPost).Path("/contacts").HandlerFunc(handler.UpsertContact)
	apiV1.Methods(http.MethodGet).Path("/contacts/{accountId}").HandlerFunc(handler.GetContact)
	apiV1.Methods(http.MethodPut).Path("/contacts").HandlerFunc(handler.UpsertContact)
	apiV1.Methods(http.MethodDelete).Path("/contacts/{accountId}").HandlerFunc(handler.DeleteContact)
	apiV1.Methods(http.MethodGet).Path("/contacts").HandlerFunc(handler.GetContacts)

	//entitlements
	apiV1.Methods(http.MethodPost).Path("/entitlements").HandlerFunc(handler.UpsertEntitlement)
	apiV1.Methods(http.MethodGet).Path("/entitlements/{entitlementId}").HandlerFunc(handler.GetEntitlement)
	apiV1.Methods(http.MethodPut).Path("/entitlements").HandlerFunc(handler.UpsertEntitlement)
	apiV1.Methods(http.MethodDelete).Path("/entitlements/{entitlementId}").HandlerFunc(handler.DeleteEntitlement)
	apiV1.Methods(http.MethodGet).Path("/entitlements").HandlerFunc(handler.GetEntitlements)
	apiV1.Methods(http.MethodGet).Path("/accounts/{accountId}/entitlements").HandlerFunc(handler.GetAccountEntitlements)

	apiV1.Methods(http.MethodGet).Path("/healthz").HandlerFunc(handler.Healthz)

	//swagger
	webService.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+webServiceEndpoint+"/swagger/doc.json"), //The url pointing to API definition"
	))

	return http.ListenAndServe(":"+webServiceEndpoint, webService)
}