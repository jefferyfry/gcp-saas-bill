package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http/httputil"
	"strings"

	"context"
	"encoding/base64"
	"errors"
	"github.com/cloudbees/cloud-bill-saas/frontend-service/auth"
	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/lestrrat/go-jwx/jwk"
	"html/template"
	"log"

	"crypto/rand"
	"net/http"
)

var (
	Store *sessions.CookieStore
)

type SubscriptionFrontendHandler struct {
	SubscriptionServiceUrl string
	ClientId string
	ClientSecret string
	CallbackUrl string
	Issuer string
	CloudCommerceProcurementUrl string
	PartnerId string
}

type Contact struct {
	//google fields
	AccountName  	string     	`json:"accountName"`
	FirstName 		string     	`json:"firstName"`
	LastName		string     	`json:"lastName"`
	EmailAddress	string     	`json:"emailAddress"`
	Phone			string     	`json:"phone"`
	Company			string     	`json:"company"`
	Timezone		string     	`json:"timezone"`
}

func GetSubscriptionFrontendHandler(subscriptionServiceUrl string,clientId string, clientSecret string, callbackUrl string, issuer string, sessionKey string, cloudCommerceProcurementUrl string, partnerId string) *SubscriptionFrontendHandler {
	Store = sessions.NewCookieStore([]byte(sessionKey))

	Store.Options = &sessions.Options{
		MaxAge:   60 * 60,
		HttpOnly: true,
	}
	return &SubscriptionFrontendHandler{
		subscriptionServiceUrl,
		clientId,
		clientSecret,
		callbackUrl,
		issuer,
		cloudCommerceProcurementUrl,
		partnerId,
	}
}


