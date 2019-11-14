package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http/httputil"
	"regexp"
	"strings"
	"time"

	"context"
	"encoding/base64"
	"errors"
	"github.com/cloudbees/cloud-bill-saas/frontend-service/auth"
	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/jefferyfry/funclog"
	"github.com/lestrrat/go-jwx/jwk"
	"html/template"

	"crypto/rand"
	"net/http"
)

var (
	Store *sessions.CookieStore

	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

type SubscriptionFrontendHandler struct {
	SubscriptionServiceUrl string
	ClientId string
	ClientSecret string
	CallbackUrl string
	Issuer string
	CloudCommerceProcurementUrl string
	PartnerId string
	FinishUrl string
	FinishUrlTitle string
}

type Account struct {
	Id  			string     	`json:"id"`
	UpdateTime   	string    	`json:"updateTime,omitempty"`
	State 	 		string      `json:"state,omitempty"`
	Provider     	string		`json:"provider,omitempty"`
}

type Product struct {
	Id     				string	`json:"id"`
	Name     			string	`json:"name"`
	Account   			string	`json:"account"`
	Product  			string	`json:"product"`
	Plan     	  		string	`json:"plan"`
	State    	  		string	`json:"state"`
	UpdateTime    	  	string	`json:"updateTime"`
	CreateTime    	  	string	`json:"createTime"`
}

type Contact struct {
	AccountId 		string     	`json:"accountId"`
	FirstName 		string     	`json:"firstName,omitempty"`
	LastName		string     	`json:"lastName,omitempty"`
	EmailAddress	string     	`json:"emailAddress"`
	Phone			string     	`json:"phone,omitempty"`
	Company			string     	`json:"company,omitempty"`
	Timezone		string     	`json:"timezone,omitempty"`
}

type Finish struct {
	FinishUrl 		string `json:"finishUrl"`
	FinishUrlTitle 	string `json:"finishUrlTitle"`
}

func GetSubscriptionFrontendHandler(subscriptionServiceUrl string,clientId string, clientSecret string, callbackUrl string, issuer string, sessionKey string, cloudCommerceProcurementUrl string, partnerId string, finishUrl string, finishUrlTitle string) *SubscriptionFrontendHandler {
	Store = sessions.NewCookieStore([]byte(sessionKey))

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
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
		finishUrl,
		finishUrlTitle,
	}
}


