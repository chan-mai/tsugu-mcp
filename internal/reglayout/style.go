// Package reglayoutはtouki.Applicationを登記申請書の描画指示([]scene.Scene)へ変換
// 項目を縦に流し、不動産が多い場合はページ送り。純粋な幾何のみでPDF非依存
package reglayout

import (
	"github.com/chan-mai/tsugu-mcp/internal/scene"
	"github.com/chan-mai/tsugu-mcp/internal/wareki"
)

// pt→mm換算(1pt=1/72inch)
const ptToMM = 25.4 / 72.0

// レイアウトの寸法・書式パラメータ(座標mm・フォントpt)
type Style struct {
	PageW, PageH float64
	MarginX      float64
	MarginTop    float64
	MarginBottom float64

	ReceiptH     float64   // 受付番号表の枠高(1枚目上部)
	ReceiptInset float64   // 受付枠の左右インセット
	ReceiptColor scene.RGB // 受付枠の線色

	TitlePt float64
	BodyPt  float64
	LineH   float64 // 行高

	LabelChars  int     // ラベル均等割付の基準文字数
	ValueGap    float64 // ラベル列と値の間隔
	PropIndent  float64 // 不動産ブロックの字下げ
	BlockGap    float64 // 不動産ブロック間の間隔
	TableLabelW float64 // 申請人表の左列幅

	LineWidth float64
	Era       wareki.Style
}

// A4縦・登記申請書向けの既定スタイル
func DefaultStyle() Style {
	return Style{
		PageW: 210, PageH: 297,
		MarginX: 22, MarginTop: 18, MarginBottom: 18,

		ReceiptH: 48, ReceiptInset: 26,
		ReceiptColor: scene.RGB{R: 0xbc, G: 0xe2, B: 0xe8}, // 水色(#bce2e8)

		TitlePt: 14,
		BodyPt:  10.5,
		LineH:   7,

		LabelChars:  5,
		ValueGap:    8,
		PropIndent:  4,
		BlockGap:    3,
		TableLabelW: 28,

		LineWidth: 0.3,
		Era:       wareki.Wareki,
	}
}

func (s Style) em() float64      { return s.BodyPt * ptToMM }
func (s Style) labelW() float64  { return float64(s.LabelChars) * s.em() }
func (s Style) valueX() float64  { return s.MarginX + s.labelW() + s.ValueGap }
func (s Style) bottomY() float64 { return s.PageH - s.MarginBottom }