//gets google jwt token and stores. sends to signup.html
func (hdlr *SubscriptionFrontendHandler) SignupSaas(w http.ResponseWriter, r *http.Request) {
	tknStr := r.Header.Get("x-gcp-marketplace-token")
	if tknStr == "" {
		http.Error(w, "x-gcp-marketplace-token not found.", http.StatusInternalServerError)
		return
	}

	tkn, err := jwt.Parse(tknStr, func(token *jwt.Token) (interface{}, error) {
		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("Expecting JWT header to have string kid.")
		}
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			return nil, errors.New("Unable to get JWT payload.")
		}

		iss, ok := claims["iss"].(string)

		if iss == "" || !ok {
			return nil, errors.New("Unable to get JWT payload.")
		}

		keySet, err := jwk.FetchHTTP(iss);

		if err != nil {
			return nil, err
		}

		if key := keySet.LookupKeyID(keyID); len(key) == 1 {
			return key[0].Materialize()
		}

		return nil, errors.New("Unable to find key.")
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !tkn.Valid {
		http.Error(w, "JWT token is not valid.", http.StatusInternalServerError)
		return
	}

	//TODO validate aud vs domain https://cloud.google.com/marketplace/docs/partners/integrated-saas/frontend-integration#verify-jwt

	claims, ok := tkn.Claims.(jwt.MapClaims)

	if !ok {
		http.Error(w, "Unable to get JWT payload.", http.StatusInternalServerError)
		return
	}

	sub, ok := claims["sub"].(string)

	if sub == "" || !ok {
		http.Error(w, "Unable to get sub from JWT payload.", http.StatusInternalServerError)
		return
	}

	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["acct"] = sub
	session.Save(r,w)

	tmpl, err := template.ParseFiles("templates/signup.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,sub)
}

func (hdlr *SubscriptionFrontendHandler) SignupSaasTest(w http.ResponseWriter, r *http.Request) {
	acct, ok := r.URL.Query()["acct"]

	if !ok || len(acct[0]) < 1 {
		http.Error(w, "Missing acct parameter.", http.StatusBadRequest)
		return
	}

	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["acct"] = acct[0]
	session.Save(r,w)

	tmpl, err := template.ParseFiles("templates/signup.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,acct)
}

func (hdlr *SubscriptionFrontendHandler) SignupProd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountName := vars["accountName"]

	if accountName == "" {
		http.Error(w,`{"error": "missing account name in path"}`,400)
		return
	}

	prod, ok := r.URL.Query()["prod"]

	if !ok || len(prod[0]) < 1 {
		http.Error(w, "Missing prod parameter.", http.StatusBadRequest)
		return
	}

	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["acct"] = accountName
	session.Values["prod"] = prod[0]
	session.Save(r,w)

	tmpl, err := template.ParseFiles("templates/signup.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,accountName)
}

func (hdlr *SubscriptionFrontendHandler) EmailConfirm(w http.ResponseWriter, r *http.Request) {
	email, ok := r.URL.Query()["email"]

	if !ok {
		http.Error(w, "Missing email parameter.", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("templates/emailConfirm.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,email)
}

//redirects to Auth0 for authentication, this should not be called unless istio fails
func (hdlr *SubscriptionFrontendHandler) Auth0Login(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authenticator, err := auth.NewAuthenticator(hdlr.Issuer,hdlr.ClientId,hdlr.ClientSecret,hdlr.CallbackUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

//handles auth0 callback, stores profile data, confirms account and redirects to confirmation
func (hdlr *SubscriptionFrontendHandler) Auth0Callback(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	authenticator, err := auth.NewAuthenticator(hdlr.Issuer,hdlr.ClientId,hdlr.ClientSecret,hdlr.CallbackUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("no token found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: hdlr.ClientId,
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

	if err != nil {
		http.Error(w, "Failed to verify ID Token: " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//need to check if we have first last
	if profile["firstName"] == nil && profile["name"] != nil {
		split := strings.Split(profile["name"].(string), " ")
		ln := len(split)
		switch(ln){
			case 2:
				profile["lastName"] = split[1]
				fallthrough
			case 1:
				profile["firstName"] = split[0]
		}
	}

	profile["acct"] = session.Values["acct"]
	prod := session.Values["prod"]
	if prod == nil {
		profile["prod"] = "saas"
	} else {
		profile["prod"] = session.Values["prod"]
	}
	fmt.Println("map:", profile)
	tmpl, err := template.ParseFiles("templates/confirm.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,profile)
}

func (hdlr *SubscriptionFrontendHandler) FinishTest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/confirm.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var profile map[string]interface{}
	tmpl.Execute(w,profile)
}

func (hdlr *SubscriptionFrontendHandler) Finish(w http.ResponseWriter, r *http.Request) {
	contact := Contact{}
	contact.AccountName = r.PostFormValue("acct")
	contact.Company = r.PostFormValue("company")
	contact.EmailAddress = r.PostFormValue("emailAddress")
	contact.FirstName = r.PostFormValue("firstName")
	contact.LastName = r.PostFormValue("lastName")
	contact.Phone = r.PostFormValue("phone")
	contact.Timezone = r.PostFormValue("timezone")
	prod := r.PostFormValue("prod")

	if !createContact(contact, hdlr.SubscriptionServiceUrl, w) {
		http.Error(w, "Failed to store contact info", http.StatusInternalServerError)
		return
	}

	if prod == "saas" {
		sendAccountApprove(hdlr.CloudCommerceProcurementUrl, hdlr.PartnerId, contact.AccountName, w)
	} else { //integrated prod support
		//TODO query subscription API for more details
		//createProduct(prod,contact.AccountName,hdlr.SubscriptionServiceUrl, w)
	}

	//delete the session
	session, err := Store.Get(r, "auth-session")
	if err == nil {
		session.Options.MaxAge = -1
		err = session.Save(r, w)
		if err != nil {
			fmt.Println("failed to delete session")
		}
	}

	tmpl, err := template.ParseFiles("templates/finish.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, contact)
}

/*func createProduct(prod string, accountName string, subscriptionServiceUrl string, w http.ResponseWriter) bool {
	contactReq, err := http.NewRequest(http.MethodPut, subscriptionServiceUrl+"/contacts", bytes.NewBuffer(contactBytes))
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error": "error with creating contact upsert request %s"}`, err)
		return false
	}
	contactResp, err := http.DefaultClient.Do(contactReq)
	if nil != err {
		w.WriteHeader(contactResp.StatusCode)
		fmt.Fprintf(w, `{"error": "error received from subscription service %s"}`, err)
		return false
	}
	defer contactResp.Body.Close()
	return true
}*/

func createContact(contact Contact, subscriptionServiceUrl string, w http.ResponseWriter) bool {
	//submit to subscript service
	contactBytes, err := json.Marshal(contact)
	if err != nil {
		fmt.Println(err)
		return false
	}
	contactReq, err := http.NewRequest(http.MethodPut, subscriptionServiceUrl+"/contacts", bytes.NewBuffer(contactBytes))
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error": "error with creating contact upsert request %s"}`, err)
		return false
	}
	contactResp, err := http.DefaultClient.Do(contactReq)
	if nil != err {
		w.WriteHeader(contactResp.StatusCode)
		fmt.Fprintf(w, `{"error": "error received from subscription service %s"}`, err)
		return false
	}
	defer contactResp.Body.Close()
	return true
}

func sendAccountApprove(cloudCommerceProcurementUrl string, partnerId string,accountName string, w http.ResponseWriter) bool {
	procurementUrl := cloudCommerceProcurementUrl + "/" + partnerId + "/accounts/" + accountName + ":approve"
	jsonApproval := []byte(`
			{
				"approvalName": "signup"
			}`)
	procReq, err := http.NewRequest(http.MethodPut, procurementUrl, bytes.NewBuffer(jsonApproval))
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error": "error with creating account approval request %s"}`, err)
		return false
	}
	procResp, err := http.DefaultClient.Do(procReq)
	if nil != err {
		w.WriteHeader(procResp.StatusCode)
		fmt.Fprintf(w, `{"error": "error received from procurement api service %s"}`, err)
		return false
	}
	defer procResp.Body.Close()
	responseDump, err := httputil.DumpResponse(procResp, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(responseDump))
	return true
}






