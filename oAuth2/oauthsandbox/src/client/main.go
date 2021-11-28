package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
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
	"time"

	"github.com/google/uuid"
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
	servicesEndpoint    string
}{
	authURL:             "http://192.168.56.10:8080/auth/realms/learningApp/protocol/openid-connect/auth",
	logout:              "http://192.168.56.10:8080/auth/realms/learningApp/protocol/openid-connect/logout",
	afterLogoutRedirect: "http://localhost:8081/home",
	appId:               "billingApp",
	appPassword:         "62b6af59-59d6-4a13-a076-80d7a91aaa9f",
	authCodeCallback:    "http://localhost:8081/authCodeRedirect",
	tokenEndpoint:       "http://192.168.56.10:8080/auth/realms/learningApp/protocol/openid-connect/token",
	servicesEndpoint:    "http://localhost:8082/billing/v1/services",
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
	Services     []string
	State        map[string]struct{}
}

func newAppVar() AppVar {
	return AppVar{State: make(map[string]struct{})}
}

var appVar = newAppVar()

func init() {
	log.SetFlags(log.Ltime)
}

func main() {
	fmt.Println("Listing server on port 8081")
	http.HandleFunc("/home", enableLog(home))
	http.HandleFunc("/login", enableLog(login))
	http.HandleFunc("/refreshToken", enableLog(refreshToken))
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

var codeVerifier = "code-challange43128unreserved-._~nge43128dX"

func login(w http.ResponseWriter, r *http.Request) {
	// create a redirect URL for authentication endpoint.
	req, err := http.NewRequest("GET", config.authURL, nil)

	if err != nil {
		log.Print(err)
		return
	}

	qs := url.Values{}
	state := uuid.New().String()
	appVar.State[state] = struct{}{}
	qs.Add("state", state)
	qs.Add("client_id", config.appId)
	qs.Add("response_type", "code")
	qs.Add("redirect_uri", config.authCodeCallback)
	// qs.Add("scope", "evil-service")
	qs.Add("scope", "billingService")

	// PKCE
	codeChallenge := makeCodeChallenge(codeVerifier)
	qs.Add("code_challenge", codeChallenge)
	qs.Add("code_challenge_method", "S256")

	req.URL.RawQuery = qs.Encode()
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}

func authCodeRedirect(w http.ResponseWriter, r *http.Request) {
	appVar.AuthCode = r.URL.Query().Get("code")
	callbackState := r.URL.Query().Get("state")
	if _, ok := appVar.State[callbackState]; !ok {
		fmt.Fprintf(w, "Error")
		return
	}
	delete(appVar.State, callbackState)

	appVar.SessionState = r.URL.Query().Get("session_state")
	r.URL.RawQuery = ""
	fmt.Printf("Request queries: %+v\n", appVar)
	// http.Redirect(w, r, "http://localhost:8081", http.StatusFound)
	// exchange token here

	exchangeToken(w, r)
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
	appVar = newAppVar()
	http.Redirect(w, r, logoutURL.String(), http.StatusFound)
}

func services(w http.ResponseWriter, r *http.Request) {
	// request
	req, err := http.NewRequest("GET", config.servicesEndpoint, nil)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+appVar.AccessToken)

	// client
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}

	// process response
	if res.StatusCode != 200 {
		log.Println(string(byteBody))
		appVar.Services = []string{string(byteBody)}
		tServices.Execute(w, appVar)
		return
	}

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
	form.Add("code_verifier", codeVerifier)

	req, err := http.NewRequest("POST", config.tokenEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		log.Println(err)
		return
	}

	req.SetBasicAuth(config.appId, config.appPassword)

	// Client
	ctx, cancelFunc := context.WithTimeout(context.Background(), 500*time.Millisecond)
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

func refreshToken(w http.ResponseWriter, r *http.Request) {
	// Request
	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", appVar.RefreshToken)
	req, err := http.NewRequest("POST", config.tokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.appId, config.appPassword)

	// client
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}

	// process response
	byteBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	body := &model.AccessTokenResponse{}
	json.Unmarshal(byteBody, body)
	appVar.AccessToken = body.AccessToken
	appVar.SessionState = body.SessionState
	appVar.RefreshToken = body.RefreshToken
	appVar.Scope = body.Scope

	t.Execute(w, appVar)
}

// BASE64URL-ENCODE(SHA256(ASCII(plain)))
func makeCodeChallenge(plain string) string {
	h := sha256.Sum256([]byte(plain))
	hs := base64.RawURLEncoding.EncodeToString(h[:])
	return hs
}
