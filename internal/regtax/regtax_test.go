package regtax

import (
	"strings"
	"testing"
)

func TestCalculate_KnowledgeExamples(t *testing.T) {
	cases := []struct {
		name         string
		in           Input
		wantTaxable  int
		wantTax      int
		wantStmtPart string // ExemptStatementsに含まれる文言(空なら検査しない)
	}{
		{
			"例A: 土地+建物", // 71,100
			Input{Properties: []Property{
				{Kind: "land", Value: 12_345_678},
				{Kind: "building", Value: 5_432_100},
			}},
			17_777_000, 71_100, "",
		},
		{
			"例B: 複数土地", // 31,800
			Input{Properties: []Property{
				{Kind: "land", Value: 4_478_400},
				{Kind: "land", Value: 3_489_100},
			}},
			7_967_000, 31_800, "",
		},
		{
			"例C: 私道(持分1/7)", // 1,000
			Input{Properties: []Property{
				{Kind: "land", ShareNum: 1, ShareDen: 7, PrivateRoad: &PrivateRoad{NeighborUnitPrice: 125_000, Area: 50}},
			}},
			267_000, 1_000, "",
		},
		{
			"例D: 100万円以下の土地(措置2)", // 非課税
			Input{Properties: []Property{
				{Kind: "land", Value: 850_000, Exemption: "small_value"},
			}},
			0, 0, "第84条の2の2第2項",
		},
		{
			"例E: 最低税額(免税対象外)", // 1,000
			Input{Properties: []Property{
				{Kind: "land", Value: 90_000, Exemption: "none"},
			}},
			90_000, 1_000, "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := Calculate(c.in)
			if r.TaxableTotal != c.wantTaxable {
				t.Errorf("TaxableTotal = %d, want %d", r.TaxableTotal, c.wantTaxable)
			}
			if r.Tax != c.wantTax {
				t.Errorf("Tax = %d, want %d", r.Tax, c.wantTax)
			}
			if c.wantStmtPart != "" && !strings.Contains(strings.Join(r.ExemptStatements, " "), c.wantStmtPart) {
				t.Errorf("ExemptStatements = %v, want contains %q", r.ExemptStatements, c.wantStmtPart)
			}
		})
	}
}

// 100万円以下の土地は自動免税せず、注意のみ(税額は課税)
func TestCalculate_SmallValueEligibilityNoteOnly(t *testing.T) {
	r := Calculate(Input{Properties: []Property{{Kind: "land", Value: 90_000}}})
	if r.Tax != 1_000 {
		t.Errorf("must not auto-exempt: Tax = %d, want 1000", r.Tax)
	}
	if len(r.EligibilityNotes) != 1 {
		t.Errorf("expected 1 small-value eligibility note, got %v", r.EligibilityNotes)
	}
}

// 数次相続の中間者免税(措置1)
func TestCalculate_Intermediate(t *testing.T) {
	r := Calculate(Input{Properties: []Property{{Kind: "land", Value: 5_000_000, Exemption: "intermediate"}}})
	if r.Tax != 0 {
		t.Errorf("intermediate exemption should give tax 0, got %d", r.Tax)
	}
	if !strings.Contains(strings.Join(r.ExemptStatements, " "), "第84条の2の2第1項") {
		t.Errorf("missing intermediate-exemption statement: %v", r.ExemptStatements)
	}
}
