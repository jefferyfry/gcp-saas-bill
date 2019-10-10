package web

import (
	"encoding/json"
	"fmt"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/persistence"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type SubscriptionServiceHandler struct {
	dbHandler             persistence.DatabaseHandler
}

func GetSubscriptionServiceHandler(dbHandler persistence.DatabaseHandler) *SubscriptionServiceHandler {
	return &SubscriptionServiceHandler {
		dbHandler,
	}
}

// @Summary Get an account
// @Description Retrieves an account by account ID
// @ID cloud-bill-saas-subscription-service-get-account
// @Accept  json
// @Produce  json
// @Param accountId path string true "Account ID"
// @Success 200 {object} persistence.Account
// @Failure 400 {string} string "Missing account ID in path"
// @Failure 500 {string} string "Internal server error"
// @Router /accounts/{accountId} [get]
func (hdlr *SubscriptionServiceHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["accountId"]

	if accountId == "" {
		http.Error(w,`{"error": "missing account ID"}`,400)
		return
	}

	if account, dbErr := hdlr.dbHandler.GetAccount(accountId); nil != dbErr {
		if dbErr.Error() == "datastore: no such entity" {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "error occured while getting account %s", dbErr)
		}
	} else {
		if account == nil {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(&account)
		}
	}
}

// @Summary GetAccounts
// @Description Gets an array of accounts
// @ID cloud-bill-saas-subscription-service-get-accounts
// @Accept  json
// @Produce  json
// @Param filters query string false "optional comma separated list of filter"
// @Param order query string false "optional order"
// @Success 200 {array} persistence.Account
// @Failure 500 {string} string "Error"
// @Router /accounts [get]
func (hdlr *SubscriptionServiceHandler) GetAccounts(w http.ResponseWriter, r *http.Request){
	filtersParam, ok := r.URL.Query()["filters"]
	var filters []string = nil
	if ok || len(filtersParam) > 0 {
		filters = strings.Split(filtersParam[0],",")
	}

	ordersParam, ok := r.URL.Query()["order"]
	var order = ""
	if ok || len(ordersParam) > 0 {
		order = ordersParam[0]
	}

	if accounts, dbErr := hdlr.dbHandler.QueryAccounts(filters,order); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting account %s", dbErr)
	} else {
		if accounts == nil {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(&accounts)
		}
	}
}

// @Summary Upsert an account
// @Description Upsert an account passing account json
// @ID cloud-bill-saas-subscription-service-upsert-account
// @Accept  json
// @Produce  json
// @Param account body persistence.Account true "Account"
// @Success 204 {string} string "Upserted"
// @Failure 500 {string} string "Error"
// @Router /accounts [put]
func (hdlr *SubscriptionServiceHandler) UpsertAccount(w http.ResponseWriter, r *http.Request) {
	account := persistence.Account{}
	if err := json.NewDecoder(r.Body).Decode(&account); nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding account data %s", err)
		return
	}
	if dbErr := hdlr.dbHandler.UpsertAccount(&account); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while persisting account %s", dbErr)
	} else {
		w.WriteHeader(204)
	}
}

// @Summary Delete an account
// @Description Delete an account
// @ID cloud-bill-saas-subscription-service-delete-account
// @Accept  json
// @Produce  json
// @Param accountId path string true "Account ID"
// @Success 204 {string} string "Deleted"
// @Failure 400 {string} string "Missing account ID in path"
// @Failure 500 {string} string "Internal server error"
// @Router /accounts/{accountId} [delete]
func (hdlr *SubscriptionServiceHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["accountId"]

	if accountId == "" {
		http.Error(w,`{"error": "missing account ID"}`,400)
		return
	}

	if dbErr := hdlr.dbHandler.DeleteAccount(accountId); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while deleting account %s", dbErr)
	} else {
		w.WriteHeader(204)
	}
}

// @Summary Get an contact
// @Description Retrieves an contact by account ID
// @ID cloud-bill-saas-subscription-service-get-contact
// @Accept  json
// @Produce  json
// @Param accountId path string true "Account ID"
// @Success 200 {object} persistence.Contact
// @Failure 400 {string} string "Missing account ID in path"
// @Failure 500 {string} string "Error"
// @Router /contacts/{accountId} [get]
func (hdlr *SubscriptionServiceHandler) GetContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["accountId"]

	if accountId == "" {
		http.Error(w,`{"error": "missing account ID in path"}`,400)
		return
	}

	if contact, dbErr := hdlr.dbHandler.GetContact(accountId); nil != dbErr {
		if dbErr.Error() == "datastore: no such entity" {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "error occured while getting contact %s", dbErr)
		}
	} else {
		if contact == nil {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(&contact)
		}
	}
}

