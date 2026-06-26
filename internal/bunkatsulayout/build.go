// Package bunkatsulayoutはbunkatsu.Agreementを遺産分割協議書の描画指示へ変換
// 本文は折返し付きの流し込みで、長文や財産多数は複数ページへ
package bunkatsulayout

import (
	"fmt"
	"strings"

	"tsugu-mcp/bunkatsu"
	"tsugu-mcp/internal/scene"
	"tsugu-mcp/internal/wareki"
	"tsugu-mcp/touki"
	"tsugu-mcp/ymd"
)

const ptToMM = 25.4 / 72.0

// レイアウトの寸法・書式パラメータ
type Style struct {
	PageW, PageH float64
	Margin       float64
	TitlePt      float64
	BodyPt       float64
	LineH        float64
	ParaGap      float64
	Indent       float64
	LineWidth    float64
	Era          wareki.Style
}

// A4縦・遺産分割協議書向けの既定スタイル
func DefaultStyle() Style {
	return Style{
		PageW: 210, PageH: 297, Margin: 22,
		TitlePt: 15, BodyPt: 10.5, LineH: 7,
		ParaGap: 4, Indent: 4,
		LineWidth: 0.3,
		Era:       wareki.Wareki,
	}
}

func (s Style) em() float64       { return s.BodyPt * ptToMM }
func (s Style) contentW() float64 { return s.PageW - 2*s.Margin }
func (s Style) bottomY() float64  { return s.PageH - s.Margin }

// 遺産分割協議書の描画指示をページ配列で構築
func Build(a bunkatsu.Agreement, st Style) []scene.Scene {
	b := &builder{st: st}
	b.newPage()
	b.y += 4
	b.put(st.PageW/2, b.y, st.TitlePt, spaceOut("遺産分割協議書"), scene.AlignCenter)
	b.y += st.TitlePt*ptToMM + 8
	b.body(a)
	return b.finish()
}

type builder struct {
	st    Style
	pages []scene.Scene
	cur   scene.Scene
	y     float64
	page  int
}

func (b *builder) newPage() {
	if b.page > 0 {
		b.pages = append(b.pages, b.cur)
	}
	b.cur = scene.Scene{Width: b.st.PageW, Height: b.st.PageH}
	b.page++
	b.y = b.st.Margin
}

func (b *builder) finish() []scene.Scene {
	b.pages = append(b.pages, b.cur)
	return b.pages
}

func (b *builder) flow(h float64) {
	if b.y+h > b.st.bottomY() {
		b.newPage()
	}
}

func (b *builder) put(x, y, size float64, s string, align scene.Align) {
	if s == "" {
		return
	}
	b.cur.Labels = append(b.cur.Labels, scene.Label{X: x, Y: y, Text: s, Size: size, Align: align})
}

// 折返し付き段落
func (b *builder) para(text string, indent float64) {
	maxW := b.st.contentW() - indent
	for _, ln := range wrapText(text, maxW, b.st.em()) {
		b.flow(b.st.LineH)
		b.put(b.st.Margin+indent, b.y, b.st.BodyPt, ln, scene.AlignLeft)
		b.y += b.st.LineH
	}
	b.y += b.st.ParaGap
}

func (b *builder) body(a bunkatsu.Agreement) {
	st := b.st
	d := a.Decedent
	opening := fmt.Sprintf("%s、%s　%s の死亡によって開始した相続の共同相続人である%sは、本日、その相続財産について、次のとおり遺産分割の協議を行った。",
		b.wareki(d.DeathDate), d.Address, d.Name, joinJp(heirNames(a.Heirs)))
	b.para(opening, 0)

	for _, al := range a.Allocations {
		kind := "財産"
		if len(al.Properties) > 0 {
			kind = "不動産"
		}
		b.para(fmt.Sprintf("相続財産のうち、下記の%sは、%sが相続する。", kind, joinJp(acquirerStrs(al.Acquirers))), 0)
		for _, p := range al.Properties {
			b.propertyBlock(p)
		}
		for _, it := range al.Items {
			b.para(it, st.Indent)
		}
		b.y += st.ParaGap
	}

	b.para(fmt.Sprintf("この協議を証するため、本協議書を%d通作成して、それぞれに署名、押印し、各自1通を保有するものとする。", a.CopyCount()), 0)

	if !a.AgreedDate.IsZero() {
		b.y += st.ParaGap
		b.flow(st.LineH)
		b.put(st.Margin, b.y, st.BodyPt, b.wareki(a.AgreedDate), scene.AlignLeft)
		b.y += st.LineH + st.ParaGap
	}

	for _, h := range a.Heirs {
		b.signature(h)
	}
}

