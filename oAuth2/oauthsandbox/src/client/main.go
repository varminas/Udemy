package main

import (
	"time"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"

	"learn.oauth.client/model"
)

var config = struct {
	authURL             string
	logout              string
	afterLogoutRedirect string
	appId               string
	appPassword         string
	authCodeCallback    string
	tokenEndpoint       string
	servicesEndpoint	string
}{
	authURL:             "http://192.168.2.10:8080/auth/realms/learningApp/protocol/openid-connect/auth",
	logout:              "http://192.168.2.10:8080/auth/realms/learningApp/protocol/openid-connect/logout",
	afterLogoutRedirect: "http://localhost:8081",
	appId:               "billingApp",
	appPassword:         "62b6af59-59d6-4a13-a076-80d7a91aaa9f",
	authCodeCallback:    "http://localhost:8081/authCodeRedirect",
	tokenEndpoint:       "http://192.168.2.10:8080/auth/realms/learningApp/protocol/openid-connect/token",
	servicesEndpoint:	 "http://localhost:8082/billing/v1/services",
}

var t = template.Must(template.ParseFiles("template/index.html"))
var tServices = template.Must(template.ParseFiles("template/index.html", "template/services.html"))

// Application private variables
type AppVar struct {
	AuthCode     string
	SessionState string
	AccessToken  string
	RefreshToken string
	Scope        string
	Services	 []string
}

var appVar = AppVar{}

func init() {
	log.SetFlags(log.Ltime)
}

func main() {
	fmt.Println("Listing server on port 8081")
	http.HandleFunc("/", enableLog(home))
	http.HandleFunc("/login", enableLog(login))
	http.HandleFunc("/exchangeToken", enableLog(exchangeToken))
	http.HandleFunc("/services", enableLog(services))
	http.HandleFunc("/logout", enableLog(logout))
	http.HandleFunc("/authCodeRedirect", enableLog(authCodeRedirect))
	http.ListenAndServe(":8081", nil)
}

func enableLog(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		log.SetPrefix(handlerName + " ")
		log.Println("--> " + handlerName)
		log.Printf("request: %+v\n", r.RequestURI)
		// log.Println("response: %+v\n", w)
		handler(w, r)
		log.Println("<--" + handlerName + "\n")
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	t.Execute(w, appVar)
}

func login(w http.ResponseWriter, r *http.Request) {
	// create a redirect URL for authentication endpoint.
	req, err := http.NewRequest("GET", config.authURL, nil)

	if err != nil {
		log.Print(err)
		return
	}

	qs := url.Values{}
	qs.Add("state", "123")
	qs.Add("client_id", config.appId)
	qs.Add("response_type", "code")
	qs.Add("redirect_uri", config.authCodeCallback)

	req.URL.RawQuery = qs.Encode()
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}

func authCodeRedirect(w http.ResponseWriter, r *http.Request) {
	appVar.AuthCode = r.URL.Query().Get("code")
	appVar.SessionState = r.URL.Query().Get("session_state")
	r.URL.RawQuery = ""
	fmt.Printf("Request queries: %+v\n", appVar)
	http.Redirect(w, r, "http://localhost:8081", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	q := url.Values{}
	q.Add("redirect_uri", config.afterLogoutRedirect)

	// fmt.Printf("logout 1 %s\n", q)

	logoutURL, err := url.Parse(config.logout)
	if err != nil {
		log.Println(err)
	}
	// fmt.Printf("logout 2 %s\n", logoutURL)
	logoutURL.RawQuery = q.Encode()

	// fmt.Printf("logout 3 %s\n\n", logoutURL)
	appVar = AppVar{}
	http.Redirect(w, r, logoutURL.String(), http.StatusFound)
}

func services(w http.ResponseWriter, r *http.Request) {
	// request
	req,err := http.NewRequest("GET", config.servicesEndpoint, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// client
	c := http.Client{}
	res,err := c.Do(req)
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}

	byteBody,err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}

	// process response
	billingResponse := &model.BillingResponse{}
	err = json.Unmarshal(byteBody, billingResponse)
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}
	appVar.Services = billingResponse.Services

	tServices.Execute(w, appVar)
}

func exchangeToken(w http.ResponseWriter, r *http.Request) {
	// Request
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", appVar.AuthCode)
	form.Add("redirect_uri", config.authCodeCallback)
	form.Add("client_id", config.appId)

	req, err := http.NewRequest("POST", config.tokenEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		log.Println(err)
		return
	}

	req.SetBasicAuth(config.appId, config.appPassword)

	// Client
	ctx,cancelFunc := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancelFunc()
	c := http.Client{}
	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		log.Println("Could not get access token", err)
		return
	}

	// Process response
	byteBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Println(err)
		return
	}

	accessTokenResponse := &model.AccessTokenResponse{}
	json.Unmarshal(byteBody, accessTokenResponse)

	appVar.AccessToken = accessTokenResponse.AccessToken
	appVar.Scope = accessTokenResponse.Scope
	appVar.RefreshToken = accessTokenResponse.RefreshToken
	log.Println(string(byteBody))
	t.Execute(w, appVar)
}