// @Summary Upsert a contact
// @Description Upsert a contact passing contact json
// @ID cloud-bill-saas-subscription-service-upsert-contact
// @Accept  json
// @Produce  json
// @Success 204 {string} string "Upserted"
// @Failure 500 {string} string "Error"
// @Router /contacts [put]
func (hdlr *SubscriptionServiceHandler) UpsertContact(w http.ResponseWriter, r *http.Request) {
	contact := persistence.Contact{}
	if err := json.NewDecoder(r.Body).Decode(&contact); nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding contact data %s", err)
		return
	}

	if dbErr := hdlr.dbHandler.UpsertContact(&contact); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while persisting contact %s", dbErr)
	} else {
		w.WriteHeader(204)
	}
}

// @Summary Delete an contact
// @Description Delete an contact
// @ID cloud-bill-saas-subscription-service-delete-contact
// @Accept  json
// @Produce  json
// @Param accountId path string true "Account ID"
// @Success 204 {string} string "Deleted"
// @Failure 400 {string} string "Missing account ID in path"
// @Failure 500 {string} string "Error"
// @Router /contacts/{accountId} [delete]
func (hdlr *SubscriptionServiceHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["accountId"]

	if accountId == "" {
		http.Error(w,`{"error": "missing contact name in path"}`,400)
		return
	}

	if dbErr := hdlr.dbHandler.DeleteContact(accountId); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while deleting contact %s", dbErr)
	} else {
		w.WriteHeader(204)
	}
}

// @Summary GetContacts
// @Description Gets an array of contacts
// @ID cloud-bill-saas-subscription-service-get-contacts
// @Accept  json
// @Produce  json
// @Param filters query string false "optional comma separated list of filter"
// @Param order query string false "optional order"
// @Success 200 {array} persistence.Contact
// @Failure 500 {string} string "Error"
// @Router /contacts [get]
func (hdlr *SubscriptionServiceHandler) GetContacts(w http.ResponseWriter, r *http.Request){
	filtersParam, ok := r.URL.Query()["filters"]
	var filters []string = nil
	if ok || len(filtersParam) > 0 {
		filters = strings.Split(filtersParam[0],",")
	}

	ordersParam, ok := r.URL.Query()["order"]
	var order = ""
	if ok || len(ordersParam) > 0 {
		order = ordersParam[0]
	}

	if contacts, dbErr := hdlr.dbHandler.QueryContacts(filters,order); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting account %s", dbErr)
	} else {
		if contacts == nil {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(&contacts)
		}
	}
}

// @Summary Get an entitlement
// @Description Retrieves an entitlement by entitlement ID
// @ID cloud-bill-saas-subscription-service-get-entitlement
// @Accept  json
// @Produce  json
// @Param entitlementId path string true "Entitlement ID"
// @Success 200 {object} persistence.Entitlement
// @Failure 400 {string} string "Missing entitlement ID in path"
// @Failure 500 {string} string "Error"
// @Router /entitlements/{entitlementId} [get]
func (hdlr *SubscriptionServiceHandler) GetEntitlement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entitlementId := vars["entitlementId"]

	if entitlementId == "" {
		http.Error(w,`{"error": "missing entitlement ID in path"}`,400)
		return
	}

	if entitlement, dbErr := hdlr.dbHandler.GetEntitlement(entitlementId); nil != dbErr {
		if dbErr.Error() == "datastore: no such entity" {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "error occured while getting entitlement %s", dbErr)
		}
	} else {
		if entitlement == nil {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(&entitlement)
		}
	}
}

