package mcpserver

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"tsugu-mcp/agreement"
)

type certDecedent struct {
	Name      string `json:"name" jsonschema:"被相続人の氏名"`
	DeathDate string `json:"deathDate" jsonschema:"死亡日 YYYY-MM-DD"`
}

type certDoc struct {
	Decedent certDecedent `json:"decedent" jsonschema:"被相続人"`
	Acquirer string       `json:"acquirer" jsonschema:"不動産を取得する相続人の氏名"`
	Signers  []string     `json:"signers" jsonschema:"共同相続人全員の氏名(各自1ページを生成)"`
	SignDate string       `json:"signDate,omitempty" jsonschema:"署名日 YYYY-MM-DD(年のみ印字、月日は空欄。任意)"`
}

type certToolInput struct {
	Document   certDoc `json:"document" jsonschema:"遺産分割協議証明書の内容"`
	OutputPath string  `json:"outputPath,omitempty" jsonschema:"出力PDFのパス(省略時は一時ファイル)"`
	Era        string  `json:"era,omitempty" jsonschema:"日付表記 wareki|both|seireki(既定 wareki)"`
}

func handleCertificate(_ context.Context, _ *mcp.CallToolRequest, in certToolInput) (*mcp.CallToolResult, toolResult, error) {
	opt, err := parseEra(in.Era)
	if err != nil {
		return errorResult(err)
	}
	data, err := json.Marshal(in.Document)
	if err != nil {
		return errorResult(err)
	}
	pdf, err := agreement.GenerateCertificateFromJSON(data, agreement.Options{Era: opt})
	if err != nil {
		return errorResult(err)
	}
	path, err := writePDF(pdf, in.OutputPath, "certificate")
	if err != nil {
		return errorResult(err)
	}
	return okResult(path, len(pdf))
}
