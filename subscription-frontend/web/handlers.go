package web

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/cloudbees/jenkins-support-saas/frontend-service/auth"
	"encoding/json"
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
	Store *sessions.FilesystemStore
)

type SubscriptionFrontendHandler struct {
	SubscriptionServiceUrl string
	ClientId string
	ClientSecret string
	CallbackUrl string
	Issuer string
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

func Init() error {
	Store = sessions.NewFilesystemStore("", []byte("something-very-secret"))
	gob.Register(map[string]interface{}{})
	return nil
}

func GetSubscriptionFrontendHandler(subscriptionServiceUrl string,clientId string, clientSecret string, callbackUrl string, issuer string) *SubscriptionFrontendHandler {
	return &SubscriptionFrontendHandler{
		subscriptionServiceUrl,
		clientId,
		clientSecret,
		callbackUrl,
		issuer,
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

	tmpl, err := template.ParseFiles("../templates/signup.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,sub)
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

	tmpl, err := template.ParseFiles("../templates/confirm.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	jsonBytes, err := json.Marshal(account)
	if err != nil {
		fmt.Println(err)
		return
	}

	url := hdlr.SubscriptionServiceUrl + "/accounts"

	req, err := http.NewRequest(http.MethodPut,url,bytes.NewBuffer(jsonBytes))
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error": "error with creating account upsert request %s"}`, err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		w.WriteHeader(resp.StatusCode)
		fmt.Fprintf(w, `{"error": "error received from subscription service %s"}`, err)
		return
	}
	defer resp.Body.Close()

	//confirm procurement api

	tmpl, err := template.ParseFiles("../templates/finish.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w,account)
}






