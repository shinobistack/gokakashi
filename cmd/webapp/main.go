package main

import (
	"log"
	"net/http"

	"github.com/shinobistack/gokakashi/webapp"
)

// This is a simple server that serves the React app
// It is used for development and testing purposes only

func main() {
	reactApp, err := webapp.ReactApp()
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		reactApp.ServeHTTP(w, r)
	}))
}
