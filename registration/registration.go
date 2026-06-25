// Package registrationは相続登記申請書PDF生成の公開API
// touki.ApplicationまたはJSONからA4縦の複数ページPDFバイト列を返す
package registration

import (
	"tsugu-mcp/internal/reginput"
	"tsugu-mcp/internal/reglayout"
	"tsugu-mcp/internal/render"
	"tsugu-mcp/internal/wareki"
	"tsugu-mcp/touki"
)

// 日付の表記形式
type EraStyle = wareki.Style

const (
	EraWareki            = wareki.Wareki
	EraWarekiWithSeireki = wareki.WarekiWithSeireki
	EraSeireki           = wareki.Seireki
)

// 生成時の設定
type Options struct {
	Era EraStyle
}

// 既定設定(和暦)
func DefaultOptions() Options {
	return Options{Era: EraWareki}
}

// touki.Applicationを検証し登記申請書PDFを生成
func Generate(app touki.Application, opt Options) ([]byte, error) {
	if err := app.Validate(); err != nil {
		return nil, err
	}
	st := reglayout.DefaultStyle()
	st.Era = opt.Era
	return render.ToPDFMulti(reglayout.Build(app, st))
}

// JSONを解釈しPDF生成(MCPツール等からの利用想定)
func GenerateFromJSON(data []byte, opt Options) ([]byte, error) {
	app, err := reginput.Decode(data)
	if err != nil {
		return nil, err
	}
	return Generate(app, opt)
}
