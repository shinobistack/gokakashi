package v1

import (
	"encoding/json"
	"fmt"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/shinobistack/gokakashi/ent"
	configv1 "github.com/shinobistack/gokakashi/internal/config/v1"
	"github.com/shinobistack/gokakashi/internal/restapi/server/middleware"
	agentlabels1 "github.com/shinobistack/gokakashi/internal/restapi/v1/agentlabels"
	agents1 "github.com/shinobistack/gokakashi/internal/restapi/v1/agents"
	agenttasks1 "github.com/shinobistack/gokakashi/internal/restapi/v1/agenttasks"
	"github.com/shinobistack/gokakashi/webapp"

	integrations1 "github.com/shinobistack/gokakashi/internal/restapi/v1/integrations"
	integrationtype1 "github.com/shinobistack/gokakashi/internal/restapi/v1/integrationtype"
	policies1 "github.com/shinobistack/gokakashi/internal/restapi/v1/policies"
	policylabels1 "github.com/shinobistack/gokakashi/internal/restapi/v1/policylabels"
	scanlabels1 "github.com/shinobistack/gokakashi/internal/restapi/v1/scanlabels"
	scannotify1 "github.com/shinobistack/gokakashi/internal/restapi/v1/scannotify"
	scans1 "github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
	v2agents "github.com/shinobistack/gokakashi/internal/restapi/v2/agents"

	"log"
	"net/http"
	"strconv"

	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"
	swg "github.com/swaggest/swgui"
	swgui "github.com/swaggest/swgui/v5emb"
	"github.com/swaggest/usecase"
)

