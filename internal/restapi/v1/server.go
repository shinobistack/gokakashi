package v1

import (
	"context"
	"encoding/json"
	"entgo.io/ent/dialect"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/shinobistack/gokakashi/ent"
	"github.com/shinobistack/gokakashi/internal/restapi/server/middleware"
	integrations1 "github.com/shinobistack/gokakashi/internal/restapi/v1/integrations"
	policies1 "github.com/shinobistack/gokakashi/internal/restapi/v1/policies"
	policylabels1 "github.com/shinobistack/gokakashi/internal/restapi/v1/policylabels"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"
	swg "github.com/swaggest/swgui"
	swgui "github.com/swaggest/swgui/v5emb"
	"github.com/swaggest/usecase"
	"log"
	"net/http"
	"strconv"
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

	v1Reflector := openapi31.NewReflector()
	v1Reflector.SpecEns().WithServers(
		openapi31.Server{URL: "/api/v1/"},
	).WithInfo(openapi31.Info{Title: "GoKakashi API v1"})
	apiV1 := web.NewService(v1Reflector)

	bearerAuth := &middleware.BearerTokenAuth{AuthToken: srv.AuthToken}
	// Auth applied to routers under /api/v1
	apiV1.Wrap(bearerAuth.Middleware)

	// Define API endpoints
	apiV1.Get("/integrations", usecase.NewInteractor(integrations1.ListIntegrations(srv.DB)))
	apiV1.Get("/integrations/{id}", usecase.NewInteractor(integrations1.GetIntegration(srv.DB)))
	apiV1.Post("/integrations", usecase.NewInteractor(integrations1.CreateIntegration(srv.DB)))
	apiV1.Put("/integrations/{id}", usecase.NewInteractor(integrations1.UpdateIntegration(srv.DB)))

	apiV1.Post("/policies", usecase.NewInteractor(policies1.CreatePolicy(srv.DB)))
	apiV1.Get("/policies", usecase.NewInteractor(policies1.ListPolicies(srv.DB)))
	apiV1.Get("/policies/{id}", usecase.NewInteractor(policies1.GetPolicy(srv.DB)))
	apiV1.Put("/policies/{id}", usecase.NewInteractor(policies1.UpdatePolicy(srv.DB)))
	apiV1.Delete("/policies/{id}", usecase.NewInteractor(policies1.DeletePolicy(srv.DB)))

	apiV1.Post("/policies/{policy_id}/labels", usecase.NewInteractor(policylabels1.CreatePolicyLabel(srv.DB)))
	apiV1.Get("/policies/{policy_id}/labels", usecase.NewInteractor(policylabels1.ListPolicyLabels(srv.DB)))
	apiV1.Get("/policies/{policy_id}/labels/{key}", usecase.NewInteractor(policylabels1.GetPolicyLabel(srv.DB)))
	apiV1.Put("/policies/{policy_id}/labels", usecase.NewInteractor(policylabels1.UpdatePolicyLabels(srv.DB)))
	apiV1.Delete("/policies/{policy_id}/labels/{key}", usecase.NewInteractor(policylabels1.DeletePolicyLabel(srv.DB)))

	s.Mount("/api/v1/openapi.json", specHandler(apiV1.OpenAPICollector.SpecSchema().(*openapi31.Spec)))
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

func specHandler(s *openapi31.Spec) http.Handler {
	j, err := json.Marshal(s)
	if err != nil {
		log.Printf("Failed to marshal OpenAPI schema: %v", err)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Failed to load OpenAPI schema", http.StatusInternalServerError)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(j)
	})
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
