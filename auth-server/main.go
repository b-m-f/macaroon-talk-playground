package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
		type macaroon_json struct {
			macaroon string
		}
		decoder := json.NewDecoder(r.Body)
		//"gopkg.in/macaroon.v2"
		fmt.Printf("%v", r.Body)
		var t macaroon_json
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	}
}
