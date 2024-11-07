package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/go-chi/chi/v5/middleware"
	"github.com/shinobistack/gokakashi/internal/restapi/server/middleware"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/policies"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"
	swg "github.com/swaggest/swgui"
	swgui "github.com/swaggest/swgui/v5emb"
	"github.com/swaggest/usecase"
)

type Server struct {
	ServerOpts
}

type ServerOpts struct {
	Host      string
	Port      int
	AuthToken string
}

func NewServer(opt *ServerOpts) *Server {
	defaultHost := "0.0.0.0"
	defaultPort := 8080
	defaultAuthToken := "gokakashi-auth-token" // TODO: generate a cryptographically secure string by default

	if opt == nil {
		opt = &ServerOpts{}
	}

	if len(opt.Host) == 0 {
		opt.Host = defaultHost
	}
	if opt.Port == 0 {
		opt.Port = defaultPort
	}
	if len(opt.AuthToken) == 0 {
		opt.AuthToken = defaultAuthToken
	}

	return &Server{*opt}
}

func (srv *Server) Service() *web.Service {
	s := web.NewService(openapi31.NewReflector())

	s.OpenAPISchema().SetTitle("GoKakashi API")
	s.OpenAPISchema().SetDescription("This is the GoKakashi REST API.")
	s.OpenAPISchema().SetVersion("v0.1.0")

	v1r := openapi3.NewReflector()
	v1r.SpecEns().WithServers(openapi3.Server{URL: "/api/v1/"}).WithInfo(openapi3.Info{Title: "GoKakashi API v1"})
	apiV1 := web.NewService(v1r)

	bearerAuth := &middleware.BearerTokenAuth{AuthToken: srv.AuthToken}
	apiV1.Wrap(
		bearerAuth.Middleware,
	)
	apiV1.Post("/policies", usecase.NewInteractor(policies.Post))
	s.Mount("/api/v1/openapi.json", specHandler(apiV1.OpenAPICollector.SpecSchema()))
	s.Mount("/api/v1", apiV1)

	s.Docs("/docs", swgui.NewWithConfig(swg.Config{
		ShowTopBar: true,
		SettingsUI: map[string]string{
			"urls": `[
	{"url": "/api/v1/openapi.json", "name": "APIv1"}
]`,
			`"urls.primaryName"`: `"APIv1"`,
		},
	}))
	return s
}

func (s *Server) PrintInfo() {
	log.Printf("  API docs: http://%s:%d/docs\n", s.Host, s.Port)
	log.Printf("  API Auth token: %s\n", s.AuthToken)
}

func (s *Server) Serve() {
	host := ""
	if len(s.Host) > 0 {
		host = s.Host
	}
	port := 8000
	if s.Port != 0 {
		port = s.Port
	}
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), s.Service())
	if err != nil {
		log.Println("Error starting up the server", err)
		return
	}
}

func specHandler(s openapi.SpecSchema) http.Handler {
	j, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(j)
	})
}
