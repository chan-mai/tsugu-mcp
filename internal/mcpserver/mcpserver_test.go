package mcpserver

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

// 入力スキーマ推論がpanicしないこと
func TestNewServer_SchemaInference(t *testing.T) {
	if newServer() == nil {
		t.Fatal("newServer returned nil")
	}
}

func TestHandleChart_Generates(t *testing.T) {
	out := filepath.Join(t.TempDir(), "chart.pdf")
	in := chartToolInput{
		Document: chartDoc{
			Decedent: chartDecedent{Name: "甲", DeathDate: "2025-01-01"},
			Spouse:   &chartPerson{Name: "乙", Relationship: "妻"},
			Children: []chartChild{{Name: "丙", Relationship: "長男", Outcome: "inherit"}},
		},
		OutputPath: out,
	}
	res, structured, err := handleChart(context.Background(), nil, in)
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error result: %+v", res.Content)
	}
	assertPDF(t, structured.Path)
}

func TestHandleTouki_GeneratesAndTemp(t *testing.T) {
	// outputPath省略時は一時ファイル
	in := toukiToolInput{
		Document: toukiDoc{
			Decedent:        toukiDecedent{Name: "甲"},
			Applicants:      []toukiApplicant{{Name: "乙", Contact: true}},
			ApplicationDate: "2026-06-26",
			Registry:        "東京法務局",
			Properties: []toukiProperty{
				{Kind: "land", Location: "東京都…", LotNumber: "1番", LandCategory: "宅地", Area: "100.00"},
			},
		},
	}
	res, structured, err := handleTouki(context.Background(), nil, in)
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error result: %+v", res.Content)
	}
	if filepath.Ext(structured.Path) != ".pdf" {
		t.Errorf("temp file is not .pdf: %s", structured.Path)
	}
	defer os.Remove(structured.Path)
	assertPDF(t, structured.Path)
}

func TestHandleTouki_InvalidIsError(t *testing.T) {
	// 申請人・不動産なし(Validateエラー)はIsError
	res, _, err := handleTouki(context.Background(), nil, toukiToolInput{
		Document: toukiDoc{Registry: "東京法務局"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError {
		t.Fatal("validation error should set IsError")
	}
}

func assertPDF(t *testing.T, path string) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read output PDF: %v", err)
	}
	if !bytes.HasPrefix(b, []byte("%PDF")) {
		t.Errorf("not a PDF: %s", path)
	}
}
