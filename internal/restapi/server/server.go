package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ashwiniag/goKakashi/internal/restapi/v0/scan"
	"github.com/ashwiniag/goKakashi/pkg/config/v0"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
	"github.com/swaggest/usecase"
)

type Server struct {
	AuthToken string
	Websites  map[string]config.Website

	Port int
}

func (srv *Server) Service() *web.Service {
	s := web.NewService(openapi31.NewReflector())

	s.OpenAPISchema().SetTitle("GoKakashi API")
	s.OpenAPISchema().SetDescription("This is the GoKakashi REST API.")
	s.OpenAPISchema().SetVersion("v0.0.1")

	bearerAuth := &BearerTokenAuth{AuthToken: srv.AuthToken}
	s.Wrap(bearerAuth.Middleware)

	s.Post("/api/v0/scan", usecase.NewInteractor(scan.Post))
	s.Get("/api/v0/scan/{scan_id}", usecase.NewInteractor(scan.Get))

	s.Docs("/docs", swgui.New)
	return s
}

func (s *Server) Serve() {
	err := http.ListenAndServe(":"+strconv.Itoa(s.Port), s.Service())
	if err != nil {
		log.Println("Error starting up the server", err)
		return
	}
}
