package mcpserver

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestAutoFillTax(t *testing.T) {
	doc := &toukiDoc{Properties: []toukiProperty{
		{Kind: "land", Value: 12_345_678},
		{Kind: "building", Value: 5_432_100},
	}}
	note := autoFillTax(doc)
	if doc.TaxValue != "17,777,000" {
		t.Errorf("TaxValue = %q, want 17,777,000", doc.TaxValue)
	}
	if doc.RegistrationTax != "71,100" {
		t.Errorf("RegistrationTax = %q, want 71,100", doc.RegistrationTax)
	}
	if !strings.Contains(note, "71,100") {
		t.Errorf("note missing amount: %q", note)
	}
}

func TestAutoFillTax_NoValue(t *testing.T) {
	doc := &toukiDoc{Properties: []toukiProperty{{Kind: "land"}}}
	if autoFillTax(doc) != "" {
		t.Error("no value should yield empty note")
	}
	if doc.TaxValue != "" || doc.RegistrationTax != "" {
		t.Error("must not modify doc when no value")
	}
}

func TestAutoFillTax_DoesNotOverwrite(t *testing.T) {
	doc := &toukiDoc{TaxValue: "999", RegistrationTax: "888", Properties: []toukiProperty{{Kind: "land", Value: 12_345_678}}}
	autoFillTax(doc)
	if doc.TaxValue != "999" || doc.RegistrationTax != "888" {
		t.Error("must not overwrite existing values")
	}
}

func TestAutoFillTax_Exempt(t *testing.T) {
	doc := &toukiDoc{Properties: []toukiProperty{{Kind: "land", Value: 850_000, Exemption: "small_value"}}}
	autoFillTax(doc)
	if !strings.Contains(doc.RegistrationTax, "非課税") {
		t.Errorf("RegistrationTax should carry 免税文言: %q", doc.RegistrationTax)
	}
}

func TestHandleTouki_AutoTax(t *testing.T) {
	in := toukiToolInput{Document: toukiDoc{
		Decedent:   toukiDecedent{Name: "山田 太郎"},
		Applicants: []toukiApplicant{{Name: "山田 花子"}},
		Properties: []toukiProperty{{Kind: "land", Location: "東京都千代田区一番町1番", LotNumber: "1番", Area: "100.00", Value: 12_345_678}},
	}}
	res, out, err := handleTouki(context.Background(), nil, in)
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error: %+v", res.Content)
	}
	if out.Path == "" {
		t.Error("no output path")
	}
	var joined string
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			joined += tc.Text
		}
	}
	if !strings.Contains(joined, "登録免許税を自動計算") {
		t.Errorf("missing auto-tax note: %q", joined)
	}
}

func TestGuidePrompt(t *testing.T) {
	r, err := guidePrompt(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Messages) == 0 {
		t.Fatal("no messages")
	}
	tc, ok := r.Messages[0].Content.(*mcp.TextContent)
	if !ok || !strings.Contains(tc.Text, "list_required_documents") {
		t.Error("guide text missing workflow steps")
	}
}