// @Summary GetEntitlements
// @Description Gets an array of entitlements
// @ID cloud-bill-saas-subscription-service-get-entitlements
// @Accept  json
// @Produce  json
// @Param filters query string false "optional comma separated list of filter"
// @Param order query string false "optional order"
// @Success 200 {array} persistence.Entitlement
// @Failure 500 {string} string "Error"
// @Router /entitlements [get]
func (hdlr *SubscriptionServiceHandler) GetEntitlements(w http.ResponseWriter, r *http.Request){
	filtersParam, ok := r.URL.Query()["filters"]
	var filters []string = nil
	if ok || len(filtersParam) > 0 {
		filters = strings.Split(filtersParam[0],",")
	}

	ordersParam, ok := r.URL.Query()["order"]
	var order = ""
	if ok || len(ordersParam) > 0 {
		order = ordersParam[0]
	}

	if entitlements, dbErr := hdlr.dbHandler.QueryEntitlements(filters,order); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting account %s", dbErr)
	} else {
		if entitlements == nil {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(&entitlements)
		}
	}
}

// @Summary GetAccountEntitlements
// @Description Gets an array of entitlements for an account
// @ID cloud-bill-saas-subscription-service-get-account-entitlements
// @Accept  json
// @Produce  json
// @Param acct path string true "account"
// @Param filters query string false "optional comma separated list of filter"
// @Param order query string false "optional order"
// @Success 200 {array} persistence.Entitlement
// @Failure 500 {string} string "Error"
// @Router /accounts/{accountId}/entitlements [get]
func (hdlr *SubscriptionServiceHandler) GetAccountEntitlements(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	accountId := vars["accountId"]

	if accountId == "" {
		http.Error(w,`{"error": "missing account ID"}`,400)
		return
	}

	filtersParam, ok := r.URL.Query()["filters"]
	var filters []string = nil
	if ok || len(filtersParam) > 0 {
		filters = strings.Split(filtersParam[0],",")
	}

	ordersParam, ok := r.URL.Query()["order"]
	var order = ""
	if ok || len(ordersParam) > 0 {
		order = ordersParam[0]
	}

	if entitlements, dbErr := hdlr.dbHandler.QueryAccountEntitlements(accountId,filters,order); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting account %s", dbErr)
	} else {
		if entitlements == nil {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(&entitlements)
		}
	}
}

// @Summary Upsert an entitlement
// @Description Upsert an entitlement passing entitlement json
// @ID cloud-bill-saas-subscription-service-upsert-entitlement
// @Accept  json
// @Produce  json
// @Success 204 {string} string "Upserted"
// @Failure 500 {string} string "Error"
// @Router /entitlements [put]
func (hdlr *SubscriptionServiceHandler) UpsertEntitlement(w http.ResponseWriter, r *http.Request) {
	entitlement := persistence.Entitlement{}
	if err := json.NewDecoder(r.Body).Decode(&entitlement); nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding entitlement data %s", err)
		return
	}
	if dbErr := hdlr.dbHandler.UpsertEntitlement(&entitlement); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while persisting entitlement %s", dbErr)
	} else {
		w.WriteHeader(204)
	}
}

// @Summary Delete an entitlement
// @Description Delete an entitlement
// @ID cloud-bill-saas-subscription-service-delete-entitlement
// @Accept  json
// @Produce  json
// @Param entitlementId path string true "Entitlement ID"
// @Success 204 {string} string "Deleted"
// @Failure 400 {string} string "Missing entitlement ID in path"
// @Failure 500 {string} string "Error"
// @Router /entitlements/{entitlementId} [delete]
func (hdlr *SubscriptionServiceHandler) DeleteEntitlement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entitlementId := vars["entitlementId"]

	if entitlementId == "" {
		http.Error(w,`{"error": "missing entitlement ID in path"}`,400)
		return
	}

	if dbErr := hdlr.dbHandler.DeleteEntitlement(entitlementId); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while deleting entitlement %s", dbErr)
	} else {
		w.WriteHeader(204)
	}
}

// @Summary Check the health of the subscription service
// @Description Check the health of the subscription service
// @ID cloud-bill-saas-subscription-service-healthz
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Ok"
// @Failure 500 {string} string "Error"
// @Router /healthz [get]
func (hdlr *SubscriptionServiceHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	if dbErr := hdlr.dbHandler.Healthz(); nil != dbErr {
		log.Printf("Healthz failed. Datastore check failed: %#v \n", dbErr)
		http.Error(w,dbErr.Error(),http.StatusInternalServerError)
	} else {
		w.WriteHeader(200)
	}
}