type Server struct {
	AuthToken string
	Websites  string
	Port      int
	DB        *ent.Client
	DBConfig  configv1.DbConnection
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

	v2Reflector := openapi31.NewReflector()
	v2Reflector.SpecEns().WithServers(
		openapi31.Server{URL: "/api/v2/"},
	).WithInfo(openapi31.Info{Title: "GoKakashi API v2"})
	apiV2 := web.NewService(v2Reflector)

	bearerAuth := &middleware.BearerTokenAuth{AuthToken: srv.AuthToken}
	apiV1.Wrap(bearerAuth.Middleware)
	apiV1.Wrap(middleware.NewRequestLogger().Middleware)

	apiV2.Wrap(bearerAuth.Middleware)
	apiV2.Wrap(middleware.NewRequestLogger().Middleware)

	// Define API endpoints
	apiV1.Get("/integrations", usecase.NewInteractor(integrations1.ListIntegrations(srv.DB)))
	apiV1.Get("/integrations/{id}", usecase.NewInteractor(integrations1.GetIntegration(srv.DB)))
	apiV1.Post("/integrations", usecase.NewInteractor(integrations1.CreateIntegration(srv.DB)))
	apiV1.Put("/integrations/{id}", usecase.NewInteractor(integrations1.UpdateIntegration(srv.DB)))

	apiV1.Get("/integrationtypes", usecase.NewInteractor(integrationtype1.ListIntegrationType(srv.DB)))
	apiV1.Get("/integrationtypes/{id}", usecase.NewInteractor(integrationtype1.GetIntegrationType(srv.DB)))
	apiV1.Post("/integrationtypes", usecase.NewInteractor(integrationtype1.CreateIntegrationType(srv.DB)))
	apiV1.Put("/integrationtypes/{id}", usecase.NewInteractor(integrationtype1.UpdateIntegrationType(srv.DB)))

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

	apiV1.Post("/scans", usecase.NewInteractor(scans1.CreateScan(srv.DB)))
	apiV1.Get("/scans", usecase.NewInteractor(scans1.ListScans(srv.DB)))
	apiV1.Get("/scans/{id}", usecase.NewInteractor(scans1.GetScan(srv.DB)))
	apiV1.Put("/scans/{id}", usecase.NewInteractor(scans1.UpdateScan(srv.DB)))
	apiV1.Delete("/scans/{id}", usecase.NewInteractor(scans1.DeleteScan(srv.DB)))

	apiV1.Post("/scans/{scan_id}/labels", usecase.NewInteractor(scanlabels1.CreateScanLabel(srv.DB)))
	apiV1.Get("/scans/{scan_id}/labels", usecase.NewInteractor(scanlabels1.ListScanLabels(srv.DB)))
	apiV1.Get("/scans/{scan_id}/labels/{key}", usecase.NewInteractor(scanlabels1.GetScanLabel(srv.DB)))
	apiV1.Put("/scans/{scan_id}/labels", usecase.NewInteractor(scanlabels1.UpdateScanLabel(srv.DB)))
	apiV1.Delete("/scans/{scan_id}/labels/{key}", usecase.NewInteractor(scanlabels1.DeleteScanLabel(srv.DB)))

	apiV1.Post("/agents", usecase.NewInteractor(agents1.RegisterAgent(srv.DB)))
	apiV1.Get("/agents", usecase.NewInteractor(agents1.PollAgents(srv.DB)))
	apiV1.Get("/agents/{id}", usecase.NewInteractor(agents1.GetAgent(srv.DB)))
	apiV1.Put("/agents/{id}", usecase.NewInteractor(agents1.UpdateAgent(srv.DB)))
	apiV1.Put("/agents/{id}/heartbeat", usecase.NewInteractor(agents1.UpdateAgentHeartbeat(srv.DB)))
	apiV1.Delete("/agents", usecase.NewInteractor(agents1.DeleteAgent(srv.DB)))

	apiV1.Post("/agents/{agent_id}/tasks", usecase.NewInteractor(agenttasks1.CreateAgentTask(srv.DB)))
	apiV1.Get("/agents/tasks", usecase.NewInteractor(agenttasks1.ListAgentTasks(srv.DB)))
	apiV1.Get("/agents/{agent_id}/tasks", usecase.NewInteractor(agenttasks1.ListAgentTasksByAgentID(srv.DB)))
	apiV1.Get("/agents/{agent_id}/tasks/{id}", usecase.NewInteractor(agenttasks1.GetAgentTask(srv.DB)))
	apiV1.Put("/agents/{agent_id}/tasks/{id}", usecase.NewInteractor(agenttasks1.UpdateAgentTask(srv.DB)))
	apiV1.Delete("/agents/{agent_id}/tasks/{id}", usecase.NewInteractor(agenttasks1.DeleteAgentTask(srv.DB)))

	apiV1.Post("/agents/{agent_id}/labels", usecase.NewInteractor(agentlabels1.CreateAgentLabel(srv.DB)))
	apiV1.Get("/agents/{agent_id}/labels", usecase.NewInteractor(agentlabels1.ListAgentLabels(srv.DB)))
	apiV1.Get("/agents/{agent_id}/labels/{key}", usecase.NewInteractor(agentlabels1.GetAgentLabel(srv.DB)))
	apiV1.Put("/agents/{agent_id}/labels", usecase.NewInteractor(agentlabels1.UpdateAgentLabel(srv.DB)))
	apiV1.Delete("/agents/{agent_id}/labels/{key}", usecase.NewInteractor(agentlabels1.DeleteAgentLabel(srv.DB)))

	apiV1.Post("/scannotify", usecase.NewInteractor(scannotify1.CreateScanNotify(srv.DB)))
	apiV1.Get("/scannotify", usecase.NewInteractor(scannotify1.GetScanNotify(srv.DB)))
	apiV1.Put("/scannotify/{scan_id}", usecase.NewInteractor(scannotify1.UpdateScanNotify(srv.DB)))
	apiV1.Delete("/scannotify/{scan_id}", usecase.NewInteractor(scannotify1.DeleteScanNotify(srv.DB)))

	apiV2.Post("/agents", usecase.NewInteractor(v2agents.Register(srv.DB)))
	apiV2.Patch("/agents/{agent_id}/heartbeat", usecase.NewInteractor(v2agents.Heartbeat(srv.DB)))

	s.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5555"},
		AllowedMethods: []string{http.MethodOptions, http.MethodGet},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler)
	s.Mount("/api/v1/openapi.json", specHandler(apiV1.OpenAPICollector.SpecSchema().(*openapi31.Spec)))
	s.Mount("/api/v1", apiV1)
	s.Mount("/api/v2/openapi.json", specHandler(apiV2.OpenAPICollector.SpecSchema().(*openapi31.Spec)))
	s.Mount("/api/v2", apiV2)

	s.Docs("/docs", swgui.NewWithConfig(swg.Config{
		ShowTopBar: true,
		SettingsUI: map[string]string{
			"urls": `[
	{"url": "/api/v1/openapi.json", "name": "APIv1"},
	{"url": "/api/v2/openapi.json", "name": "APIv2"}
]`,
			`"urls.primaryName"`: `"APIv1"`,
		},
	}))

	frontend, err := webapp.ReactApp()
	if err != nil {
		return nil
	}
	s.Mount("/", frontend)

	return s
}

// InitDB defaults to postgres
func InitDB(dbConfig configv1.DbConnection) *ent.Client {
	// Build the database connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
	)

	client, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
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
	s.DB = InitDB(s.DBConfig)
	defer s.DB.Close()

	// Start the server
	err := http.ListenAndServe(":"+strconv.Itoa(s.Port), s.Service())
	if err != nil {
		log.Println("Error starting up the server", err)
		return
	}
}
