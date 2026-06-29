// Package renderはscene.SceneをPDFへ描画
// 描画プリミティブをCanvasに抽象化しsignintech/gopdf依存を隔離
package render

import (
	"math"

	"github.com/chan-mai/tsugu-mcp/internal/scene"
)

// 描画バックエンドの抽象(差し替え・テスト可能、座標mm・フォントpt)
type Canvas interface {
	Rect(x, y, w, h float64)
	Line(x1, y1, x2, y2, width float64, dashed bool, color scene.RGB)
	Text(x, y float64, text string, sizePt float64, align scene.Align)
}

// scene.SceneをCanvasへ描画(順序: 関係線→枠→単独テキスト)
func Draw(s scene.Scene, c Canvas) {
	for _, e := range s.Edges {
		drawEdge(c, e)
	}
	for _, b := range s.Boxes {
		if b.Border {
			c.Rect(b.X, b.Y, b.W, b.H)
		}
		for i, line := range b.Lines {
			c.Text(b.X+b.Pad, b.Y+b.Pad+float64(i)*b.LineHeight, line, b.FontSize, scene.AlignLeft)
		}
	}
	for _, l := range s.Labels {
		c.Text(l.X, l.Y, l.Text, l.Size, l.Align)
	}
}

const doubleLineOffset = 0.6 // 婚姻二重線の片側オフセット(mm)

func drawEdge(c Canvas, e scene.Edge) {
	for i := 0; i+1 < len(e.Points); i++ {
		a, b := e.Points[i], e.Points[i+1]
		if !e.Double {
			c.Line(a.X, a.Y, b.X, b.Y, e.Width, e.Dashed, e.Color)
			continue
		}
		ox, oy := perpOffset(a, b, doubleLineOffset)
		c.Line(a.X+ox, a.Y+oy, b.X+ox, b.Y+oy, e.Width, e.Dashed, e.Color)
		c.Line(a.X-ox, a.Y-oy, b.X-ox, b.Y-oy, e.Width, e.Dashed, e.Color)
	}
}

// 線分abに垂直方向の(dx,dy)を長さdで返す
func perpOffset(a, b scene.Pt, d float64) (float64, float64) {
	dx, dy := b.X-a.X, b.Y-a.Y
	length := math.Hypot(dx, dy)
	if length == 0 {
		return 0, d
	}
	return -dy / length * d, dx / length * d
}
