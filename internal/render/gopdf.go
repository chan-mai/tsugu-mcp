package render

import (
	"fmt"

	"github.com/signintech/gopdf"

	"tsugu-mcp/assets"
	"tsugu-mcp/internal/scene"
)

const fontName = "ipaexg"

// scene.SceneをA4縦のPDFバイト列へ描画
func ToPDF(s scene.Scene) ([]byte, error) {
	pdf, err := newDoc()
	if err != nil {
		return nil, err
	}
	Draw(s, &gopdfCanvas{pdf: pdf})
	return pdf.GetBytesPdf(), nil
}

func newDoc() (*gopdf.GoPdf, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: gopdf.UnitMM, PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	if err := pdf.AddTTFFontData(fontName, assets.IPAexGothic); err != nil {
		return nil, fmt.Errorf("failed to load font: %w", err)
	}
	return pdf, nil
}

// Canvasのsignintech/gopdf実装
type gopdfCanvas struct {
	pdf *gopdf.GoPdf
}

func (c *gopdfCanvas) Rect(x, y, w, h float64) {
	c.pdf.SetStrokeColor(0, 0, 0)
	c.pdf.SetLineWidth(0.3)
	c.pdf.SetLineType("solid")
	c.pdf.RectFromUpperLeftWithStyle(x, y, w, h, "D")
}

func (c *gopdfCanvas) Line(x1, y1, x2, y2, width float64, dashed bool) {
	c.pdf.SetStrokeColor(0, 0, 0)
	c.pdf.SetLineWidth(width)
	if dashed {
		c.pdf.SetLineType("dashed")
	} else {
		c.pdf.SetLineType("solid")
	}
	c.pdf.Line(x1, y1, x2, y2)
}

func (c *gopdfCanvas) Text(x, y float64, text string, sizePt float64, align scene.Align) {
	if text == "" {
		return
	}
	if err := c.pdf.SetFont(fontName, "", sizePt); err != nil {
		return
	}
	c.pdf.SetTextColor(0, 0, 0)
	switch align {
	case scene.AlignCenter:
		if w, err := c.pdf.MeasureTextWidth(text); err == nil {
			x -= w / 2
		}
	case scene.AlignRight:
		if w, err := c.pdf.MeasureTextWidth(text); err == nil {
			x -= w
		}
	}
	c.pdf.SetXY(x, y)
	_ = c.pdf.Cell(nil, text)
}

// 実フォントによる正確な文字幅(mm)をレイアウト層へ提供
type Measurer struct {
	pdf *gopdf.GoPdf
}

// 描画と同一フォントで採寸するMeasurerを生成
func NewMeasurer() (*Measurer, error) {
	pdf, err := newDoc()
	if err != nil {
		return nil, err
	}
	return &Measurer{pdf: pdf}, nil
}

func (m *Measurer) Measure(text string, sizePt float64) float64 {
	if text == "" {
		return 0
	}
	if err := m.pdf.SetFont(fontName, "", sizePt); err != nil {
		return 0
	}
	w, err := m.pdf.MeasureTextWidth(text)
	if err != nil {
		return 0
	}
	return w
}
