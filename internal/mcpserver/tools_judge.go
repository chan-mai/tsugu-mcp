package mcpserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/chan-mai/tsugu-mcp/internal/casepattern"
	"github.com/chan-mai/tsugu-mcp/internal/notification"
	"github.com/chan-mai/tsugu-mcp/internal/shares"
	"github.com/chan-mai/tsugu-mcp/ymd"
)

// --- 法定相続分の算定 ---

// 再代襲(2階層目、子の直系卑属のみ)
type shareSub2 struct {
	Name         string `json:"name" jsonschema:"氏名"`
	Alive        bool   `json:"alive,omitempty" jsonschema:"生存ならtrue"`
	Renounced    bool   `json:"renounced,omitempty" jsonschema:"相続放棄ならtrue"`
	Disqualified bool   `json:"disqualified,omitempty" jsonschema:"欠格・廃除ならtrue"`
}

// 代襲(1階層目、孫・甥姪)
type shareSub struct {
	Name         string      `json:"name" jsonschema:"氏名"`
	Alive        bool        `json:"alive,omitempty" jsonschema:"生存ならtrue"`
	Renounced    bool        `json:"renounced,omitempty" jsonschema:"相続放棄ならtrue"`
	Disqualified bool        `json:"disqualified,omitempty" jsonschema:"欠格・廃除ならtrue"`
	Substitutes  []shareSub2 `json:"substitutes,omitempty" jsonschema:"再代襲者(子の系統のみ)"`
}

type shareHeir struct {
	Name         string     `json:"name" jsonschema:"氏名"`
	Alive        bool       `json:"alive,omitempty" jsonschema:"生存ならtrue"`
	Renounced    bool       `json:"renounced,omitempty" jsonschema:"相続放棄ならtrue(代襲しない)"`
	Disqualified bool       `json:"disqualified,omitempty" jsonschema:"欠格・廃除ならtrue(代襲する)"`
	HalfBlood    bool       `json:"halfBlood,omitempty" jsonschema:"半血の兄弟姉妹ならtrue(相続分は全血の半分)"`
	Substitutes  []shareSub `json:"substitutes,omitempty" jsonschema:"代襲者(子→孫、兄弟姉妹→甥姪)"`
}

type sharesToolInput struct {
	DeathDate  string      `json:"deathDate,omitempty" jsonschema:"被相続人の死亡日 YYYY-MM-DD(適用法の確認)"`
	SpouseName string      `json:"spouseName,omitempty" jsonschema:"配偶者の氏名(いなければ空)"`
	Children   []shareHeir `json:"children,omitempty" jsonschema:"子(第1順位)"`
	Ascendants []shareHeir `json:"ascendants,omitempty" jsonschema:"直系尊属(第2順位)"`
	Siblings   []shareHeir `json:"siblings,omitempty" jsonschema:"兄弟姉妹(第3順位)"`
}

func handleShares(_ context.Context, _ *mcp.CallToolRequest, in sharesToolInput) (*mcp.CallToolResult, shares.Result, error) {
	dd, err := parseDateOpt(in.DeathDate)
	if err != nil {
		return textErr(err), shares.Result{}, nil
	}
	r := shares.Calculate(shares.Input{
		DeathDate:  dd,
		SpouseName: in.SpouseName,
		Children:   toHeirs(in.Children),
		Ascendants: toHeirs(in.Ascendants),
		Siblings:   toHeirs(in.Siblings),
	})
	return textOK(formatShares(r)), r, nil
}

func toHeirs(hs []shareHeir) []shares.Heir {
	out := make([]shares.Heir, 0, len(hs))
	for _, h := range hs {
		e := shares.Heir{Name: h.Name, Alive: h.Alive, Renounced: h.Renounced, Disqualified: h.Disqualified, HalfBlood: h.HalfBlood}
		for _, s := range h.Substitutes {
			sub := shares.Heir{Name: s.Name, Alive: s.Alive, Renounced: s.Renounced, Disqualified: s.Disqualified}
			for _, s2 := range s.Substitutes {
				sub.Substitutes = append(sub.Substitutes, shares.Heir{Name: s2.Name, Alive: s2.Alive, Renounced: s2.Renounced, Disqualified: s2.Disqualified})
			}
			e.Substitutes = append(e.Substitutes, sub)
		}
		out = append(out, e)
	}
	return out
}

