// Package layoutはfamily.Documentを横型・枠なしの相続関係説明図へ変換
// 純粋な幾何計算のみでPDF・描画ライブラリ非依存
package layout

import (
	"github.com/chan-mai/tsugu-mcp/family"
	"github.com/chan-mai/tsugu-mcp/internal/scene"
	"github.com/chan-mai/tsugu-mcp/internal/wareki"
)

type layouter struct {
	st      Style
	m       Measurer
	columnX []float64 // 世代ごとの左端X(列位置)
}

// 相続関係説明図1ページ分の描画指示を構築
func Build(doc family.Document, st Style, m Measurer) scene.Scene {
	l := &layouter{st: st, m: m}
	root := toTree(doc)
	l.computeColumns(root)

	graph, _, _ := l.hlayout(root, 0)

	out := scene.Scene{Width: st.PageW, Height: st.PageH}
	l.addTitle(&out, doc)

	// 作成欄分を下部に確保し残りを作図領域に
	createLines := l.createLines(doc)
	createH := blockHeight(l.st, len(createLines))
	reserve := 0.0
	if createH > 0 {
		reserve = createH + st.FooterGap
	}
	areaTop := st.Margin + st.TitlePt*ptToMM + st.TitleGap
	area := scene.Rect{
		X: st.Margin,
		Y: areaTop,
		W: st.PageW - 2*st.Margin,
		H: (st.PageH - st.Margin - reserve) - areaTop,
	}
	l.fitTopLeft(&graph, area)
	out.Append(graph)

	// 以下余白・作成欄は家系図直下に追従(下端固定にしない)
	bottom := graph.BBox().Y + graph.BBox().H
	l.addZanyo(&out, area.X+area.W, bottom+st.FooterGap)
	if createH > 0 {
		createY := bottom + st.FooterGap + st.BodyPt*ptToMM + 6
		if maxY := st.PageH - st.Margin - createH; createY > maxY {
			createY = maxY
		}
		l.addCreateBox(&out, createLines, createY)
	}
	return out
}

// ノードの部分木を横型配置し、局所座標(y上端=0)のシーン・接続点Y・高さを返却
func (l *layouter) hlayout(n *treeNode, depth int) (scene.Scene, float64, float64) {
	colX := l.columnX[depth]
	pc := l.buildCard(n.primary, false)
	var sc *card
	if n.spouse != nil {
		c := l.buildCard(*n.spouse, true)
		sc = &c
	}

	ownHeight := pc.height
	if sc != nil {
		ownHeight = pc.height + l.st.CoupleVGap + sc.height
	}

	// 自身(+配偶者)を上端topへ配置し接続点Yを返却
	emitOwn := func(s *scene.Scene, top float64) (primaryY, spouseY float64) {
		s.Boxes = append(s.Boxes, l.cardBox(pc, colX, top))
		primaryY = top + pc.anchorY
		if sc != nil {
			sTop := top + pc.height + l.st.CoupleVGap
			s.Boxes = append(s.Boxes, l.cardBox(*sc, colX, sTop))
			spouseY = sTop + sc.anchorY
			mx := colX + l.st.MarriageIndent
			// 二重線端と上下文字のクリアランス確保
			clear := l.st.BodyPt*ptToMM*0.5 + l.st.MarriageClear
			y1, y2 := primaryY+clear, spouseY-clear
			if y2 <= y1 {
				y1, y2 = primaryY, spouseY
			}
			s.Edges = append(s.Edges, scene.Edge{
				Points: []scene.Pt{{X: mx, Y: y1}, {X: mx, Y: y2}},
				Width:  l.st.LineWidth, Double: true,
			})
		}
		return
	}

	if len(n.children) == 0 {
		var s scene.Scene
		primaryY, _ := emitOwn(&s, 0)
		return s, primaryY, ownHeight
	}

	// 子部分木を縦積み
	var childScenes []scene.Scene
	var childAnchors []float64
	y := 0.0
	for i, c := range n.children {
		cs, ca, ch := l.hlayout(c, depth+1)
		cs.Translate(0, y)
		childScenes = append(childScenes, cs)
		childAnchors = append(childAnchors, y+ca)
		y += ch
		if i != len(n.children)-1 {
			y += l.st.SiblingVGap
		}
	}
	childrenBottom := y
	mid := (childAnchors[0] + childAnchors[len(childAnchors)-1]) / 2

	// 続柄行を子の縦中央へ整列、上はみ出しは全体を下げて補正
	ownTop := mid - pc.anchorY
	if ownTop < 0 {
		shift := -ownTop
		ownTop = 0
		for i := range childScenes {
			childScenes[i].Translate(0, shift)
		}
		for i := range childAnchors {
			childAnchors[i] += shift
		}
		childrenBottom += shift
	}

	var s scene.Scene
	primaryY, spouseY := emitOwn(&s, ownTop)

	connectStartX := colX + pc.width
	connectY := primaryY
	if sc != nil {
		connectStartX = colX + l.st.MarriageIndent
		connectY = (primaryY + spouseY) / 2
	}
	bracketX := l.columnX[depth+1] - l.st.BracketGap

	s.Edges = append(s.Edges, hLine(connectStartX, bracketX, connectY, l.st.LineWidth))
	top, bot := childAnchors[0], childAnchors[len(childAnchors)-1]
	top = min(top, connectY)
	bot = max(bot, connectY)
	if bot > top {
		s.Edges = append(s.Edges, vLine(bracketX, top, bot, l.st.LineWidth))
	}
	for _, ca := range childAnchors {
		s.Edges = append(s.Edges, hLine(bracketX, l.columnX[depth+1], ca, l.st.LineWidth))
	}
	for _, cs := range childScenes {
		s.Append(cs)
	}

	return s, primaryY, max(ownHeight, childrenBottom)
}

