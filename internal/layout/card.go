package layout

import "tsugu-mcp/internal/scene"

// 枠なし人物欄(anchor=続柄+氏名行、その中心Yで関係線接続)
type card struct {
	lines   []string
	width   float64
	height  float64
	anchorY float64 // カード上端から続柄行の中心までの距離(mm)
}

// 人物情報から表示行を生成(詳細は続柄行の上、死亡日は下)
func (l *layouter) buildCard(p person) card {
	var above []string
	if p.applicant {
		above = append(above, "(申立人)")
	}
	if p.isDecedent {
		if p.address != "" {
			above = append(above, "最後の住所", "　"+p.address)
		}
		if p.honseki != "" {
			above = append(above, "最後の本籍", "　"+p.honseki)
		}
		// 被相続人の出生・死亡は改行せずスペース併記
		if !p.birth.IsZero() {
			above = append(above, "出生　"+l.fmtDate(p.birth))
		}
		if p.death != nil {
			above = append(above, "死亡　"+l.fmtDate(*p.death))
		}
	} else {
		if p.address != "" {
			above = append(above, "住所", "　"+p.address)
		}
		if !p.birth.IsZero() {
			above = append(above, "出生", "　"+l.fmtDate(p.birth))
		}
	}

	anchor := p.name
	if p.relationship != "" {
		anchor = "(" + p.relationship + ")　" + p.name
	}
	if tag := p.outcome.Label(); tag != "" {
		anchor += " " + tag
	}

	var below []string
	if !p.isDecedent && p.death != nil {
		below = append(below, l.fmtDate(*p.death)+"　死亡")
	}

	return l.makeCard(above, anchor, below)
}

func (l *layouter) makeCard(above []string, anchor string, below []string) card {
	lines := make([]string, 0, len(above)+1+len(below))
	lines = append(lines, above...)
	lines = append(lines, anchor)
	lines = append(lines, below...)

	lh := l.st.lineHeight()
	var maxW float64
	for _, s := range lines {
		if w := l.m.Measure(s, l.st.BodyPt); w > maxW {
			maxW = w
		}
	}
	return card{
		lines:   lines,
		width:   maxW,
		height:  float64(len(lines)) * lh,
		anchorY: float64(len(above))*lh + l.st.BodyPt*ptToMM*0.5,
	}
}

// cardを枠なしscene.Box(テキストのみ)へ変換
func (l *layouter) cardBox(c card, x, top float64) scene.Box {
	return scene.Box{
		X: x, Y: top, W: c.width, H: c.height,
		Lines:      c.lines,
		FontSize:   l.st.BodyPt,
		LineHeight: l.st.lineHeight(),
		Pad:        0,
		Border:     false,
	}
}
