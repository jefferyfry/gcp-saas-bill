package rest

import (
	"encoding/json"
	"fmt"
	"github.com/cloudbees/jenkins-support-saas/subscription-service/persistence"
	"github.com/gorilla/mux"
	"net/http"
)

type SubscriptionServiceHandler struct {
	dbHandler    persistence.DatabaseHandler
	cloudCommerceProcurementUrl string
	partnerId                   string
}

func GetSubscriptionServiceHandler(dbHandler persistence.DatabaseHandler,cloudCommerceProcurementUrl string,partnerId string) *SubscriptionServiceHandler {
	return &SubscriptionServiceHandler {
		dbHandler,cloudCommerceProcurementUrl, partnerId ,
	}
}

// @Summary Add a new account
// @Description Adds a new account
// @ID jenkins-support-saas-subscription-service-add-account
// @Accept  json
// @Produce  json
// @Success 201 {string} string "Added"
// @Failure 500 {string} string "Error"
// @Router /accounts [post]
func (hdlr *SubscriptionServiceHandler) AddAccount(w http.ResponseWriter, r *http.Request) {
	account := persistence.Account{}
	err := json.NewDecoder(r.Body).Decode(&account)

	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding account data %s", err)
		return
	}
	dbErr := hdlr.dbHandler.AddAccount(&account)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while persisting account %s", dbErr)
		return
	}

	w.WriteHeader(201)
}

// @Summary Get an account
// @Description Retrieves an account by account name
// @ID jenkins-support-saas-subscription-service-get-account
// @Accept  json
// @Produce  json
// @Param accountName path string true "Account Name"
// @Success 200 {object} persistence.Account
// @Failure 500 {string} string "Error"
// @Router /accounts/{accountName} [get]
func (hdlr *SubscriptionServiceHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountName := vars["accountName"]

	account, dbErr := hdlr.dbHandler.GetAccount(accountName)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting account %s", dbErr)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&account)
}

// @Summary Update an account
// @Description Update an account passing account json
// @ID jenkins-support-saas-subscription-service-update-account
// @Accept  json
// @Produce  json
// @Success 204 {string} string "Updated"
// @Failure 500 {string} string "Error"
// @Router /accounts [put]
func (hdlr *SubscriptionServiceHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	account := persistence.Account{}
	err := json.NewDecoder(r.Body).Decode(&account)

	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding account data %s", err)
		return
	}
	dbErr := hdlr.dbHandler.UpdateAccount(&account)
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
// @Failure 500 {string} string "Error"
// @Router /accounts/{accountName} [delete]
func (hdlr *SubscriptionServiceHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountName := vars["accountName"]

	dbErr := hdlr.dbHandler.DeleteAccount(accountName)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while deleting account %s", dbErr)
		return
	}

	w.WriteHeader(204)
}

// @Summary Add a new entitlement
// @Description Adds a new entitlement
// @ID jenkins-support-saas-subscription-service-add-entitlement
// @Accept  json
// @Produce  json
// @Success 201 {string} string "Added"
// @Failure 500 {string} string "Error"
// @Router /entitlements [post]
func (hdlr *SubscriptionServiceHandler) AddEntitlement(w http.ResponseWriter, r *http.Request) {
	entitlement := persistence.Entitlement{}
	err := json.NewDecoder(r.Body).Decode(&entitlement)

	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding entitlement data %s", err)
		return
	}
	dbErr := hdlr.dbHandler.AddEntitlement(&entitlement)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while persisting entitlement %s", dbErr)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&entitlement)
}

// @Summary Get an entitlement
// @Description Retrieves an entitlement by entitlement name
// @ID jenkins-support-saas-subscription-service-get-entitlement
// @Accept  json
// @Produce  json
// @Param entitlementName path string true "Entitlement Name"
// @Success 200 {object} persistence.Entitlement
// @Failure 500 {string} string "Error"
// @Router /entitlements/{entitlementName} [get]
func (hdlr *SubscriptionServiceHandler) GetEntitlement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entitlementName := vars["entitlementName"]

	entitlement, dbErr := hdlr.dbHandler.GetEntitlement(entitlementName)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while getting entitlement %s", dbErr)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&entitlement)
}

// @Summary Update an entitlement
// @Description Update an entitlement passing entitlement json
// @ID jenkins-support-saas-subscription-service-update-entitlement
// @Accept  json
// @Produce  json
// @Success 204 {string} string "Updated"
// @Failure 500 {string} string "Error"
// @Router /entitlements [put]
func (hdlr *SubscriptionServiceHandler) UpdateEntitlement(w http.ResponseWriter, r *http.Request) {
	entitlement := persistence.Entitlement{}
	err := json.NewDecoder(r.Body).Decode(&entitlement)

	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding entitlement data %s", err)
		return
	}
	dbErr := hdlr.dbHandler.UpdateEntitlement(&entitlement)
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
// @Failure 500 {string} string "Error"
// @Router /entitlements/{entitlementName} [delete]
func (hdlr *SubscriptionServiceHandler) DeleteEntitlement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entitlementName := vars["entitlementName"]

	dbErr := hdlr.dbHandler.DeleteEntitlement(entitlementName)
	if nil != dbErr {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while deleting entitlement %s", dbErr)
		return
	}

	w.WriteHeader(204)
}




