package assigner_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/shinobistack/gokakashi/internal/assigner"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockScan struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type MockAgent struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func TestAssignTasks(t *testing.T) {
	// Mock data
	mockScans := []MockScan{
		{ID: uuid.New(), Status: "scan_pending"},
		{ID: uuid.New(), Status: "scan_pending"},
		{ID: uuid.New(), Status: "scan_pending"},
		{ID: uuid.New(), Status: "scan_pending"},
		{ID: uuid.New(), Status: "scan_pending"},
		{ID: uuid.New(), Status: "scan_pending"},
		{ID: uuid.New(), Status: "scan_pending"},
		{ID: uuid.New(), Status: "scan_pending"},
	}
	mockAgents := []MockAgent{
		{ID: 1, Status: "connected"},
		{ID: 2, Status: "connected"},
	}

	// Mock server
	scanHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/scans" && r.URL.Query().Get("status") == "scan_pending" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockScans)
		} else if r.URL.Path == "/api/v1/agents" && r.URL.Query().Get("status") == "connected" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockAgents)
		} else if r.URL.Path == "/api/v1/agents/tasks" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	mockServer := httptest.NewServer(http.HandlerFunc(scanHandler))
	defer mockServer.Close()

	// Run the assigner logic
	assigner.AssignTasks(mockServer.URL, 0, "mock-token")

	t.Log("Ensure tasks are assigned to agents in round-robin fashion.")
}

func TestAssignTasksWithNoAgents(t *testing.T) {
	mockScans := []MockScan{
		{ID: uuid.New(), Status: "scan_pending"},
	}

	scanHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/scans" && r.URL.Query().Get("status") == "scan_pending" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockScans)
		} else if r.URL.Path == "/api/v1/agents" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]MockAgent{}) // No agents available
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	mockServer := httptest.NewServer(http.HandlerFunc(scanHandler))
	defer mockServer.Close()

	assigner.AssignTasks(mockServer.URL, 0, "mock-token")

	t.Log("Ensure no assignments are made when no agents are available.")
}

func TestAssignTasksWithNoScans(t *testing.T) {
	mockAgents := []MockAgent{
		{ID: 1, Status: "connected"},
	}

	scanHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/scans" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]MockScan{}) // No scans available
		} else if r.URL.Path == "/api/v1/agents" && r.URL.Query().Get("status") == "connected" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockAgents)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	mockServer := httptest.NewServer(http.HandlerFunc(scanHandler))
	defer mockServer.Close()

	assigner.AssignTasks(mockServer.URL, 0, "mock-token")

	t.Log("Ensure no assignments are made when no scans are pending.")
}
