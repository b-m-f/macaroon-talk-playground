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

		macaroonObject, err := macaroon.New([]byte("test"), []byte("test"), "", 2)

		macaroonObject.UnmarshalJSON([]byte(requestData.Macaroon))

		marshal, err := macaroonObject.MarshalJSON()

		fmt.Fprintf(w, "Hello, %q", marshal)
	}
}
