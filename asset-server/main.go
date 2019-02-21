package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/macaroon.v2"
)

// First we create the root macaroon which allows access to all Photos
var root_macaroon *macaroon.Macaroon

func main() {
	var err error
	root_macaroon, err = macaroon.New([]byte("AliceKey"), []byte("root"), "http://localhost:8080", 2)
	if err != nil {
		panic(err)
	}
	root_macaroon.AddFirstPartyCaveat([]byte("photos = all"))
	// add third party caveat
	// this should ideally be done by sending a caveat_key and predicate
	// to the auth server and retrieving an identifier
	// which then gets added together with the caveat_key
	// In this case we use simple cleartext identifier
	// and caveat key as hardcoded values,
	// that we can use in both services
	root_macaroon.AddThirdPartyCaveat([]byte("Alice3rdKey"), []byte("Auth = tops"), "http://localhost:9999")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/get-image", GetImage)
	router.HandleFunc("/macaroon", Macaroon)
	fmt.Printf("Starting server for assets at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func Macaroon(w http.ResponseWriter, r *http.Request) {
	root_macaroon_json, err := root_macaroon.MarshalJSON()
	if err != nil {
		fmt.Fprintf(w, "Error:, %q", err)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"macaroon": string(root_macaroon_json)})
	}
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	var discharges []*macaroon.Macaroon

	err := root_macaroon.Verify([]byte("AliceKey"), func(caveat string) error {
		if caveat != "photos = all" {
			return fmt.Errorf("Verification failed")
		} else {
			return nil
		}
	}, discharges)

	if err != nil {
		fmt.Fprintf(w, "Error:, %q", err)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "Success, %q", "Here is your picture")
	}
}
