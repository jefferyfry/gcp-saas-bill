package rest

import (
	_ "github.com/cloudbees/jenkins-support-saas/subscription-service/docs"
	"github.com/cloudbees/jenkins-support-saas/subscription-service/persistence"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"net/http"
)

//SetUpService sets up the subscription service.
func SetUpService(dbHandler persistence.DatabaseHandler,serviceEndpoint string, cloudCommerceProcurementUrl string, partnerId string) error {
	handler := GetSubscriptionServiceHandler(dbHandler,cloudCommerceProcurementUrl,partnerId)
	r := mux.NewRouter()
	subscriptionRouter := r.PathPrefix("/api/v1").Subrouter()

	//accounts
	subscriptionRouter.Methods(http.MethodPost).Path("/accounts").HandlerFunc(handler.UpsertAccount)
	subscriptionRouter.Methods(http.MethodGet).Path("/accounts/{accountName}").HandlerFunc(handler.GetAccount)
	subscriptionRouter.Methods(http.MethodPut).Path("/accounts").HandlerFunc(handler.UpsertAccount)
	subscriptionRouter.Methods(http.MethodDelete).Path("/accounts/{accountName}").HandlerFunc(handler.DeleteAccount)

	//entitlements
	subscriptionRouter.Methods(http.MethodPost).Path("/entitlements").HandlerFunc(handler.UpsertEntitlement)
	subscriptionRouter.Methods(http.MethodGet).Path("/entitlements/{entitlementName}").HandlerFunc(handler.GetEntitlement)
	subscriptionRouter.Methods(http.MethodPut).Path("/entitlements").HandlerFunc(handler.UpsertEntitlement)
	subscriptionRouter.Methods(http.MethodDelete).Path("/entitlements/{entitlementName}").HandlerFunc(handler.DeleteEntitlement)

	//swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8085/swagger/doc.json"), //The url pointing to API definition"
	))

	return http.ListenAndServe(serviceEndpoint, r)
}