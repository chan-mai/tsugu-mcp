package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// 相続登記準備のガイド導線をMCPプロンプトとして登録
func addPrompts(s *mcp.Server) {
	s.AddPrompt(&mcp.Prompt{
		Name:        "prepare_inheritance_registration",
		Title:       "相続登記の準備ガイド",
		Description: "ヒアリングからケース判定・必要書類・税額・書類生成まで一気通貫で進める導線",
	}, guidePrompt)
}

func guidePrompt(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "相続登記の準備ガイド",
		Messages: []*mcp.PromptMessage{{
			Role:    "user",
			Content: &mcp.TextContent{Text: guideText},
		}},
	}, nil
}

const guideText = `あなたは日本の相続登記(不動産の名義変更)の準備を支援するアシスタントです。次の手順で進めてください。

1. 前提知識: 必要に応じて knowledge:// リソース(index・01〜07)を参照する。
2. ヒアリング: 被相続人(氏名・死亡日・最後の住所/本籍・登記上の住所)、相続人全員、対象不動産(登記事項証明書の表示)、相続方法(法定相続/遺産分割協議/遺言)を確認する。
3. 相続人パターン判定: children / substitution(代襲) / ascendants(第2順位) / siblings(第3順位) / siblings_substitution のいずれか。
4. 必要書類: list_required_documents で添付書類・戸籍範囲・入手先・手数料を提示する。
5. 登録免許税: 各不動産の固定資産評価額を確認し calculate_registration_tax で税額を算出する。
6. 書類生成:
   - generate_relationship_chart(相続関係説明図)
   - 遺産分割協議なら generate_division_agreement(遺産分割協議書)
   - generate_registration_application(相続登記申請書。各不動産に value=評価額 を渡すと登録免許税を自動計算)
7. 境界: 法定相続人や相続分の最終判断・配分の助言には踏み込まない。免税は申請書への条文記載が必要(自動適用されない)。

【免責】本ガイドは本人申請の準備を支援する情報提供であり法的助言ではありません。最終確認は司法書士・弁護士へ。`
