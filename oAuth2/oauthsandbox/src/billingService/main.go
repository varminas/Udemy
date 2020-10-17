package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"runtime"
)

// Billing list of services to pay
type Billing struct {
	Services []string `json:"services"`
}

func main() {
	http.HandleFunc("/billing/v1/services", enableLog(services))
	http.ListenAndServe(":8082", nil)
}

func services(w http.ResponseWriter, r *http.Request) {
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
