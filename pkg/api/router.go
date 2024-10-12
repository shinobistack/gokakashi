package api

import (
	"github.com/ashwiniag/goKakashi/pkg/config"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func StartAPIServer(port int, websites map[string]config.Website, validToken string) {
	r := mux.NewRouter()
	// ToDo: to restructure and parse the api tokens/authentication provided from websites config

	// Wrap StartScan handler to pass the config
	r.Handle("/api/v0/scan", BearerTokenAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StartScan(w, r, websites)
	}), validToken)).Methods("POST")

	// Wrap StatusHandler similarly
	r.Handle("/api/v0/scan/{scan_id}/status", BearerTokenAuth(http.HandlerFunc(StatusHandler), validToken)).Methods("GET")

	http.ListenAndServe(":"+strconv.Itoa(port), r)
}
