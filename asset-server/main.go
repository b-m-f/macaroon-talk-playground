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
var rootMacaroon *macaroon.Macaroon

func main() {
	var err error
	rootMacaroon, err = macaroon.New([]byte("AliceKey"), []byte("root"), "http://localhost:8080", 2)
	if err != nil {
		panic(err)
	}
	rootMacaroon.AddFirstPartyCaveat([]byte("photos = all"))
	// add third party caveat
	// this should ideally be done by sending a caveat_key and predicate
	// to the auth server and retrieving an identifier
	// which then gets added together with the caveat_key
	// In this case we use simple cleartext identifier
	// and caveat key as hardcoded values,
	// that we can use in both services
	rootMacaroon.AddThirdPartyCaveat([]byte("Alice3rdKey"), []byte("Auth"), "http://localhost:9999")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", _Index)
	router.HandleFunc("/get-image", _GetImage)
	router.HandleFunc("/macaroon", _Macaroon)
	fmt.Printf("Starting server for assets at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func _Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func _Macaroon(w http.ResponseWriter, r *http.Request) {
	rootMacaroonJSON, err := rootMacaroon.MarshalJSON()
	if err != nil {
		fmt.Fprintf(w, "Error:, %q", err)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"macaroon": string(rootMacaroonJSON)})
	}
}

func _GetImage(w http.ResponseWriter, r *http.Request) {
	receivedDischargeMacaroon, err := macaroon.New([]byte("test"), []byte("test"), "", 2)
	if err != nil {
		fmt.Fprintf(w, "Error:, %q", err)
		return
	}
	if r.FormValue("macaroon") != "" {
		dischargeMacaroonJSON := r.FormValue("macaroon")
		receivedDischargeMacaroon.UnmarshalJSON([]byte(dischargeMacaroonJSON))
		log.Print(receivedDischargeMacaroon)
		var discharges = []*macaroon.Macaroon{receivedDischargeMacaroon}

		verificationError := rootMacaroon.Verify([]byte("AliceKey"), func(caveat string) error {
			if caveat != "photos = all" {
				return fmt.Errorf("Verification failed")
			}
			return nil

		}, discharges)

		if verificationError != nil {
			fmt.Fprintf(w, "Error:, %q", verificationError)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			fmt.Fprintf(w, "Success, %q", "Here is your picture")
		}
	} else {
		fmt.Fprintf(w, "Error: Please provide a valid Discharge Macaroon to prove you have access")

	}
}
