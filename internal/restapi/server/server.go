package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ashwiniag/goKakashi/internal/restapi/v0/scan"
	"github.com/ashwiniag/goKakashi/pkg/config"
	"github.com/gorilla/mux"
)

type Server struct {
	AuthToken string
	Websites  map[string]config.Website

	Port int
}

func (s *Server) Router() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/api/v0/scan", BearerTokenAuth(&scan.PostHandler{Websites: s.Websites}, s.AuthToken)).Methods("POST")
	r.Handle("/api/v0/scan/{scan_id}", BearerTokenAuth(http.HandlerFunc(scan.GetHandler), s.AuthToken)).Methods("GET")

	return r
}

func (s *Server) Serve() {
	err := http.ListenAndServe(":"+strconv.Itoa(s.Port), s.Router())
	if err != nil {
		log.Println("Error starting up the server", err)
		return
	}
}