// 各世代の最大カード幅から列の左端Xを決定
func (l *layouter) computeColumns(root *treeNode) {
	maxW := map[int]float64{}
	var walk func(n *treeNode, d int)
	walk = func(n *treeNode, d int) {
		w := l.buildCard(n.primary, false).width
		if n.spouse != nil {
			if sw := l.buildCard(*n.spouse, true).width; sw > w {
				w = sw
			}
		}
		if w > maxW[d] {
			maxW[d] = w
		}
		for _, c := range n.children {
			walk(c, d+1)
		}
	}
	walk(root, 0)

	depth := 0
	for {
		if _, ok := maxW[depth]; !ok {
			break
		}
		depth++
	}
	l.columnX = make([]float64, depth+1)
	x := 0.0
	for d := 0; d <= depth; d++ {
		l.columnX[d] = x
		if d < depth {
			x += maxW[d] + l.st.GenGap
		}
	}
}

func (l *layouter) fitTopLeft(g *scene.Scene, area scene.Rect) {
	bb := g.BBox()
	if bb.W <= 0 || bb.H <= 0 {
		g.Translate(area.X-bb.X, area.Y-bb.Y)
		return
	}
	scale := 1.0
	if sx := area.W / bb.W; sx < scale {
		scale = sx
	}
	if sy := area.H / bb.H; sy < scale {
		scale = sy
	}
	if scale < l.st.MinScale {
		scale = l.st.MinScale
	}
	g.Scale(scale)
	bb = g.BBox()
	g.Translate(area.X-bb.X, area.Y-bb.Y)
}

func (l *layouter) addTitle(out *scene.Scene, doc family.Document) {
	out.Labels = append(out.Labels, scene.Label{
		X: l.st.Margin, Y: l.st.Margin,
		Text: "被相続人 " + doc.Decedent.Name + " 相続関係説明図",
		Size: l.st.TitlePt, Align: scene.AlignLeft,
	})
}

// 作成欄の枠内余白
const createPad = 2.5

// 作成日・作成者の行を生成(全角スペースで日付・住所・氏名の値を桁揃え)
func (l *layouter) createLines(doc family.Document) []string {
	var lines []string
	if !doc.PreparedAt.IsZero() {
		lines = append(lines, "作成日：　　　"+l.fmtDate(doc.PreparedAt))
	}
	switch {
	case doc.Preparer.Address != "" && doc.Preparer.Name != "":
		lines = append(lines, "作成者：住所　"+doc.Preparer.Address, "　　　　氏名　"+doc.Preparer.Name)
	case doc.Preparer.Name != "":
		lines = append(lines, "作成者：　　　"+doc.Preparer.Name)
	case doc.Preparer.Address != "":
		lines = append(lines, "作成者：住所　"+doc.Preparer.Address)
	}
	return lines
}

func blockHeight(st Style, n int) float64 {
	if n == 0 {
		return 0
	}
	return float64(n)*st.FootPt*ptToMM*st.LineSpacing + 2*createPad
}

// 作成欄を左に枠付き配置(上端Y=top)
func (l *layouter) addCreateBox(out *scene.Scene, lines []string, top float64) {
	st := l.st
	lh := st.FootPt * ptToMM * st.LineSpacing
	var maxW float64
	for _, s := range lines {
		if w := l.m.Measure(s, st.FootPt); w > maxW {
			maxW = w
		}
	}
	out.Boxes = append(out.Boxes, scene.Box{
		X: st.Margin, Y: top, W: maxW + 2*createPad, H: float64(len(lines))*lh + 2*createPad,
		Lines: lines, FontSize: st.FootPt, LineHeight: lh, Pad: createPad, Border: true,
	})
}

// 「以下余白」を右端rightXに右揃え配置
func (l *layouter) addZanyo(out *scene.Scene, rightX, y float64) {
	out.Labels = append(out.Labels, scene.Label{
		X: rightX, Y: y, Text: "以下余白", Size: l.st.BodyPt, Align: scene.AlignRight,
	})
}

func vLine(x, y1, y2, w float64) scene.Edge {
	return scene.Edge{Points: []scene.Pt{{X: x, Y: y1}, {X: x, Y: y2}}, Width: w}
}

func hLine(x1, x2, y, w float64) scene.Edge {
	return scene.Edge{Points: []scene.Pt{{X: x1, Y: y}, {X: x2, Y: y}}, Width: w}
}

func (l *layouter) fmtDate(d family.Date) string {
	return wareki.Format(d.Year, d.Month, d.Day, l.st.EraStyle)
}
