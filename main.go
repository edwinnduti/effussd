package main

import (
	"encoding/json"
	"log"
	"fmt"
	"os"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// constants
const (
	footer = "\n0:Back 00:Home"
)

type Form struct {
	Text string `json:"text"`
	SessionId string `json:"sessionId"`
	ServiceCode string `json:"serviceCode"`
	PhoneNumber string `json:"phoneNumber"`
}

// the main function
func main() {
	// instanciate the router
	router := mux.NewRouter()
	router.HandleFunc("/", RootEndpoint).Methods("GET")
	router.HandleFunc("/ussd", UssdEndPoint).Methods("POST")

	// establish portnumber
	var Port string = os.Getenv("PORT")

	if Port == "" {
		Port = "8040"
	}

	// set server
	n := negroni.Classic()
	n.UseHandler(router)
	server := &http.Server{
		Handler: n, // router if not negroni
		Addr:    ":" + Port,
	}

	// log server output
	log.Printf("Listening on PORT: %s", Port)
	server.ListenAndServe()
}

// / GET ndpoint
func RootEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

// /ussd POST ussd endpoint
func UssdEndPoint(w http.ResponseWriter, r *http.Request) {
	// set header content to type x-www-form-urlencoded
	// Allow all origin to handle cors
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// instantiate form struct
	var form Form

	// parse json
	err := json.NewDecoder(r.Body).Decode(&form)
	Check(err)

	// create text variable
	text := form.Text

	// message to be printed
	var msg string = ""

	// decide based on the text value
	if text == "" || text == "1*00" || text == "1*0" || text == "1*1*00" {
		phoneNo := fmt.Sprintf("CON Welcome %s",form.PhoneNumber)
		msg += phoneNo
		msg += "\n1: Buy data bundles"
		msg += "\n2: Buy calls and sms bundles"

	} else if text == "1" || text == "1*1*0" {
		msg += "CON 1: Daily bundles"
		msg += "\n2: Weekly bundles"
		msg += footer

	} else if text == "1*1" || text == "1*1*1*0" {
		msg += "CON 1: Buy for my number"
		msg += "\n2: Buy for other number"
		msg += footer

	} else {
		msg += "END Invalid response!"
		msg += footer

	}

	// print on the browser
	w.Write([]byte(msg))
}

// Handle error function
func Check(ERR error) {
	if ERR != nil {
		log.Fatal(ERR)
	}
}
