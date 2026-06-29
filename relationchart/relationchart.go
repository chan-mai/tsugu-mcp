// Package relationchart は相続関係説明図PDF生成の公開API
// family.DocumentまたはJSONからA4縦1ページのPDFバイト列を返す
package relationchart

import (
	"github.com/chan-mai/tsugu-mcp/family"
	"github.com/chan-mai/tsugu-mcp/internal/inputjson"
	"github.com/chan-mai/tsugu-mcp/internal/layout"
	"github.com/chan-mai/tsugu-mcp/internal/render"
	"github.com/chan-mai/tsugu-mcp/internal/wareki"
)

// 日付の表記形式
type EraStyle = wareki.Style

const (
	EraWarekiWithSeireki = wareki.WarekiWithSeireki // 令和7年(2025年)3月15日
	EraWareki            = wareki.Wareki            // 令和7年3月15日
	EraSeireki           = wareki.Seireki           // 2025年3月15日
)

// 生成時の設定
type Options struct {
	Era EraStyle
}

// 既定設定(和暦のみ・伝統的様式)
func DefaultOptions() Options {
	return Options{Era: EraWareki}
}

// family.Documentを検証し相続関係説明図PDFを生成
func Generate(doc family.Document, opt Options) ([]byte, error) {
	if err := doc.Validate(); err != nil {
		return nil, err
	}
	m, err := render.NewMeasurer()
	if err != nil {
		return nil, err
	}
	st := layout.DefaultStyle()
	st.EraStyle = opt.Era
	return render.ToPDF(layout.Build(doc, st, m))
}

// JSONを解釈しPDF生成(MCPツール等からの利用想定)
func GenerateFromJSON(data []byte, opt Options) ([]byte, error) {
	doc, err := inputjson.Decode(data)
	if err != nil {
		return nil, err
	}
	return Generate(doc, opt)
}
