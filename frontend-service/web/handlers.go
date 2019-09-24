package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http/httputil"
	"strings"
	"time"

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

type Product struct {
	Name     			string	`json:"name"`
	Account   			string	`json:"account"`
	Product  			string	`json:"product"`
	State    	  		string	`json:"state"`
	UpdateTime    	  	string	`json:"updateTime"`
	CreateTime    	  	string	`json:"createTime"`
}

type Contact struct {
	AccountName 	string     	`json:"accountName"`
	FirstName 		string     	`json:"firstName,omitempty"`
	LastName		string     	`json:"lastName,omitempty"`
	EmailAddress	string     	`json:"emailAddress"`
	Phone			string     	`json:"phone,omitempty"`
	Company			string     	`json:"company,omitempty"`
	Timezone		string     	`json:"timezone,omitempty"`
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

func (hdlr *SubscriptionFrontendHandler) ResetSaas(w http.ResponseWriter, r *http.Request) {
	acct, acctOk := r.URL.Query()["acct"]

	if !acctOk || len(acct[0]) < 1 {
		http.Error(w, "Missing acct parameter.", http.StatusBadRequest)
		return
	}

	err := postAccountReset(hdlr.CloudCommerceProcurementUrl, hdlr.PartnerId, acct[0], w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w,"Account %s has been reset.",acct[0])
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
		switch ln {
			case 2:
				profile["lastName"] = split[1]
				fallthrough
			case 1:
				profile["firstName"] = split[0]
		}
	}

	profile["acct"] = session.Values["acct"]
	profile["prod"] = "saas"
	prod := session.Values["prod"]
	tmplHtml := "templates/confirmSaas.html"
	if prod != nil {
		profile["prod"] = session.Values["prod"]
		tmplHtml = "templates/confirmProd.html"
	}

	tmpl, err := template.ParseFiles(tmplHtml)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,profile)
}

func (hdlr *SubscriptionFrontendHandler) FinishSaas(w http.ResponseWriter, r *http.Request) {
	contact := Contact{}
	contact.AccountName = r.PostFormValue("acct")
	contact.Company = r.PostFormValue("company")
	contact.EmailAddress = r.PostFormValue("emailAddress")
	contact.FirstName = r.PostFormValue("firstName")
	contact.LastName = r.PostFormValue("lastName")
	contact.Phone = r.PostFormValue("phone")
	contact.Timezone = r.PostFormValue("timezone")

	if !createContact(contact, hdlr.SubscriptionServiceUrl, w) {
		http.Error(w, "Failed to store contact info", http.StatusInternalServerError)
		return
	}

	postAccountApproval(hdlr.CloudCommerceProcurementUrl, hdlr.PartnerId, contact.AccountName, w)

	tmpl, err := template.ParseFiles("templates/finish.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, contact)
}

func (hdlr *SubscriptionFrontendHandler) FinishProd(w http.ResponseWriter, r *http.Request) {
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

	//TODO query subscription API

	createProduct(prod,contact.AccountName,hdlr.SubscriptionServiceUrl, w)

	tmpl, err := template.ParseFiles("templates/finish.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, contact)
}

func createProduct(prod string, accountName string, subscriptionServiceUrl string, w http.ResponseWriter) bool {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc).String()
	product := Product {
		Name: accountName+"-"+prod,
		Product: prod,
		Account: accountName,
		CreateTime: now,
		UpdateTime: now,
	}
	productBytes, err := json.Marshal(product)
	prodReq, err := http.NewRequest(http.MethodPut, subscriptionServiceUrl+"/entitlements", bytes.NewBuffer(productBytes))
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error": "error with creating product upsert request %s"}`, err)
		return false
	}
	prodResp, err := http.DefaultClient.Do(prodReq)
	if nil != err {
		w.WriteHeader(prodResp.StatusCode)
		fmt.Fprintf(w, `{"error": "error received from subscription service %s"}`, err)
		return false
	}
	defer prodResp.Body.Close()
	return true
}

func createContact(contact Contact, subscriptionServiceUrl string, w http.ResponseWriter) bool {
	//submit to subscript service
	contactBytes, err := json.Marshal(contact)
	if err != nil {
		log.Println(err)
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

func postAccountApproval(cloudCommerceProcurementUrl string, partnerId string,accountName string, w http.ResponseWriter) error {
	procurementUrl := cloudCommerceProcurementUrl +  "/providers/" +  partnerId + "/accounts/" + accountName + ":approve"
	jsonApproval := []byte(`
			{
				"approvalName": "signup"
			}`)
	client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform")

	if clientErr != nil {
		log.Printf("Failed to create oath2 client for the procurement API %#v \n", clientErr)
		return clientErr
	}

	log.Printf("Sending account approval: %s \n", procurementUrl)
	procResp, err := client.Post(procurementUrl,"",bytes.NewBuffer(jsonApproval))
	if nil != err {
		log.Printf("Failed sending entitlement approval request %s %#v \n",procurementUrl, err)
		return err
	}
	defer procResp.Body.Close()
	if procResp.StatusCode != 200 {
		log.Println("Account approval received error response: ",procResp.StatusCode)
		responseDump, _ := httputil.DumpResponse(procResp, true)
		log.Println(string(responseDump))
		return errors.New("Account approval received error response: "+procResp.Status)
	} else {
		log.Printf("%s %s",procurementUrl,procResp.Status)
	}
	return nil
}

func postAccountReset(cloudCommerceProcurementUrl string, partnerId string,accountName string,w http.ResponseWriter) error {
	procurementUrl := cloudCommerceProcurementUrl +  "/providers/" +  partnerId + "/accounts/" + accountName + ":reset"
	client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform")

	if clientErr != nil {
		log.Printf("Failed to create oath2 client for the procurement API %#v \n", clientErr)
		return clientErr
	}

	log.Printf("Sending account reset: %s \n", procurementUrl)
	procResp, err := client.Post(procurementUrl,"",nil)
	if nil != err {
		log.Printf("Failed sending account reset request %s %#v \n",procurementUrl, err)
		return err
	}
	defer procResp.Body.Close()
	if procResp.StatusCode != 200 {
		log.Println("Account reset received error response: ",procResp.StatusCode)
		responseDump, _ := httputil.DumpResponse(procResp, true)
		log.Println(string(responseDump))
		return errors.New("Account reset received error response: "+procResp.Status)
	}
	return nil
}






