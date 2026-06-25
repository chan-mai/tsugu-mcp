// Package scene は描画指示の中間表現(PDF非依存の図形プリミティブ)
// 座標系mm、原点左上、Y下向き
package scene

import "math"

// テキストの水平揃え
type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// 座標点(mm)
type Pt struct{ X, Y float64 }

// 矩形領域(mm)
type Rect struct{ X, Y, W, H float64 }

// テキスト欄(Linesを上から左揃え、Border時のみ枠描画)
type Box struct {
	X, Y, W, H float64
	Lines      []string
	FontSize   float64 // pt
	LineHeight float64 // mm 行送り
	Pad        float64 // mm 内側余白
	Border     bool
}

// 関係線(Points折れ線、Double=婚姻二重線)
type Edge struct {
	Points []Pt
	Width  float64 // mm 線幅
	Double bool
	Dashed bool
}

// 単独テキスト(タイトル・以下余白など)
type Label struct {
	X, Y  float64
	Text  string
	Size  float64 // pt
	Align Align
}

// 1ページ分の描画指示
type Scene struct {
	Width, Height float64 // ページサイズ(mm)
	Boxes         []Box
	Edges         []Edge
	Labels        []Label
}

// otherの要素を取り込む
func (s *Scene) Append(other Scene) {
	s.Boxes = append(s.Boxes, other.Boxes...)
	s.Edges = append(s.Edges, other.Edges...)
	s.Labels = append(s.Labels, other.Labels...)
}

// 全要素を囲む最小矩形(要素無しならゼロ値)
func (s Scene) BBox() Rect {
	first := true
	var minX, minY, maxX, maxY float64
	upd := func(x, y float64) {
		if first {
			minX, minY, maxX, maxY = x, y, x, y
			first = false
			return
		}
		minX, minY = math.Min(minX, x), math.Min(minY, y)
		maxX, maxY = math.Max(maxX, x), math.Max(maxY, y)
	}
	for _, b := range s.Boxes {
		upd(b.X, b.Y)
		upd(b.X+b.W, b.Y+b.H)
	}
	for _, e := range s.Edges {
		for _, p := range e.Points {
			upd(p.X, p.Y)
		}
	}
	for _, l := range s.Labels {
		upd(l.X, l.Y)
	}
	if first {
		return Rect{}
	}
	return Rect{X: minX, Y: minY, W: maxX - minX, H: maxY - minY}
}

// 全要素を平行移動
func (s *Scene) Translate(dx, dy float64) {
	for i := range s.Boxes {
		s.Boxes[i].X += dx
		s.Boxes[i].Y += dy
	}
	for i := range s.Edges {
		for j := range s.Edges[i].Points {
			s.Edges[i].Points[j].X += dx
			s.Edges[i].Points[j].Y += dy
		}
	}
	for i := range s.Labels {
		s.Labels[i].X += dx
		s.Labels[i].Y += dy
	}
}

// 原点中心に全要素(座標・寸法・フォント)を一様拡縮
func (s *Scene) Scale(f float64) {
	for i := range s.Boxes {
		b := &s.Boxes[i]
		b.X, b.Y, b.W, b.H = b.X*f, b.Y*f, b.W*f, b.H*f
		b.FontSize, b.LineHeight, b.Pad = b.FontSize*f, b.LineHeight*f, b.Pad*f
	}
	for i := range s.Edges {
		e := &s.Edges[i]
		e.Width *= f
		for j := range e.Points {
			e.Points[j].X *= f
			e.Points[j].Y *= f
		}
	}
	for i := range s.Labels {
		s.Labels[i].X, s.Labels[i].Y, s.Labels[i].Size = s.Labels[i].X*f, s.Labels[i].Y*f, s.Labels[i].Size*f
	}
}
