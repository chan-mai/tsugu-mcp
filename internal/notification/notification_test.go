package notification

import (
	"strings"
	"testing"

	"github.com/chan-mai/tsugu-mcp/ymd"
)

func TestGuide_PreEnforcement(t *testing.T) {
	r := Guide(Input{DeathDate: ymd.Date{Year: 2020, Month: 1, Day: 1}})
	if !strings.HasPrefix(r.Deadline, "2027-03-31") {
		t.Errorf("pre-enforcement deadline = %q, want 2027-03-31", r.Deadline)
	}
}

func TestGuide_PostEnforcement(t *testing.T) {
	r := Guide(Input{DeathDate: ymd.Date{Year: 2025, Month: 6, Day: 15}})
	if !strings.HasPrefix(r.Deadline, "2028-06-15") {
		t.Errorf("post-enforcement deadline = %q, want 2028-06-15", r.Deadline)
	}
}

func TestGuide_Provisional(t *testing.T) {
	r := Guide(Input{DeathDate: ymd.Date{Year: 2025, Month: 6, Day: 15}})
	joined := strings.Join(r.Provisional, " ")
	if !strings.Contains(joined, "非課税") || !strings.Contains(joined, "76条の3") {
		t.Errorf("provisional missing key facts: %v", r.Provisional)
	}
}
