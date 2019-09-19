package web

import (
	"encoding/json"
	"fmt"
	"github.com/cloudbees/jenkins-support-saas/subscription-service/persistence"
	"github.com/gorilla/mux"
	"net/http"
)

type SubscriptionServiceHandler struct {
	dbHandler    persistence.DatabaseHandler
}

func GetSubscriptionServiceHandler(dbHandler persistence.DatabaseHandler) *SubscriptionServiceHandler {
	return &SubscriptionServiceHandler {
		dbHandler,
	}
}

// @Summary Get an account
// @Description Retrieves an account by account name
// @ID jenkins-support-saas-subscription-service-get-account
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

	account, dbErr := hdlr.dbHandler.GetAccount(accountName)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting account %s", dbErr)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&account)
}

// @Summary Upsert an account
// @Description Upsert an account passing account json
// @ID jenkins-support-saas-subscription-service-upsert-account
// @Accept  json
// @Produce  json
// @Param account body persistence.Account true "Account"
// @Success 204 {string} string "Upserted"
// @Failure 500 {string} string "Error"
// @Router /accounts [put]
func (hdlr *SubscriptionServiceHandler) UpsertAccount(w http.ResponseWriter, r *http.Request) {
	account := persistence.Account{}
	err := json.NewDecoder(r.Body).Decode(&account)

	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding account data %s", err)
		return
	}
	dbErr := hdlr.dbHandler.UpsertAccount(&account)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while persisting account %s", dbErr)
		return
	}

	w.WriteHeader(204)
}

// @Summary Delete an account
// @Description Delete an account
// @ID jenkins-support-saas-subscription-service-delete-account
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

	dbErr := hdlr.dbHandler.DeleteAccount(accountName)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while deleting account %s", dbErr)
		return
	}

	w.WriteHeader(204)
}

// @Summary Get an entitlement
// @Description Retrieves an entitlement by entitlement name
// @ID jenkins-support-saas-subscription-service-get-entitlement
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

	entitlement, dbErr := hdlr.dbHandler.GetEntitlement(entitlementName)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting entitlement %s", dbErr)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&entitlement)
}

// @Summary Upsert an entitlement
// @Description Upsert an entitlement passing entitlement json
// @ID jenkins-support-saas-subscription-service-upsert-entitlement
// @Accept  json
// @Produce  json
// @Success 204 {string} string "Upserted"
// @Failure 500 {string} string "Error"
// @Router /entitlements [put]
func (hdlr *SubscriptionServiceHandler) UpsertEntitlement(w http.ResponseWriter, r *http.Request) {
	entitlement := persistence.Entitlement{}
	err := json.NewDecoder(r.Body).Decode(&entitlement)

	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding entitlement data %s", err)
		return
	}
	dbErr := hdlr.dbHandler.UpsertEntitlement(&entitlement)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while persisting entitlement %s", dbErr)
		return
	}

	w.WriteHeader(204)
}

// @Summary Delete an entitlement
// @Description Delete an entitlement
// @ID jenkins-support-saas-subscription-service-delete-entitlement
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

	dbErr := hdlr.dbHandler.DeleteEntitlement(entitlementName)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while deleting entitlement %s", dbErr)
		return
	}

	w.WriteHeader(204)
}




