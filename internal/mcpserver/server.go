// Package mcpserverは相続書類PDF生成のMCPサーバー(stdio)
package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const version = "0.1.0"

// ツールを登録したMCPサーバーを構築
func newServer() *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "tsugu-mcp", Version: version}, nil)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_relationship_chart",
		Description: "相続関係説明図(A4縦1ページ)のPDFを生成しファイルパスを返す。被相続人を中心に配偶者・子孫(代襲)・直系尊属・兄弟姉妹を作図する。法定相続人や相続分の算定は行わず、与えた内容を忠実に描画する。",
	}, handleChart)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_registration_application",
		Description: "相続登記申請書(所有権移転、不動産が多ければ複数ページ)のPDFを生成しファイルパスを返す。原因・相続人・申請人(複数+持分)・不動産の表示(土地/建物)を流し込む。課税価格・登録免許税・相続分の算定は行わない。",
	}, handleTouki)

	return s
}

// MCPサーバーをstdioで起動(クライアント切断までブロック)
func Run(ctx context.Context) error {
	return newServer().Run(ctx, &mcp.StdioTransport{})
}
