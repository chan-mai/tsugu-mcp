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
		Description: "相続関係説明図(A4縦1ページ)のPDFを生成しファイルパスを返す。被相続人を中心に配偶者・子孫(代襲)・直系尊属・兄弟姉妹を作図する。法定相続人や相続分の算定は行わず、与えた内容を忠実に描画する。" + disclaimer,
	}, handleChart)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_registration_application",
		Description: "相続登記申請書(所有権移転、不動産が多ければ複数ページ)のPDFを生成しファイルパスを返す。原因・相続人・申請人(複数+持分)・不動産の表示(土地/建物)を流し込む。各不動産に評価額(value)を渡すと登録免許税を自動計算して課税価格・登録免許税欄を補完する。相続分の算定は行わない。" + disclaimer,
	}, handleTouki)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_division_agreement",
		Description: "遺産分割協議書(財産が多ければ複数ページ)のPDFを生成しファイルパスを返す。被相続人・共同相続人・取得の対応(取得者+持分)・不動産/財産を流し込み、署名押印欄を付す。" + disclaimer,
	}, handleBunkatsu)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_division_certificate",
		Description: "遺産分割協議証明書(相続人ごとに1ページの個別書面型)のPDFを生成しファイルパスを返す。被相続人・取得者・共同相続人を渡すと、各相続人が取得を証明する書面を全員分の連続ページで出力する(住所欄は手書き用に空欄)。" + disclaimer,
	}, handleCertificate)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "calculate_registration_tax",
		Description: "相続登記の登録免許税を計算する(課税標準の合算・端数処理・免税措置の文言)。免税は自動適用せず、対象になり得る土地は注意のみ示す。" + disclaimer,
	}, handleTax)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_required_documents",
		Description: "相続登記の必要書類を相続方法・相続人パターンから案内する(添付情報4分類・入手先・手数料目安)。" + disclaimer,
	}, handleDocs)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "calculate_statutory_shares",
		Description: "法定相続分を算定する(現行・民法900/901条。配偶者と順位の組合せ・半血・代襲を既約分数で計算)。これは法律上の原則的割合であり、実際の配分判断には踏み込まない。" + disclaimer,
	}, handleShares)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "select_application_pattern",
		Description: "相続登記申請書のケース別様式(A法定相続〜H配偶者居住権)を選び、登記の目的・原因・申請構造・税率・添付の要点を返す。" + disclaimer,
	}, handlePattern)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "guide_heir_notification",
		Description: "相続登記の申請義務期限を算定し、期限内に難しい場合の相続人申告登記(簡易な義務履行)を案内する。" + disclaimer,
	}, handleNotify)

	addKnowledgeResources(s)
	addPrompts(s)
	return s
}

// MCPサーバーをstdioで起動(クライアント切断までブロック)
func Run(ctx context.Context) error {
	return newServer().Run(ctx, &mcp.StdioTransport{})
}
