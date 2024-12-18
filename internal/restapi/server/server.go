package server

import (
	"context"
	"entgo.io/ent/dialect"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/internal/restapi/v0/integrations"
	"github.com/swaggest/usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
)

type Server struct {
	AuthToken string
	Websites  string
	Port      int
	DB        *ent.Client
}

func (srv *Server) Service() *web.Service {
	s := web.NewService(openapi31.NewReflector())

	s.OpenAPISchema().SetTitle("GoKakashi API")
	s.OpenAPISchema().SetDescription("This is the GoKakashi REST API.")
	s.OpenAPISchema().SetVersion("v0.0.1")

	bearerAuth := &BearerTokenAuth{AuthToken: srv.AuthToken}
	s.Wrap(bearerAuth.Middleware)

	// Define API endpoints
	s.Get("/api/v0/integrations", usecase.NewInteractor(integrations.ListIntegrations(srv.DB)))
	s.Get("/api/v0/integration/{id}", usecase.NewInteractor(integrations.GetIntegration(srv.DB)))
	s.Post("/api/v0/integration/create", usecase.NewInteractor(integrations.CreateIntegration(srv.DB)))
	s.Put("/api/v0/integration/update/{id}", usecase.NewInteractor(integrations.UpdateIntegration(srv.DB)))

	//s.Post("/api/v0/scan", usecase.NewInteractor(scan.Post))
	//s.Get("/api/v0/scan/{scan_id}", usecase.NewInteractor(scan.Get))

	s.Docs("/docs", swgui.New)
	return s
}

// InitDB defaults to postgres
func InitDB() *ent.Client {
	fmt.Println("test: inside InitDB")
	// ToDo: To take DB connection as input
	client, err := ent.Open(dialect.Postgres, "host=localhost port=5432 user=postgres password=secret dbname=postgres sslmode=disable")

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automatically run migrations
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Failed to create database schema: %v", err)
	}

	log.Println("Database initialized successfully")
	return client

}

func (s *Server) Serve() {
	// Initialize the database client
	s.DB = InitDB()
	defer s.DB.Close()

	// Start the server
	err := http.ListenAndServe(":"+strconv.Itoa(s.Port), s.Service())
	if err != nil {
		log.Println("Error starting up the server", err)
		return
	}
}
