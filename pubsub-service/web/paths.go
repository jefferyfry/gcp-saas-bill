package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

//SetUpService sets up the subscription service.
func SetUpService(healthCheckEndpoint string, pubSubSubscription string, subscriptionServiceUrl string, cloudCommerceProcurementUrl string, partnerId string, gcpProjectId string) error {
	handler := GetPubSubServiceHandler(pubSubSubscription, subscriptionServiceUrl, cloudCommerceProcurementUrl, partnerId, gcpProjectId)
	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/healthz").HandlerFunc(handler.Healthz)

	return http.ListenAndServe(":"+healthCheckEndpoint, r)
}