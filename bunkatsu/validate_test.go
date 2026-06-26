package bunkatsu

import (
	"strings"
	"testing"

	"tsugu-mcp/touki"
	"tsugu-mcp/ymd"
)

func okAgreement() Agreement {
	return Agreement{
		Decedent: Decedent{Name: "甲", DeathDate: ymd.Date{Year: 2024, Month: 6, Day: 15}},
		Heirs:    []Heir{{Name: "乙"}, {Name: "丙"}},
		Allocations: []Allocation{{
			Acquirers:  []Acquirer{{Name: "乙", Share: "2分の1"}},
			Properties: []touki.Property{{Kind: touki.Land, Location: "東京都", LotNumber: "1番"}},
		}},
	}
}

func TestValidate_OK(t *testing.T) {
	if err := okAgreement().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_Errors(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*Agreement)
		want string
	}{
		{"死亡日なし", func(a *Agreement) { a.Decedent.DeathDate = ymd.Date{} }, "death date is required"},
		{"相続人なし", func(a *Agreement) { a.Heirs = nil }, "at least one required"},
		{"取得者なし", func(a *Agreement) { a.Allocations[0].Acquirers = nil }, "at least one acquirer"},
		{"財産なし", func(a *Agreement) { a.Allocations[0].Properties = nil; a.Allocations[0].Items = nil }, "at least one property or item"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := okAgreement()
			c.mut(&a)
			err := a.Validate()
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("err = %v, want contains %q", err, c.want)
			}
		})
	}
}

func TestCopyCount(t *testing.T) {
	a := okAgreement()
	if a.CopyCount() != 2 {
		t.Errorf("CopyCount = %d, want 2 (heir count)", a.CopyCount())
	}
	a.Copies = 5
	if a.CopyCount() != 5 {
		t.Errorf("CopyCount = %d, want 5", a.CopyCount())
	}
}
