package agent

import (
	"testing"
)

func TestStatuses(t *testing.T) {
	expected := []string{"connected", "scan_in_progress", "disconnected"}
	got := Statuses()
	if len(got) != len(expected) {
		t.Fatalf("expected %d statuses, got %d", len(expected), len(got))
	}
	for i, v := range expected {
		if got[i] != v {
			t.Errorf("expected status %q at index %d, got %q", v, i, got[i])
		}
	}
}

func TestValidateStatus(t *testing.T) {
	valid := []string{"connected", "scan_in_progress", "disconnected"}
	for _, s := range valid {
		if err := ValidateStatus(s); err != nil {
			t.Errorf("ValidateStatus(%q) returned error: %v, want nil", s, err)
		}
	}

	invalid := []string{"foo", "", "CONNECTED", "scan-in-progress"}
	for _, s := range invalid {
		if err := ValidateStatus(s); err != ErrInvalidAgentStatus {
			t.Errorf("ValidateStatus(%q) = %v, want ErrInvalidAgentStatus", s, err)
		}
	}
}