func formatShares(r shares.Result) string {
	var b strings.Builder
	b.WriteString("法定相続分:\n")
	for _, s := range r.Shares {
		fmt.Fprintf(&b, "  %s %s: %s\n", s.Role, s.Name, s.Fraction)
	}
	fmt.Fprintf(&b, "検算合計: %s\n", r.Sum)
	for _, w := range r.Warnings {
		fmt.Fprintf(&b, "warning: %s\n", w)
	}
	fmt.Fprintf(&b, "注: %s\n\n%s", r.Note, disclaimer)
	return b.String()
}

// --- ケース別様式の選択 ---

type patternToolInput struct {
	Method           string `json:"method" jsonschema:"相続方法 legal(法定相続) / agreement(遺産分割) / will_specified(相続させる遺言) / bequest(遺贈)"`
	BequestToHeir    bool   `json:"bequestToHeir,omitempty" jsonschema:"遺贈で受遺者が相続人ならtrue(第三者ならfalse)"`
	Multilevel       bool   `json:"multilevel,omitempty" jsonschema:"数次相続ならtrue"`
	Substitution     bool   `json:"substitution,omitempty" jsonschema:"代襲相続ならtrue"`
	Renunciation     bool   `json:"renunciation,omitempty" jsonschema:"相続放棄があればtrue"`
	SpousalResidence bool   `json:"spousalResidence,omitempty" jsonschema:"配偶者居住権の設定があればtrue"`
}

func handlePattern(_ context.Context, _ *mcp.CallToolRequest, in patternToolInput) (*mcp.CallToolResult, casepattern.Result, error) {
	r := casepattern.Select(casepattern.Input{
		Method: in.Method, BequestToHeir: in.BequestToHeir, Multilevel: in.Multilevel,
		Substitution: in.Substitution, Renunciation: in.Renunciation, SpousalResidence: in.SpousalResidence,
	})
	return textOK(formatPattern(r)), r, nil
}

func formatPattern(r casepattern.Result) string {
	p := r.Primary
	var b strings.Builder
	fmt.Fprintf(&b, "ケース %s: %s\n", p.Key, p.Name)
	fmt.Fprintf(&b, "  登記の目的: %s\n  原因: %s\n  申請構造: %s\n  税率: %s\n  登記原因証明情報: %s\n", p.Purpose, p.Cause, p.Structure, p.TaxRate, p.OriginInfo)
	if p.Caveat != "" {
		fmt.Fprintf(&b, "  確実性: %s\n", p.Caveat)
	}
	for _, m := range r.Modifiers {
		fmt.Fprintf(&b, "追加注意: %s\n", m)
	}
	fmt.Fprintf(&b, "注: %s\n\n%s", r.Note, disclaimer)
	return b.String()
}

// --- 相続人申告登記の案内 ---

type notifyToolInput struct {
	DeathDate string `json:"deathDate,omitempty" jsonschema:"相続開始日 YYYY-MM-DD"`
	KnownDate string `json:"knownDate,omitempty" jsonschema:"相続開始・取得を知った日 YYYY-MM-DD(空なら死亡日)"`
}

func handleNotify(_ context.Context, _ *mcp.CallToolRequest, in notifyToolInput) (*mcp.CallToolResult, notification.Result, error) {
	dd, err := parseDateOpt(in.DeathDate)
	if err != nil {
		return textErr(err), notification.Result{}, nil
	}
	kd, err := parseDateOpt(in.KnownDate)
	if err != nil {
		return textErr(err), notification.Result{}, nil
	}
	r := notification.Guide(notification.Input{DeathDate: dd, KnownDate: kd})
	return textOK(formatNotify(r)), r, nil
}

func formatNotify(r notification.Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, "申請義務の期限: %s\n  (%s)\n", r.Deadline, r.DeadlineNote)
	b.WriteString("相続人申告登記:\n")
	for _, p := range r.Provisional {
		fmt.Fprintf(&b, "  - %s\n", p)
	}
	fmt.Fprintf(&b, "過料: %s\n注: %s\n\n%s", r.Penalty, r.Note, disclaimer)
	return b.String()
}

// --- 共通ヘルパ ---

func parseDateOpt(s string) (ymd.Date, error) {
	if strings.TrimSpace(s) == "" {
		return ymd.Date{}, nil
	}
	return ymd.Parse(s)
}

func textOK(s string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: s}}}
}

func textErr(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}}}
}
