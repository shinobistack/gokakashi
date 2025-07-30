package scan

import (
	"reflect"
	"testing"
)

func TestValidScanStatuses(t *testing.T) {
	expected := []Status{
		Pending,
		InProgress,
		Success,
		Error,
	}
	actual := validStatuses()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("validStatuses() = %v, want %v", actual, expected)
	}
}

func TestStatuses(t *testing.T) {
	expected := []string{
		string(Pending),
		string(InProgress),
		string(Success),
		string(Error),
	}
	actual := Statuses()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Statuses() = %v, want %v", actual, expected)
	}
}

func TestValidateStatus(t *testing.T) {
	valid := Statuses()
	for _, s := range valid {
		err := ValidateStatus(s)
		if err != nil {
			t.Errorf("ValidateStatus(%q) = %v, want nil", s, err)
		}
	}

	invalidStatuses := []string{"", "unknown", "pending", "in_progress", "success ", " error"}
	for _, s := range invalidStatuses {
		err := ValidateStatus(s)
		if err != ErrInvalidScanStatus {
			t.Errorf("ValidateStatus(%q) = %v, want ErrInvalidScanStatus", s, err)
		}
	}
}
