package mcpserver

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"tsugu-mcp/internal/docguide"
	"tsugu-mcp/internal/regtax"
)

const disclaimer = "【免責】本ツールは本人申請の準備を支援する情報提供であり法的助言ではありません。税額・必要書類は参考情報で、個別事案の正確性・最新性は保証しません。"

// --- 登録免許税計算 ---

type taxProperty struct {
	Kind                 string  `json:"kind" jsonschema:"land(土地) または building(建物)"`
	Value                int     `json:"value,omitempty" jsonschema:"固定資産評価額(円)。私道はprivateRoadで認定"`
	ShareNum             int     `json:"shareNum,omitempty" jsonschema:"持分の分子(省略時は全部)"`
	ShareDen             int     `json:"shareDen,omitempty" jsonschema:"持分の分母(省略時は全部)"`
	Exemption            string  `json:"exemption,omitempty" jsonschema:"免税 none / small_value(100万円以下) / intermediate(数次中間者)。省略は適用なし"`
	PrivateRoadUnitPrice int     `json:"privateRoadUnitPrice,omitempty" jsonschema:"私道の近傍宅地1㎡単価(円)"`
	PrivateRoadArea      float64 `json:"privateRoadArea,omitempty" jsonschema:"私道の地積(平方メートル)"`
}

type taxToolInput struct {
	Properties []taxProperty `json:"properties" jsonschema:"対象不動産(1件以上)"`
}

func handleTax(_ context.Context, _ *mcp.CallToolRequest, in taxToolInput) (*mcp.CallToolResult, regtax.Result, error) {
	props := make([]regtax.Property, 0, len(in.Properties))
	for _, p := range in.Properties {
		rp := regtax.Property{Kind: p.Kind, Value: p.Value, ShareNum: p.ShareNum, ShareDen: p.ShareDen, Exemption: p.Exemption}
		if p.PrivateRoadUnitPrice > 0 {
			rp.PrivateRoad = &regtax.PrivateRoad{NeighborUnitPrice: p.PrivateRoadUnitPrice, Area: p.PrivateRoadArea}
		}
		props = append(props, rp)
	}
	r := regtax.Calculate(regtax.Input{Properties: props})
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: formatTax(r) + "\n\n" + disclaimer}}}, r, nil
}

func formatTax(r regtax.Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, "課税標準: %s円\n登録免許税: %s円\n", comma(r.TaxableTotal), comma(r.Tax))
	for _, l := range r.Lines {
		status := "課税"
		if !l.Taxable {
			status = "免税(" + l.Exemption + ")"
		}
		fmt.Fprintf(&b, "  物件%d: 課税価格 %s円 [%s]\n", l.Index+1, comma(l.TaxValue), status)
	}
	for _, s := range r.ExemptStatements {
		fmt.Fprintf(&b, "申請書記載: %s\n", s)
	}
	for _, n := range r.EligibilityNotes {
		fmt.Fprintf(&b, "注意: %s\n", n)
	}
	b.WriteString(r.Note)
	return b.String()
}

// --- 必要書類ナビ ---

type docToolInput struct {
	Method                            string `json:"method" jsonschema:"相続方法 legal(法定相続) / agreement(遺産分割協議) / will(遺言)"`
	HeirPattern                       string `json:"heirPattern" jsonschema:"相続人パターン children / substitution(代襲) / ascendants(第2順位) / siblings(第3順位) / siblings_substitution"`
	RegistryAddressDiffersFromHonseki bool   `json:"registryAddressDiffersFromHonseki,omitempty" jsonschema:"登記上の住所が本籍と異なる場合true(同一性証明)"`
	UseLegalInfoNumber                bool   `json:"useLegalInfoNumber,omitempty" jsonschema:"法定相続情報番号/一覧図で戸籍束を省略する場合true"`
	ApplicantAtWindow                 bool   `json:"applicantAtWindow,omitempty" jsonschema:"本人が窓口で広域交付を使う場合true(注意喚起)"`
}

func handleDocs(_ context.Context, _ *mcp.CallToolRequest, in docToolInput) (*mcp.CallToolResult, docguide.Result, error) {
	r := docguide.RequiredDocuments(docguide.Input{
		Method:                            in.Method,
		HeirPattern:                       in.HeirPattern,
		RegistryAddressDiffersFromHonseki: in.RegistryAddressDiffersFromHonseki,
		UseLegalInfoNumber:                in.UseLegalInfoNumber,
		ApplicantAtWindow:                 in.ApplicantAtWindow,
	})
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: formatDocs(r) + "\n\n" + disclaimer}}}, r, nil
}

func formatDocs(r docguide.Result) string {
	var b strings.Builder
	for _, c := range r.Categories {
		fmt.Fprintf(&b, "【%s】\n", c.Name)
		for _, it := range c.Items {
			fmt.Fprintf(&b, "  - %s\n", it)
		}
	}
	if len(r.Notes) > 0 {
		b.WriteString("【注意】\n")
		for _, n := range r.Notes {
			fmt.Fprintf(&b, "  - %s\n", n)
		}
	}
	if len(r.Fees) > 0 {
		b.WriteString("【手数料目安】\n")
		for _, f := range r.Fees {
			fmt.Fprintf(&b, "  - %s\n", f)
		}
	}
	return b.String()
}

// 評価額が指定された不動産があれば登録免許税を計算し、申請書の課税価格・登録免許税を補完
// 既入力値は上書きせず、補足説明(免税文言・注意・免責)を返す
func autoFillTax(doc *toukiDoc) string {
	hasValue := false
	props := make([]regtax.Property, 0, len(doc.Properties))
	for _, p := range doc.Properties {
		if p.Value > 0 {
			hasValue = true
		}
		props = append(props, regtax.Property{Kind: normalizeKind(p.Kind), Value: p.Value, Exemption: p.Exemption})
	}
	if !hasValue {
		return ""
	}
	r := regtax.Calculate(regtax.Input{Properties: props})
	if doc.TaxValue == "" {
		doc.TaxValue = comma(r.TaxableTotal)
	}
	if doc.RegistrationTax == "" {
		if r.Tax == 0 && len(r.ExemptStatements) > 0 {
			doc.RegistrationTax = strings.Join(r.ExemptStatements, "　")
		} else {
			doc.RegistrationTax = comma(r.Tax)
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "登録免許税を自動計算: 課税標準 %s円 / 税額 %s円", comma(r.TaxableTotal), comma(r.Tax))
	for _, s := range r.ExemptStatements {
		fmt.Fprintf(&b, "\n免税: %s(申請書に条文記載が必要)", s)
	}
	for _, n := range r.EligibilityNotes {
		fmt.Fprintf(&b, "\n注意: %s", n)
	}
	b.WriteString("\n" + disclaimer)
	return b.String()
}

// 不動産の種別をregtaxのland|buildingへ正規化
func normalizeKind(k string) string {
	switch k {
	case "building", "建物", "condominium", "区分建物":
		return "building"
	default:
		return "land"
	}
}

// 整数を3桁区切りにする
func comma(n int) string {
	s := strconv.Itoa(n)
	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	out := strings.Join(parts, ",")
	if neg {
		return "-" + out
	}
	return out
}
