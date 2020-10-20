package main

import (
	"io/ioutil"
	"net/url"
	"strings"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
)

// Billing list of services to pay
type Billing struct {
	Services []string `json:"services"`
}

// BillingError error response
type BillingError struct {
	Error string `json:"error"`
}

// Token inspect response
type TokenIntrospec struct {
	Exp	int	`json:"exp"`
	Nbf	int	`json:"nbf"`
	Iat	int	`json:"iat"`
	Jti	string	`json:"jti"`
	Aud	string	`json:"aud"`
	Typ	string	`json:"typ"`
	Acr	string	`json:"acr"`
	Active	bool	`json:"active"`
  }  

var config = struct {
	tokenIntrospection string
}{
	tokenIntrospection: "http://192.168.2.10:8080/auth/realms/learningApp/protocol/openid-connect/token/introspect",
}

func main() {
	http.HandleFunc("/billing/v1/services", enableLog(services))
	http.ListenAndServe(":8082", nil)
}

func services(w http.ResponseWriter, r *http.Request) {
	token, err := getToken(r)
	if err != nil {
		log.Println(err)
		s := &BillingError{Error: err.Error()}
		encoder := json.NewEncoder(w)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(s)
		return
	}
	log.Println("Token: ", token)
	// Validate token
	if !validateToken(token) {
		log.Println(err)
		s := &BillingError{Error: "InvalidToken"}
		encoder := json.NewEncoder(w)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(s)
		return
	}

	s := Billing{
		Services: []string{
			"electronic",
			"phone",
			"internet",
			"water",
		},
	}

	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	encoder.Encode(s)
}

func enableLog(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		log.SetPrefix(handlerName + " ")
		log.Println("--> " + handlerName)
		log.Printf("request: %+v\n", r.RequestURI)
		handler(w, r)
		log.Println("<--" + handlerName + "\n")
	}
}

func getToken(r *http.Request) (string, error) {
	// header
	token := r.Header.Get("Authorization")
	if token != "" {
		parts := strings.Split(token, " ")
		if len(parts) != 2 {
			return "", fmt.Errorf("Invalid Authorization header format")
		}
		return parts[1], nil
	}

	// form body
	token = r.FormValue("access_token")
	if token != "" {
		return token, nil
	}

	// query
	token = r.URL.Query().Get("access_token")
	if token != "" {
		return token, nil
	}

	return token, fmt.Errorf("Missing access token")
}

func validateToken(token string) bool {
	// request
	form := url.Values{}
	form.Add("token_type_hint", "requesting_party_token")
	form.Add("token", token)
	req, err := http.NewRequest("POST", config.tokenIntrospection, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println(err)
		return false
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("tokenChecker", "79198687-c9eb-4291-a53a-682cc673f870")

	// client
	c := http.Client{}
	res,err := c.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}

	// process response
	if res.StatusCode != 200 {
		log.Println("Status is not 200: ", res.StatusCode)
		return false
	}

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	defer res.Body.Close()

	introSpect := &TokenIntrospec{}
	err = json.Unmarshal(byteBody, introSpect)
	if err != nil {
		log.Println(err)
		return false
	}

	return introSpect.Active
}
