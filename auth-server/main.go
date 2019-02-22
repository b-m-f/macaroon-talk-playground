package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/macaroon.v2"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", Authenticate)
	fmt.Printf("Starting server for assets at 9999")
	log.Fatal(http.ListenAndServe(":9999", router))
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		return
	} else {
		type MacaroonJSON struct {
			Macaroon string
		}
		var requestData MacaroonJSON
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			panic(err)
		}
		log.Print("\n")
		log.Print(requestData)

		receivedRootMacaroonObject, err := macaroon.New([]byte("test"), []byte("test"), "", 2)

		// here we use the previously hardcoded value from the asset server,
		// which in a real use case has to be previously agreed
		// between these two services by using this services public key
		// or another mechanism
		newAuthMacaroonObject, err := macaroon.New([]byte("Alice3rdKey"), []byte("Auth"), "http://localhost:9999", 2)

		receivedRootMacaroonObject.UnmarshalJSON([]byte(requestData.Macaroon))

		newAuthMacaroonObject.Bind(receivedRootMacaroonObject.Signature())
		newAuthJSON, err := newAuthMacaroonObject.MarshalJSON()

		if err != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			fmt.Fprintf(w, "Error:, %q", err)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"macaroon": string(newAuthJSON)})
		}
	}
}