func (b *builder) propertyBlock(p touki.Property) {
	st := b.st
	rows := propertyRows(p)
	b.flow(float64(len(rows))*st.LineH + 2)
	x := st.Margin + st.Indent
	valX := x + 5*st.em() + 4
	for _, r := range rows {
		b.put(x, b.y, st.BodyPt, justify(r[0], 5), scene.AlignLeft)
		b.put(valX, b.y, st.BodyPt, r[1], scene.AlignLeft)
		b.y += st.LineH
	}
	b.y += 2
}

func (b *builder) signature(h bunkatsu.Heir) {
	st := b.st
	b.flow(3 * st.LineH)
	x := st.Margin + st.Indent
	if h.Address != "" {
		b.put(x, b.y, st.BodyPt, h.Address, scene.AlignLeft)
		b.y += st.LineH
	}
	b.put(x, b.y, st.BodyPt, "氏名　"+h.Name+"　　　　　　　　㊞", scene.AlignLeft)
	b.y += st.LineH + st.ParaGap
}

func (b *builder) wareki(d ymd.Date) string {
	return wareki.Format(d.Year, d.Month, d.Day, b.st.Era)
}

// --- ヘルパ ---

func propertyRows(p touki.Property) [][2]string {
	rows := [][2]string{{"不動産番号", p.Number}, {"所在", p.Location}}
	switch p.Kind {
	case touki.Building:
		rows = append(rows,
			[2]string{"家屋番号", p.HouseNumber},
			[2]string{"種類", p.BuildingType},
			[2]string{"構造", p.Structure},
			[2]string{"床面積", p.FloorArea})
	default:
		rows = append(rows,
			[2]string{"地番", p.LotNumber},
			[2]string{"地目", p.LandCategory},
			[2]string{"地積", area(p.Area)})
	}
	return rows
}

func area(s string) string {
	if s == "" {
		return ""
	}
	return s + "平方メートル"
}

func heirNames(hs []bunkatsu.Heir) []string {
	out := make([]string, len(hs))
	for i, h := range hs {
		out[i] = h.Name
	}
	return out
}

func acquirerStrs(as []bunkatsu.Acquirer) []string {
	out := make([]string, len(as))
	for i, a := range as {
		out[i] = a.Name
		if a.Share != "" {
			out[i] = fmt.Sprintf("%s(持分%s)", a.Name, a.Share)
		}
	}
	return out
}

// 氏名を「A、B及びC」の形で連結
func joinJp(names []string) string {
	switch len(names) {
	case 0:
		return ""
	case 1:
		return names[0]
	default:
		return strings.Join(names[:len(names)-1], "、") + "及び" + names[len(names)-1]
	}
}

// 文字間に全角スペースを挿入(表題用)
func spaceOut(s string) string {
	return strings.Join(strings.Split(s, ""), "　")
}

// 行頭禁則文字(これらでは改行しない)
const gyoutoukinsoku = "、。）」』】，．）)"

// テキストをmaxWidth(mm)で折り返す(全角em・半角0.5em、簡易禁則あり)
func wrapText(s string, maxWidth, em float64) []string {
	var lines []string
	var cur []rune
	var w float64
	for _, r := range s {
		rw := em
		if r < 0x0100 {
			rw = em * 0.5
		}
		if w+rw > maxWidth && len(cur) > 0 && !strings.ContainsRune(gyoutoukinsoku, r) {
			lines = append(lines, string(cur))
			cur, w = nil, 0
		}
		cur = append(cur, r)
		w += rw
	}
	if len(cur) > 0 {
		lines = append(lines, string(cur))
	}
	return lines
}

// ラベルをtarget文字幅へ均等割付
func justify(s string, target int) string {
	runes := []rune(s)
	n := len(runes)
	if n == 0 || n >= target {
		return s
	}
	gaps := n - 1
	if gaps == 0 {
		return s
	}
	total := target - n
	per, extra := total/gaps, total%gaps
	var sb strings.Builder
	for i, r := range runes {
		sb.WriteRune(r)
		if i < gaps {
			c := per
			if i < extra {
				c++
			}
			for j := 0; j < c; j++ {
				sb.WriteRune('　')
			}
		}
	}
	return sb.String()
}
