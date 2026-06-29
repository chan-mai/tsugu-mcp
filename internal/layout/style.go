package layout

import "github.com/chan-mai/tsugu-mcp/internal/wareki"

// pt→mm換算(1pt=1/72inch)
const ptToMM = 25.4 / 72.0

// レイアウトの寸法・書式パラメータ(座標mm・フォントpt)
type Style struct {
	PageW, PageH float64 // ページサイズ
	Margin       float64 // ページ余白

	TitlePt  float64 // 表題
	TitleGap float64 // 表題下の間隔
	BodyPt   float64 // 人物・関係の本文
	FootPt   float64 // 作成欄

	LineSpacing float64 // 行送り係数
	LineWidth   float64 // 関係線の幅

	GenGap         float64 // 世代(列)間の水平間隔
	SiblingVGap    float64 // 兄弟姉妹間の垂直間隔
	CoupleVGap     float64 // 夫婦の上下間隔
	MarriageIndent float64 // 婚姻二重線の名前下インデント
	MarriageClear  float64 // 婚姻二重線の端と文字との余白
	BracketGap     float64 // 子ブラケットと子カラムの間隔
	FooterGap      float64 // 作図領域と作成欄の間隔

	MinScale float64 // AutoFitの縮小下限

	EraStyle wareki.Style
}

// A4縦・伝統的様式向けの既定スタイル
func DefaultStyle() Style {
	return Style{
		PageW: 210, PageH: 297,
		Margin: 15,

		TitlePt:  15,
		TitleGap: 9,
		BodyPt:   9,
		FootPt:   9,

		LineSpacing: 1.3,
		LineWidth:   0.3,

		GenGap:         14,
		SiblingVGap:    3.5,
		CoupleVGap:     7,
		MarriageIndent: 6,
		MarriageClear:  1.5,
		BracketGap:     5,
		FooterGap:      6,

		MinScale: 0.4,

		EraStyle: wareki.WarekiWithSeireki,
	}
}

func (s Style) lineHeight() float64 { return s.BodyPt * ptToMM * s.LineSpacing }