//gets google jwt token and stores. sends to signup.html
func (hdlr *SubscriptionFrontendHandler) SignupSaas(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tknStr := r.Form.Get("x-gcp-marketplace-token")
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

		issSub := strings.Replace(iss,"x509","jwk",-1)

		LogI.Printf("Getting keys from %s",issSub)

		keySet, err := jwk.FetchHTTP(issSub);

		if err != nil {
			LogE.Printf("Error fetching keyset from ISS %s: %#v",issSub,err)
			return nil, err
		}

		LogI.Printf("KeySet found %#v",keySet)

		if key := keySet.LookupKeyID(keyID); len(key) == 1 {
			return key[0].Materialize()
		}

		return nil, errors.New("Unable to find key.")
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	if session, err := Store.Get(r, "auth-session"); err != nil {
		LogE.Printf("Unable to get session %#v",err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		session.Values["acct"] = sub
		if err := session.Save(r,w); err != nil {
			LogE.Printf("Unable to save session %#v",err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}


	if tmpl, err := template.ParseFiles("templates/signup.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		tmpl.Execute(w,sub)
	}
}

func (hdlr *SubscriptionFrontendHandler) SignupSaasTest(w http.ResponseWriter, r *http.Request) {
	acct, ok := r.URL.Query()["acct"]

	if !ok || len(acct) < 1 {
		http.Error(w, "Missing acct parameter.", http.StatusBadRequest)
		return
	}

	if session, err := Store.Get(r, "auth-session"); err != nil {
		LogE.Printf("Unable to get session %#v",err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		session.Values["acct"] = acct[0]
		if err := session.Save(r,w); err != nil {
			LogE.Printf("Unable to save session %#v",err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if tmpl, err := template.ParseFiles("templates/signup.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		tmpl.Execute(w, acct)
	}
}

func (hdlr *SubscriptionFrontendHandler) ResetSaas(w http.ResponseWriter, r *http.Request) {
	acct, acctOk := r.URL.Query()["acct"]

	if !acctOk || len(acct) < 1 {
		http.Error(w, "Missing acct parameter.", http.StatusBadRequest)
		return
	}

	if err := postAccountReset(hdlr.CloudCommerceProcurementUrl, hdlr.PartnerId, acct[0], w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w,"Account %s has been reset.",acct[0])
	}
}

func (hdlr *SubscriptionFrontendHandler) SignupProd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["accountId"]

	if accountId == "" {
		http.Error(w,`{"error": "missing account name in path"}`,400)
		return
	}

	matched,_ := regexp.MatchString("^[a-zA-Z0-9_-]*$",accountId)

	if !matched {
		http.Error(w,`{"error": "invalid account ID"}`,400)
		return
	}

	prod, ok := r.URL.Query()["prod"]

	if !ok || len(prod) < 1 {
		http.Error(w, "Missing prod parameter.", http.StatusBadRequest)
		return
	}

	valid, err := accountValid(accountId,prod[0])

	if !valid && err != nil {
		http.Error(w,`{"error": "error validating account"}`,500)
		return
	} else if !valid {
		http.Error(w,`{"error": "invalid account"}`,500)
		return
	}

	if session, err := Store.Get(r, "auth-session"); err != nil {
		LogE.Printf("Unable to get session %#v",err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		session.Values["acct"] = accountId
		session.Values["prod"] = prod[0]
		if err := session.Save(r,w); err != nil {
			LogE.Printf("Unable to save session %#v",err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if tmpl, err := template.ParseFiles("templates/signup.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		tmpl.Execute(w, accountId)
	}
}

//redirects to Auth0 for authentication, this should not be called unless istio fails
func (hdlr *SubscriptionFrontendHandler) Auth0Login(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	if session, err := Store.Get(r, "auth-session"); err != nil {
		LogE.Printf("Unable to get session %#v",err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		session.Values["state"] = state
		if err = session.Save(r, w); err != nil {
			LogE.Printf("Unable to save session %#v",err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if authenticator, err := auth.NewAuthenticator(hdlr.Issuer,hdlr.ClientId,hdlr.ClientSecret,hdlr.CallbackUrl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
	}
}

//handles auth0 callback, stores profile data, confirms account and redirects to confirmation
func (hdlr *SubscriptionFrontendHandler) Auth0Callback(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		LogE.Printf("Unable to get session %#v",err)
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
		LogE.Printf("no token found: %v", err)
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
		profile["prod"] = prod
		tmplHtml = "templates/confirmProd.html"
	}

	if tmpl, err := template.ParseFiles(tmplHtml); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		tmpl.Execute(w, profile)
	}
}

func (hdlr *SubscriptionFrontendHandler) FinishSaas(w http.ResponseWriter, r *http.Request) {
	if session, err := Store.Get(r, "auth-session"); err != nil {
		LogE.Printf("Unable to get session %#v",err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		session.Options.MaxAge = -1
		if err = session.Save(r, w); err != nil {
			LogE.Printf("Unable to delete session %#v",err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	contact := Contact{}
	contact.AccountId = r.PostFormValue("acct")
	contact.Company = r.PostFormValue("company")
	contact.EmailAddress = r.PostFormValue("emailAddress")
	contact.FirstName = r.PostFormValue("firstName")
	contact.LastName = r.PostFormValue("lastName")
	contact.Phone = r.PostFormValue("phone")
	contact.Timezone = r.PostFormValue("timezone")

	if !createContact(contact, hdlr.SubscriptionServiceUrl, w) {
		http.Error(w, "Failed to store contact info", http.StatusInternalServerError)
	} else {
		postAccountApproval(hdlr.CloudCommerceProcurementUrl, hdlr.PartnerId, contact.AccountId, w)

		if tmpl, err := template.ParseFiles("templates/finish.html"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			finish := Finish{}
			finish.FinishUrl = hdlr.FinishUrl
			finish.FinishUrlTitle = hdlr.FinishUrlTitle
			tmpl.Execute(w, finish)
		}
	}
}

func (hdlr *SubscriptionFrontendHandler) FinishProd(w http.ResponseWriter, r *http.Request) {
	if session, err := Store.Get(r, "auth-session"); err != nil {
		LogE.Printf("Unable to get session %#v",err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		session.Options.MaxAge = -1
		if err = session.Save(r, w); err != nil {
			LogE.Printf("Unable to delete session %#v",err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	contact := Contact{
		AccountId : r.PostFormValue("acct"),
		Company : r.PostFormValue("company"),
		EmailAddress : r.PostFormValue("emailAddress"),
		FirstName : r.PostFormValue("firstName"),
		LastName : r.PostFormValue("lastName"),
		Phone : r.PostFormValue("phone"),
		Timezone : r.PostFormValue("timezone"),
	}

	prod := r.PostFormValue("prod")

	if !createContact(contact, hdlr.SubscriptionServiceUrl, w) {
		http.Error(w, "Failed to store contact info", http.StatusInternalServerError)
	} else {
		loc, _ := time.LoadLocation("UTC")
		now := time.Now().In(loc).String()
		account := Account{
			Id : contact.AccountId,
			Provider : hdlr.PartnerId,
			UpdateTime : now,
			State : "ACCOUNT_ACTIVE",
		}

		if !createAccount(account, hdlr.SubscriptionServiceUrl, w) {
			http.Error(w, "Failed to store account info", http.StatusInternalServerError)
		} else {
			if !createProduct(prod, contact.AccountId, hdlr.SubscriptionServiceUrl, w){
				http.Error(w, "Failed to store product info", http.StatusInternalServerError)
			}
		}

		if tmpl, err := template.ParseFiles("templates/finish.html");err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			finish := Finish{}
			finish.FinishUrl = hdlr.FinishUrl
			finish.FinishUrlTitle = hdlr.FinishUrlTitle
			tmpl.Execute(w, finish)
		}
	}
}

func (hdlr *SubscriptionFrontendHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	subscriptionServiceUrl := hdlr.SubscriptionServiceUrl+"/healthz"
	if subResp, err := http.Get(subscriptionServiceUrl); err == nil {
		if subResp.StatusCode == http.StatusOK {
			procurementUrl := hdlr.CloudCommerceProcurementUrl +  "/providers/" +  hdlr.PartnerId + "/accounts/"

			if client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform"); clientErr != nil {
				LogE.Printf("Healthz failed. Failed to create oath2 client for the procurement API %#v \n", clientErr)
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				if procResp, err := client.Get(procurementUrl); nil != err {
					LogE.Printf("Healthz failed. Cloud Commerce API check failed: %s %#v \n", procurementUrl, err)
					http.Error(w,err.Error(),procResp.StatusCode)
				} else {
					w.WriteHeader(procResp.StatusCode)
				}
			}
		}
	} else {
		LogE.Printf("Healthz failed. Subscription Service check failed: %s %s %#v \n", subscriptionServiceUrl,subResp.StatusCode,err)
		http.Error(w,err.Error(),subResp.StatusCode)
	}
}

func createProduct(prod string, accountId string, subscriptionServiceUrl string, w http.ResponseWriter) bool {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc).String()
	product := Product {
		Id: accountId+"-"+prod,
		Name: accountId+"-"+prod,
		Product: prod,
		Plan: prod,
		Account: accountId,
		State: "ENTITLEMENT_PROD_ACTIVE",
		CreateTime: now,
		UpdateTime: now,
	}
	if productBytes, err := json.Marshal(product); err != nil {
		fmt.Fprintf(w, `{"error": "unable to decode product object %s"}`, err)
		return false
	} else {
		if prodReq, err := http.NewRequest(http.MethodPut, subscriptionServiceUrl+"/entitlements", bytes.NewBuffer(productBytes)); nil != err {
			w.WriteHeader(500)
			fmt.Fprint(w, `{"error": "error with creating product upsert request %s"}`, err)
			return false
		} else {
			if prodResp, err := http.DefaultClient.Do(prodReq); nil != err {
				w.WriteHeader(prodResp.StatusCode)
				fmt.Fprintf(w, `{"error": "error received from subscription service %s"}`, err)
				return false
			} else {
				defer prodResp.Body.Close()
				return true
			}
		}
	}
}

func createContact(contact Contact, subscriptionServiceUrl string, w http.ResponseWriter) bool {
	//submit to subscript service
	if contactBytes, err := json.Marshal(contact); err != nil {
		LogE.Println(err)
		return false
	} else {
		if contactReq, err := http.NewRequest(http.MethodPut, subscriptionServiceUrl+"/contacts", bytes.NewBuffer(contactBytes)); nil != err {
			w.WriteHeader(500)
			fmt.Fprint(w, `{"error": "error with creating contact upsert request %s"}`, err)
			return false
		} else {
			if contactResp, err := http.DefaultClient.Do(contactReq); nil != err {
				w.WriteHeader(contactResp.StatusCode)
				fmt.Fprintf(w, `{"error": "error received from subscription service %s"}`, err)
				return false
			} else {
				defer contactResp.Body.Close()
				return true
			}
		}
	}
}

func createAccount(account Account, subscriptionServiceUrl string, w http.ResponseWriter) bool {
	//submit to subscript service
	if accountBytes, err := json.Marshal(account); err != nil {
		LogE.Println(err)
		return false
	} else {
		if accountReq, err := http.NewRequest(http.MethodPut, subscriptionServiceUrl+"/accounts", bytes.NewBuffer(accountBytes)); nil != err {
			w.WriteHeader(500)
			fmt.Fprint(w, `{"error": "error with creating account upsert request %s"}`, err)
			return false
		} else {
			if accountResp, err := http.DefaultClient.Do(accountReq); nil != err {
				w.WriteHeader(accountResp.StatusCode)
				fmt.Fprintf(w, `{"error": "error received from subscription service %s"}`, err)
				return false
			} else {
				defer accountResp.Body.Close()
				return true
			}
		}
	}
}

func postAccountApproval(cloudCommerceProcurementUrl string, partnerId string,accountName string, w http.ResponseWriter) error {
	procurementUrl := cloudCommerceProcurementUrl +  "/providers/" +  partnerId + "/accounts/" + accountName + ":approve"
	jsonApproval := []byte(`
			{
				"approvalName": "signup"
			}`)
	if client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform"); clientErr != nil {
		LogE.Printf("Failed to create oath2 client for the procurement API %#v \n", clientErr)
		return clientErr
	} else {
		LogI.Printf("Sending account approval: %s \n", procurementUrl)
		if procResp, err := client.Post(procurementUrl, "", bytes.NewBuffer(jsonApproval)); nil != err {
			LogE.Printf("Failed sending entitlement approval request %s %#v \n", procurementUrl, err)
			return err
		} else {
			defer procResp.Body.Close()
			if procResp.StatusCode != 200 {
				LogE.Println("Account approval received error response: ", procResp.StatusCode)
				responseDump, _ := httputil.DumpResponse(procResp, true)
				LogE.Println(string(responseDump))
				return errors.New("Account approval received error response: " + procResp.Status)
			} else {
				LogI.Printf("%s %s", procurementUrl, procResp.Status)
			}
			return nil
		}
	}
}

func postAccountReset(cloudCommerceProcurementUrl string, partnerId string,accountName string,w http.ResponseWriter) error {
	procurementUrl := cloudCommerceProcurementUrl +  "/providers/" +  partnerId + "/accounts/" + accountName + ":reset"
	if client, clientErr := google.DefaultClient(oauth2.NoContext,"https://www.googleapis.com/auth/cloud-platform"); clientErr != nil {
		LogE.Printf("Failed to create oath2 client for the procurement API %#v \n", clientErr)
		return clientErr
	} else {
		LogI.Printf("Sending account reset: %s \n", procurementUrl)
		if procResp, err := client.Post(procurementUrl, "", nil); nil != err {
			LogE.Printf("Failed sending account reset request %s %#v \n", procurementUrl, err)
			return err
		} else {
			defer procResp.Body.Close()
			if procResp.StatusCode != 200 {
				LogE.Println("Account reset received error response: ", procResp.StatusCode)
				responseDump, _ := httputil.DumpResponse(procResp, true)
				LogE.Println(string(responseDump))
				return errors.New("Account reset received error response: " + procResp.Status)
			}
			return nil
		}
	}
}

func accountValid(accountId string, prod string) (bool, error) {
	subscriptonsUrl := "https://cloudbilling.googleapis.com/v1/partnerSubscriptions?externalAccountId="+accountId
	if client, clientErr := google.DefaultClient(oauth2.NoContext, "https://www.googleapis.com/auth/cloud-platform"); clientErr != nil {
		LogE.Printf("Failed to create oath2 client for the procurement API %#v \n", clientErr)
		return false, clientErr
	} else {
		LogI.Printf("Validating account : %s \n", subscriptonsUrl)
		if subResp, err := client.Get(subscriptonsUrl); nil != err {
			LogE.Printf("Failed validating account request %s %#v \n", subscriptonsUrl, err)
			return false, err
		} else {
			defer subResp.Body.Close()
			if subResp.StatusCode != 200 {
				LogE.Println("Account check received error response: ", subResp.StatusCode)
				responseDump, _ := httputil.DumpResponse(subResp, true)
				LogE.Println(string(responseDump))
				return false,errors.New("Account check received error response: " + subResp.Status)
			}
			responseDump, _ := httputil.DumpResponse(subResp, true)
			responseString := string(responseDump)
			if strings.Contains(responseString,prod) {
				return true, nil
			} else {
				return false, nil
			}
		}
	}
}






