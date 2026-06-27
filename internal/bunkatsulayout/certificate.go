package bunkatsulayout

import (
	"fmt"
	"strconv"

	"tsugu-mcp/bunkatsu"
	"tsugu-mcp/internal/scene"
	"tsugu-mcp/internal/wareki"
	"tsugu-mcp/ymd"
)

// 遺産分割協議証明書の描画指示を構築(共同相続人ごとに1ページ生成)
func BuildCertificate(c bunkatsu.Certificate, st Style) []scene.Scene {
	signers := c.Signers
	if len(signers) == 0 {
		signers = []string{""} // フォールバック: 氏名空欄で1枚
	}
	b := &builder{st: st}
	for _, name := range signers {
		b.newPage()
		b.y += 4
		b.put(st.PageW/2, b.y, st.TitlePt, spaceOut("証明書"), scene.AlignCenter)
		b.y += st.TitlePt*ptToMM + 10
		b.certificateBody(c, name)
	}
	return b.finish()
}

func (b *builder) certificateBody(c bunkatsu.Certificate, signer string) {
	st := b.st
	d := c.Decedent
	body := fmt.Sprintf("%s被相続人　%s　の死亡により、同日相続が開始したが、今般共同相続人全員で遺産分割協議の結果、被相続人　%s　名義の不動産は、共同相続人　%s　が、取得することに協議が成立したことを、共同相続人として証明します。",
		b.wareki(d.DeathDate), d.Name, d.Name, c.Acquirer)
	b.para(body, 0)

	b.y += st.ParaGap
	b.put(st.Margin, b.y, st.BodyPt, "被相続人　　"+d.Name, scene.AlignLeft)
	b.y += st.LineH + 2*st.ParaGap

	b.put(st.Margin, b.y, st.BodyPt, signDateLine(c.SignDate), scene.AlignLeft)
	b.y += st.LineH + 2*st.ParaGap

	b.put(st.Margin, b.y, st.BodyPt, "住　　所", scene.AlignLeft)
	b.y += st.LineH + st.ParaGap
	b.put(st.Margin, b.y, st.BodyPt, "上記相続人　氏　名　　"+signer+"　　　　　　㊞", scene.AlignLeft)
	b.y += st.LineH
}

// 署名日行(年のみ印字、月日は空欄)
func signDateLine(d ymd.Date) string {
	if d.IsZero() {
		return "令和　　　年　　　　月　　　　日"
	}
	era, y := wareki.EraYear(d.Year, d.Month, d.Day)
	if era == "" {
		era = "令和"
	}
	yStr := strconv.Itoa(y)
	if y == 1 {
		yStr = "元"
	}
	return fmt.Sprintf("%s　%s　年　　　　月　　　　日", era, yStr)
}
