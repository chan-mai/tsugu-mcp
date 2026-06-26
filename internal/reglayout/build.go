package reglayout

import (
	"fmt"
	"strings"

	"tsugu-mcp/internal/scene"
	"tsugu-mcp/internal/wareki"
	"tsugu-mcp/touki"
	"tsugu-mcp/ymd"
)

// 登記申請書の描画指示をページ配列で構築
func Build(app touki.Application, st Style) []scene.Scene {
	b := &builder{st: st}
	b.newPage()
	b.header(app)
	for _, p := range app.Properties {
		b.property(p)
	}
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
	b.y = b.st.MarginTop
}

func (b *builder) finish() []scene.Scene {
	b.pages = append(b.pages, b.cur)
	return b.pages
}

// hの高さが現ページに収まらなければ改ページ
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

func (b *builder) body(x, y float64, s string) { b.put(x, y, b.st.BodyPt, s, scene.AlignLeft) }

func (b *builder) line(x1, y1, x2, y2 float64, dashed bool) {
	b.cur.Edges = append(b.cur.Edges, scene.Edge{
		Points: []scene.Pt{{X: x1, Y: y1}, {X: x2, Y: y2}}, Width: b.st.LineWidth, Dashed: dashed,
	})
}

// --- ヘッダ(1枚目) ---

func (b *builder) header(app touki.Application) {
	st := b.st

	// 受付番号表用の破線枠
	x0 := st.MarginX + st.ReceiptInset
	x1 := st.PageW - st.MarginX - st.ReceiptInset
	b.dashedRect(x0, st.MarginTop, x1-x0, st.ReceiptH)
	b.y = st.MarginTop + st.ReceiptH + 6

	// タイトル
	b.put(st.PageW/2, b.y, st.TitlePt, "登　記　申　請　書", scene.AlignCenter)
	b.y += st.TitlePt*ptToMM + 6

	b.fieldRow("登記の目的", "所有権移転")
	b.causeRow(app.Causes)
	b.decedentRows(app.Decedent)
	b.applicantRows(app.Applicants)
	b.attachmentRows(app.Attachments, app.DeclineIDInfo)
	b.fieldRow0(b.wareki(app.ApplicationDate) + "申請　" + app.Registry)
	b.y += 1
	b.fieldRow("課税価格", "金　"+orZero(app.TaxValue)+"　円")
	b.fieldRow("登録免許税", "金　"+orZero(app.RegistrationTax)+"　円")

	b.y += 2
	b.body(st.MarginX, b.y, "不動産の表示")
	b.y += st.LineH
}

func (b *builder) dashedRect(x, y, w, h float64) {
	seg := func(x1, y1, x2, y2 float64) {
		b.cur.Edges = append(b.cur.Edges, scene.Edge{
			Points: []scene.Pt{{X: x1, Y: y1}, {X: x2, Y: y2}},
			Width:  b.st.LineWidth, Dashed: true, Color: b.st.ReceiptColor,
		})
	}
	seg(x, y, x+w, y)
	seg(x+w, y, x+w, y+h)
	seg(x+w, y+h, x, y+h)
	seg(x, y+h, x, y)
}

// ラベル+値の1行
func (b *builder) fieldRow(label, value string) {
	b.flow(b.st.LineH)
	b.body(b.st.MarginX, b.y, justify(label, b.st.LabelChars))
	b.body(b.st.valueX(), b.y, value)
	b.y += b.st.LineH
}

// ラベル無しの1行(左余白から)
func (b *builder) fieldRow0(value string) {
	b.flow(b.st.LineH)
	b.body(b.st.MarginX, b.y, value)
	b.y += b.st.LineH
}

func (b *builder) causeRow(causes []touki.Cause) {
	var parts []string
	for _, c := range causes {
		s := c.Text
		if !c.Date.IsZero() {
			s = b.wareki(c.Date) + " " + c.Text
		}
		parts = append(parts, s)
	}
	b.fieldRow("原因", strings.Join(parts, "　"))
}

func (b *builder) decedentRows(d touki.Decedent) {
	b.flow(2 * b.st.LineH)
	b.body(b.st.MarginX, b.y, justify("相続人", b.st.LabelChars))
	b.body(b.st.valueX(), b.y, "（被相続人　"+d.Name+"　）")
	b.y += 2*b.st.LineH + 1 // 名の下に住所1行分の空きを残す
}

func (b *builder) applicantRows(apps []touki.Applicant) {
	for i, ap := range apps {
		b.flow(2 * b.st.LineH)
		label := ""
		if i == 0 {
			label = "（申請人）"
		}
		name := ap.Name
		if ap.Share != "" {
			name = "持分" + ap.Share + "　" + ap.Name
		}
		b.body(b.st.MarginX, b.y, label)
		b.body(b.st.valueX(), b.y, name)
		b.y += b.st.LineH
		if ap.Address != "" {
			b.body(b.st.valueX(), b.y, ap.Address)
			b.y += b.st.LineH
		}
		if ap.Contact {
			b.applicantTable(ap)
		}
		b.y += 1
	}
}

