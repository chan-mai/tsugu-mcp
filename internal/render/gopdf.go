package render

import (
	"fmt"

	"github.com/signintech/gopdf"

	"github.com/chan-mai/tsugu-mcp/assets"
	"github.com/chan-mai/tsugu-mcp/internal/scene"
)

const fontName = "ipaexg"

// scene.SceneをA4縦1ページのPDFバイト列へ描画
func ToPDF(s scene.Scene) ([]byte, error) {
	return ToPDFMulti([]scene.Scene{s})
}

// 複数sceneをページごとにA4縦PDFへ描画
func ToPDFMulti(pages []scene.Scene) ([]byte, error) {
	pdf, err := startDoc()
	if err != nil {
		return nil, err
	}
	c := &gopdfCanvas{pdf: pdf}
	for _, s := range pages {
		pdf.AddPage()
		Draw(s, c)
	}
	return pdf.GetBytesPdf(), nil
}

// Start+フォント読込のみ(ページ追加は呼び出し側)
func startDoc() (*gopdf.GoPdf, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: gopdf.UnitMM, PageSize: *gopdf.PageSizeA4})
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

func (c *gopdfCanvas) Line(x1, y1, x2, y2, width float64, dashed bool, color scene.RGB) {
	c.pdf.SetStrokeColor(color.R, color.G, color.B)
	c.pdf.SetLineWidth(width)
	if dashed {
		c.pdf.SetCustomLineType([]float64{1.0, 1.0}, 0) // 細かい破線
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
	pdf, err := startDoc()
	if err != nil {
		return nil, err
	}
	pdf.AddPage() // MeasureTextWidthは現在ページが必要
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
