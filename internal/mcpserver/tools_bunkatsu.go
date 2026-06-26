package mcpserver

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"tsugu-mcp/agreement"
)

type bunkatsuDecedent struct {
	Name      string `json:"name" jsonschema:"被相続人の氏名"`
	Address   string `json:"address,omitempty" jsonschema:"最後の住所"`
	DeathDate string `json:"deathDate" jsonschema:"死亡日 YYYY-MM-DD"`
}

type bunkatsuHeir struct {
	Name    string `json:"name" jsonschema:"氏名"`
	Address string `json:"address,omitempty" jsonschema:"住所(署名欄)"`
}

type bunkatsuAcquirer struct {
	Name  string `json:"name" jsonschema:"取得者の氏名"`
	Share string `json:"share,omitempty" jsonschema:"持分(例 2分の1。単独なら空)"`
}

type bunkatsuAllocation struct {
	Acquirers  []bunkatsuAcquirer `json:"acquirers" jsonschema:"取得者(1名以上)"`
	Properties []toukiProperty    `json:"properties,omitempty" jsonschema:"取得する不動産"`
	Items      []string           `json:"items,omitempty" jsonschema:"不動産以外の財産(預貯金等のフリー記述)"`
}

type bunkatsuDoc struct {
	Decedent    bunkatsuDecedent     `json:"decedent" jsonschema:"被相続人"`
	Heirs       []bunkatsuHeir       `json:"heirs" jsonschema:"共同相続人全員"`
	Allocations []bunkatsuAllocation `json:"allocations" jsonschema:"取得の対応(1つ以上)"`
	AgreedDate  string               `json:"agreedDate,omitempty" jsonschema:"協議成立日 YYYY-MM-DD"`
	Copies      int                  `json:"copies,omitempty" jsonschema:"作成通数(省略時は相続人数)"`
}

type bunkatsuToolInput struct {
	Document   bunkatsuDoc `json:"document" jsonschema:"遺産分割協議書の内容"`
	OutputPath string      `json:"outputPath,omitempty" jsonschema:"出力PDFのパス(省略時は一時ファイル)"`
	Era        string      `json:"era,omitempty" jsonschema:"日付表記 wareki|both|seireki(既定 wareki)"`
}

func handleBunkatsu(_ context.Context, _ *mcp.CallToolRequest, in bunkatsuToolInput) (*mcp.CallToolResult, toolResult, error) {
	opt, err := parseEra(in.Era)
	if err != nil {
		return errorResult(err)
	}
	data, err := json.Marshal(in.Document)
	if err != nil {
		return errorResult(err)
	}
	pdf, err := agreement.GenerateFromJSON(data, agreement.Options{Era: opt})
	if err != nil {
		return errorResult(err)
	}
	path, err := writePDF(pdf, in.OutputPath, "bunkatsu")
	if err != nil {
		return errorResult(err)
	}
	return okResult(path, len(pdf))
}
