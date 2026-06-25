package render

import (
	"bytes"
	"testing"

	"tsugu-mcp/internal/scene"
)

// 呼び出し回数を記録するフェイクCanvas
type recCanvas struct {
	rects, lines, texts int
}

func (r *recCanvas) Rect(x, y, w, h float64)                                { r.rects++ }
func (r *recCanvas) Line(x1, y1, x2, y2, width float64, dashed bool)        { r.lines++ }
func (r *recCanvas) Text(x, y float64, s string, sz float64, a scene.Align) { r.texts++ }

func TestDraw_EmitsPrimitives(t *testing.T) {
	s := scene.Scene{
		Boxes:  []scene.Box{{X: 1, Y: 1, W: 10, H: 10, Lines: []string{"a", "b"}, Border: true}},
		Edges:  []scene.Edge{{Points: []scene.Pt{{X: 0, Y: 0}, {X: 5, Y: 0}}, Double: true}},
		Labels: []scene.Label{{X: 0, Y: 0, Text: "t"}},
	}
	c := &recCanvas{}
	Draw(s, c)

	if c.rects != 1 {
		t.Errorf("rects = %d, want 1", c.rects)
	}
	if c.lines != 2 {
		t.Errorf("二重線は2本に展開されるべき: lines = %d, want 2", c.lines)
	}
	if c.texts != 3 {
		t.Errorf("texts = %d, want 3 (枠2行 + ラベル1)", c.texts)
	}
}

func TestToPDF_ProducesPDF(t *testing.T) {
	s := scene.Scene{
		Width: 210, Height: 297,
		Boxes: []scene.Box{{X: 20, Y: 20, W: 40, H: 20, Lines: []string{"山田 太郎"}, FontSize: 10, LineHeight: 5, Pad: 2, Border: true}},
	}
	b, err := ToPDF(s)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(b, []byte("%PDF")) {
		t.Errorf("PDFヘッダがありません: %q", b[:min(8, len(b))])
	}
	if len(b) < 1000 {
		t.Errorf("PDFが小さすぎます: %d bytes", len(b))
	}
}

func TestMeasurer_PositiveWidth(t *testing.T) {
	m, err := NewMeasurer()
	if err != nil {
		t.Fatal(err)
	}
	if w := m.Measure("山田太郎", 10); w <= 0 {
		t.Errorf("文字幅は正であるべき: %f", w)
	}
}
