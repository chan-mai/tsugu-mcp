package casepattern

import (
	"strings"
	"testing"
)

func TestSelect_Primary(t *testing.T) {
	cases := []struct {
		in      Input
		wantKey string
		wantTax string
	}{
		{Input{Method: "legal"}, "A", "4/1000"},
		{Input{Method: "agreement"}, "B", "4/1000"},
		{Input{Method: "will_specified"}, "C", "4/1000"},
		{Input{Method: "bequest", BequestToHeir: true}, "D-1", "4/1000"},
		{Input{Method: "bequest", BequestToHeir: false}, "D-2", "20/1000"},
	}
	for _, c := range cases {
		r := Select(c.in)
		if r.Primary.Key != c.wantKey {
			t.Errorf("%+v: key = %q, want %q", c.in, r.Primary.Key, c.wantKey)
		}
		if r.Primary.TaxRate != c.wantTax {
			t.Errorf("%+v: tax = %q, want %q", c.in, r.Primary.TaxRate, c.wantTax)
		}
	}
}

func TestSelect_Modifiers(t *testing.T) {
	r := Select(Input{Method: "agreement", Substitution: true, Renunciation: true})
	joined := strings.Join(r.Modifiers, " / ")
	if !strings.Contains(joined, "代襲") {
		t.Error("substitution must add 代襲 modifier")
	}
	if !strings.Contains(joined, "相続放棄") {
		t.Error("renunciation must add 放棄 modifier")
	}
}

func TestSelect_UnknownMethod(t *testing.T) {
	if Select(Input{Method: "xxx"}).Primary.Key != "?" {
		t.Error("unknown method should yield key '?'")
	}
}
