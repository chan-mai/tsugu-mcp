package mcpserver

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/chan-mai/tsugu-mcp/docs/knowledge"
)

// 知識系ツールが調査基準日を構造化出力に持ち、テキストに基準日とファクトチェック文言を含むこと
func TestHandleShares_AsOfAndFactCheck(t *testing.T) {
	res, out, err := handleShares(context.Background(), nil, sharesToolInput{
		SpouseName: "花子",
		Children:   []shareHeir{{Name: "一郎", Alive: true}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.AsOf != knowledge.AsOf {
		t.Errorf("AsOf=%q want %q", out.AsOf, knowledge.AsOf)
	}
	if len(out.Result.Shares) == 0 {
		t.Errorf("result should be nested under .Result: %+v", out.Result)
	}

	text := res.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(text, knowledge.AsOf) {
		t.Errorf("text should contain as-of date %q", knowledge.AsOf)
	}
	if !strings.Contains(text, "一次情報") {
		t.Errorf("text should contain fact-check guidance")
	}
}
