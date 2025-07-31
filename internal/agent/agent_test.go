package agent

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/shinobistack/gokakashi/pkg/client"
)

// --- Custom RoundTripper to intercept HTTP requests ---
type fakeRoundTripper struct {
	registerCalled  *bool
	heartbeatCalled *bool
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Path == "/api/v2/agents" {
		*f.registerCalled = true
		resp := http.Response{
			StatusCode: 200,
			Body:       http.NoBody,
			Header:     make(http.Header),
		}
		return &resp, nil
	}
	if len(req.URL.Path) > 20 && req.URL.Path[len(req.URL.Path)-10:] == "/heartbeat" {
		*f.heartbeatCalled = true
		resp := http.Response{
			StatusCode: 200,
			Body:       http.NoBody,
			Header:     make(http.Header),
		}
		return &resp, nil
	}
	resp := http.Response{
		StatusCode: 404,
		Body:       http.NoBody,
		Header:     make(http.Header),
	}
	return &resp, nil
}

// --- Test New ---
func TestNewAgent(t *testing.T) {
	mc := &client.Client{}
	a := New(mc)
	if a.client != mc {
		t.Errorf("expected client to be set")
	}
	if a.heartbeatTicker == nil {
		t.Errorf("expected heartbeatTicker to be set")
	}
	if a.done == nil {
		t.Errorf("expected done channel to be set")
	}
}

// --- Test Start ---
func TestAgentStartAndStop(t *testing.T) {
	registerCalled := false
	heartbeatCalled := false
	fakeTransport := &fakeRoundTripper{
		registerCalled:  &registerCalled,
		heartbeatCalled: &heartbeatCalled,
	}
	httpClient := &http.Client{Transport: fakeTransport}
	cli, err := client.New("http://localhost", httpClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	a := &Agent{
		client:          cli,
		heartbeatTicker: time.NewTicker(10 * time.Millisecond),
		done:            make(chan struct{}),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.Listen(ctx)
		if err != nil {
			t.Errorf("Start returned error: %v", err)
		}
	}()

	// Wait for heartbeat to be called
	time.Sleep(20 * time.Millisecond)
	a.Stop()
	wg.Wait()

	if !registerCalled {
		t.Errorf("expected Register to be called")
	}
	if !heartbeatCalled {
		t.Errorf("expected Heartbeat to be called")
	}
}

// --- Test Stop idempotency ---
func TestAgentStopIdempotent(t *testing.T) {
	mc := &client.Client{}
	a := New(mc)
	a.Stop()
	// Should not panic or block if called again
	a.Stop()
}
