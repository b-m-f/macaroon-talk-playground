package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/macaroon.v2"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/macaroon", Macaroon)
	fmt.Printf("Starting server for assets at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
func Macaroon(w http.ResponseWriter, r *http.Request) {
	root_macaroon, err := macaroon.New([]byte("AliceKey"), []byte("root"), "http://localhost:8080", macaroon.LatestVersion)
	root_macaroon_marshal, err := root_macaroon.MarshalBinary()

	if err != nil {
		fmt.Fprintf(w, "macaroon, %q", err)
	}

	fmt.Fprintf(w, "macaroon, %q", root_macaroon_marshal)
}
