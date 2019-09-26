package web

import (
	"encoding/json"
	"fmt"
	"github.com/cloudbees/cloud-bill-saas/subscription-service/persistence"
	"github.com/gorilla/mux"
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
// @Description Retrieves an account by account name
// @ID cloud-bill-saas-subscription-service-get-account
// @Accept  json
// @Produce  json
// @Param accountName path string true "Account Name"
// @Success 200 {object} persistence.Account
// @Failure 400 {string} string "Missing account name in path"
// @Failure 500 {string} string "Internal server error"
// @Router /accounts/{accountName} [get]
func (hdlr *SubscriptionServiceHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountName := vars["accountName"]

	if accountName == "" {
		http.Error(w,`{"error": "missing account name"}`,400)
		return
	}

	if account, dbErr := hdlr.dbHandler.GetAccount(accountName); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting account %s", dbErr)
		return
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&account)
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
	vars := mux.Vars(r)
	filtersParam := vars["filters"]
	var filters []string = nil
	if filtersParam != "" {
		filters = strings.Split(filtersParam,",")
	}
	orderParam := vars["order"]

	if accounts, dbErr := hdlr.dbHandler.QueryAccounts(filters,orderParam); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting account %s", dbErr)
		return
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&accounts)
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
// @Param accountName path string true "Account Name"
// @Success 204 {string} string "Deleted"
// @Failure 400 {string} string "Missing account name in path"
// @Failure 500 {string} string "Internal server error"
// @Router /accounts/{accountName} [delete]
func (hdlr *SubscriptionServiceHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountName := vars["accountName"]

	if accountName == "" {
		http.Error(w,`{"error": "missing account name"}`,400)
		return
	}

	if dbErr := hdlr.dbHandler.DeleteAccount(accountName); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while deleting account %s", dbErr)
	} else {
		w.WriteHeader(204)
	}
}

// @Summary Get an contact
// @Description Retrieves an contact by account name
// @ID cloud-bill-saas-subscription-service-get-contact
// @Accept  json
// @Produce  json
// @Param accountName path string true "Account Name"
// @Success 200 {object} persistence.Contact
// @Failure 400 {string} string "Missing account name in path"
// @Failure 500 {string} string "Error"
// @Router /contacts/{accountName} [get]
func (hdlr *SubscriptionServiceHandler) GetContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountName := vars["accountName"]

	if accountName == "" {
		http.Error(w,`{"error": "missing account name in path"}`,400)
		return
	}

	if contact, dbErr := hdlr.dbHandler.GetContact(accountName); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting contact %s", dbErr)
		return
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&contact)
	}
}

// @Summary Upsert an contact
// @Description Upsert an contact passing contact json
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
// @Param accountName path string true "Account Name"
// @Success 204 {string} string "Deleted"
// @Failure 400 {string} string "Missing account name in path"
// @Failure 500 {string} string "Error"
// @Router /contacts/{accountName} [delete]
func (hdlr *SubscriptionServiceHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountName := vars["accountName"]

	if accountName == "" {
		http.Error(w,`{"error": "missing contact name in path"}`,400)
		return
	}

	if dbErr := hdlr.dbHandler.DeleteContact(accountName); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while deleting contact %s", dbErr)
	} else {
		w.WriteHeader(204)
	}
}

// @Summary Get an entitlement
// @Description Retrieves an entitlement by entitlement name
// @ID cloud-bill-saas-subscription-service-get-entitlement
// @Accept  json
// @Produce  json
// @Param entitlementName path string true "Entitlement Name"
// @Success 200 {object} persistence.Entitlement
// @Failure 400 {string} string "Missing entitlement name in path"
// @Failure 500 {string} string "Error"
// @Router /entitlements/{entitlementName} [get]
func (hdlr *SubscriptionServiceHandler) GetEntitlement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entitlementName := vars["entitlementName"]

	if entitlementName == "" {
		http.Error(w,`{"error": "missing entitlement name in path"}`,400)
		return
	}

	if entitlement, dbErr := hdlr.dbHandler.GetEntitlement(entitlementName); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting entitlement %s", dbErr)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&entitlement)
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
	vars := mux.Vars(r)
	filtersParam := vars["filters"]
	var filters []string = nil
	if filtersParam != "" {
		filters = strings.Split(filtersParam,",")
	}
	orderParam := vars["order"]

	if entitlements, dbErr := hdlr.dbHandler.QueryEntitlements(filters,orderParam); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting account %s", dbErr)
		return
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&entitlements)
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
// @Router /accounts/{accountName}/entitlements [get]
func (hdlr *SubscriptionServiceHandler) GetAccountEntitlements(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	accountName := vars["acct"]

	if accountName == "" {
		http.Error(w,`{"error": "missing account name"}`,400)
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

	if entitlements, dbErr := hdlr.dbHandler.QueryAccountEntitlements(accountName,filters,order); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting account %s", dbErr)
		return
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&entitlements)
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
// @Param entitlementName path string true "Entitlement Name"
// @Success 204 {string} string "Deleted"
// @Failure 400 {string} string "Missing entitlement name in path"
// @Failure 500 {string} string "Error"
// @Router /entitlements/{entitlementName} [delete]
func (hdlr *SubscriptionServiceHandler) DeleteEntitlement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entitlementName := vars["entitlementName"]

	if entitlementName == "" {
		http.Error(w,`{"error": "missing entitlement name in path"}`,400)
		return
	}

	if dbErr := hdlr.dbHandler.DeleteEntitlement(entitlementName); nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while deleting entitlement %s", dbErr)
	} else {
		w.WriteHeader(204)
	}
}




