package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httputil"

	"context"
	"encoding/base64"
	"errors"
	"github.com/cloudbees/jenkins-support-saas/frontend-service/auth"
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

type Account struct {
	//google fields
	Name  			string     	`json:"name"`
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
func (hdlr *SubscriptionFrontendHandler) Signup(w http.ResponseWriter, r *http.Request) {
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
	session.Values["sub"] = sub
	session.Save(r,w)

	tmpl, err := template.ParseFiles("templates/signup.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,sub)
}

func (hdlr *SubscriptionFrontendHandler) SignupTest(w http.ResponseWriter, r *http.Request) {
	sub, ok := r.URL.Query()["sub"]

	if !ok || len(sub[0]) < 1 {
		http.Error(w, "Missing sub parameter.", http.StatusBadRequest)
		return
	}

	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["sub"] = sub[0]
	session.Save(r,w)

	tmpl, err := template.ParseFiles("templates/signup.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,sub)
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

	profile["sub"] = session.Values["sub"]
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
	account := Account{}
	account.Name = r.PostFormValue("sub")
	account.Company = r.PostFormValue("company")
	account.EmailAddress = r.PostFormValue("emailAddress")
	account.FirstName = r.PostFormValue("firstName")
	account.LastName = r.PostFormValue("lastName")
	account.Phone = r.PostFormValue("phone")
	account.Timezone = r.PostFormValue("timezone")

	//submit to subscript service
	accountBytes, err := json.Marshal(account)
	if err != nil {
		fmt.Println(err)
		return
	}

	subscriptionServiceUrl := hdlr.SubscriptionServiceUrl + "/accounts"

	subReq, err := http.NewRequest(http.MethodPut,subscriptionServiceUrl,bytes.NewBuffer(accountBytes))
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error": "error with creating account upsert request %s"}`, err)
		return
	}

	subResp, err := http.DefaultClient.Do(subReq)
	if nil != err {
		w.WriteHeader(subResp.StatusCode)
		fmt.Fprintf(w, `{"error": "error received from subscription service %s"}`, err)
		return
	}
	defer subResp.Body.Close()

	//confirm account with procurement api
	procurementUrl := hdlr.CloudCommerceProcurementUrl+"/"+hdlr.PartnerId + "/accounts/"+account.Name+":approve"
	jsonApproval := []byte(`
		{
			"approvalName": "signup"
		}`)

	procReq, err := http.NewRequest(http.MethodPut,procurementUrl,bytes.NewBuffer(jsonApproval))
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error": "error with creating account approval request %s"}`, err)
		return
	}

	procResp, err := http.DefaultClient.Do(procReq)
	if nil != err {
		w.WriteHeader(procResp.StatusCode)
		fmt.Fprintf(w, `{"error": "error received from procurement api service %s"}`, err)
		return
	}
	defer procResp.Body.Close()

	responseDump, err := httputil.DumpResponse(procResp,true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(responseDump))

	tmpl, err := template.ParseFiles("templates/finish.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,account)
}