// 申請人の3行枠表(氏名ふりがな・生年月日・メール)+ 連絡先電話
func (b *builder) applicantTable(ap touki.Applicant) {
	st := b.st
	x := st.valueX()
	w := st.PageW - st.MarginX - x
	rowH := st.LineH
	h := rowH * 3

	b.flow(h + st.LineH + 2)
	top := b.y
	b.cur.Boxes = append(b.cur.Boxes, scene.Box{X: x, Y: top, W: w, H: h, Border: true})
	b.line(x, top+rowH, x+w, top+rowH, false)
	b.line(x, top+2*rowH, x+w, top+2*rowH, false)
	b.line(x+st.TableLabelW, top, x+st.TableLabelW, top+h, false)

	rows := [3][2]string{
		{"氏名ふりがな", ap.NameKana},
		{"生年月日", seireki(ap.BirthDate)},
		{"メールアドレス", ap.Email},
	}
	const pad = 1.6
	for i, r := range rows {
		cy := top + float64(i)*rowH + pad
		b.body(x+pad, cy, r[0])
		b.body(x+st.TableLabelW+pad, cy, r[1])
	}
	b.y = top + h + 1.5

	if ap.Phone != "" {
		b.body(x, b.y, "連絡先の電話番号　"+ap.Phone)
		b.y += st.LineH
	}
}

func (b *builder) attachmentRows(atts []string, decline bool) {
	// 添付情報ラベルは均等割付せず、項目は改行+インデントしスペース区切りで列挙
	b.flow(b.st.LineH)
	b.body(b.st.MarginX, b.y, "添付情報")
	b.y += b.st.LineH
	if len(atts) > 0 {
		b.flow(b.st.LineH)
		b.body(b.st.MarginX+b.st.em()*2, b.y, strings.Join(atts, "　"))
		b.y += b.st.LineH
	}
	// 登記識別情報の通知を希望しない欄は常に描画(falseでも空欄□)
	b.flow(b.st.LineH)
	tx := b.checkbox(b.st.MarginX, b.y, decline)
	b.body(tx, b.y, "登記識別情報の通知を希望しません")
	b.y += b.st.LineH + 1
}

// 小さな四角のチェックボックスを描き、後続テキストの開始Xを返す
func (b *builder) checkbox(x, y float64, checked bool) float64 {
	const size = 3.2
	top := y + 0.3
	b.cur.Boxes = append(b.cur.Boxes, scene.Box{X: x, Y: top, W: size, H: size, Border: true})
	if checked {
		b.line(x+0.22*size, top+0.55*size, x+0.42*size, top+0.80*size, false)
		b.line(x+0.42*size, top+0.80*size, x+0.82*size, top+0.18*size, false)
	}
	return x + size + 1.6
}

// --- 不動産ブロック ---

func (b *builder) property(p touki.Property) {
	rows := propertyRows(p)
	blockH := float64(len(rows))*b.st.LineH + b.st.BlockGap
	b.flow(blockH)

	x := b.st.MarginX + b.st.PropIndent
	valX := x + b.st.labelW() + b.st.ValueGap
	for _, r := range rows {
		b.body(x, b.y, justify(r[0], b.st.LabelChars))
		b.body(valX, b.y, r[1])
		b.y += b.st.LineH
	}
	b.y += b.st.BlockGap
}

func propertyRows(p touki.Property) [][2]string {
	rows := [][2]string{
		{"不動産番号", p.Number},
		{"所在", p.Location},
	}
	switch p.Kind {
	case touki.Building:
		rows = append(rows,
			[2]string{"家屋番号", p.HouseNumber},
			[2]string{"種類", p.BuildingType},
			[2]string{"構造", p.Structure},
			[2]string{"床面積", p.FloorArea},
		)
	default: // Land
		rows = append(rows,
			[2]string{"地番", p.LotNumber},
			[2]string{"地目", p.LandCategory},
			[2]string{"地積", area(p.Area)},
		)
	}
	return rows
}

// --- 整形ヘルパ ---

func (b *builder) wareki(d ymd.Date) string {
	return wareki.Format(d.Year, d.Month, d.Day, b.st.Era)
}

// 生年月日は西暦・ゼロ埋め(添付準拠)
func seireki(d ymd.Date) string {
	if d.IsZero() {
		return ""
	}
	return fmt.Sprintf("%d年%02d月%02d日", d.Year, d.Month, d.Day)
}

func area(s string) string {
	if s == "" {
		return ""
	}
	return s + "平方メートル"
}

func orZero(s string) string {
	if s == "" {
		return "0"
	}
	return s
}

// ラベル文字を全角スペースでtarget文字幅へ均等割付
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
			cnt := per
			if i < extra {
				cnt++
			}
			for j := 0; j < cnt; j++ {
				sb.WriteRune('　')
			}
		}
	}
	return sb.String()
}
