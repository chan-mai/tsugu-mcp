package mcpserver

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func toolText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	if len(res.Content) == 0 {
		t.Fatal("no content")
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("not text content: %T", res.Content[0])
	}
	return tc.Text
}

func TestHandleTax_ExampleA(t *testing.T) {
	res, out, err := handleTax(context.Background(), nil, taxToolInput{Properties: []taxProperty{
		{Kind: "land", Value: 12_345_678},
		{Kind: "building", Value: 5_432_100},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error: %+v", res.Content)
	}
	if out.Tax != 71_100 {
		t.Errorf("Tax = %d, want 71100", out.Tax)
	}
	if !strings.Contains(toolText(t, res), "免責") {
		t.Error("missing disclaimer in text")
	}
}

func TestHandleDocs_Siblings(t *testing.T) {
	res, out, err := handleDocs(context.Background(), nil, docToolInput{Method: "legal", HeirPattern: "siblings", ApplicantAtWindow: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error: %+v", res.Content)
	}
	if len(out.Categories) == 0 {
		t.Fatal("no categories")
	}
	text := toolText(t, res)
	if !strings.Contains(text, "父・母") {
		t.Error("missing parents' register")
	}
	if !strings.Contains(text, "免責") {
		t.Error("missing disclaimer in text")
	}
}
