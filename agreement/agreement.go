// Package agreementは遺産分割協議書PDF生成の公開API
// bunkatsu.AgreementまたはJSONからA4縦の複数ページPDFバイト列を返す
package agreement

import (
	"tsugu-mcp/bunkatsu"
	"tsugu-mcp/internal/bunkatsuinput"
	"tsugu-mcp/internal/bunkatsulayout"
	"tsugu-mcp/internal/render"
	"tsugu-mcp/internal/wareki"
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

// bunkatsu.Agreementを検証し遺産分割協議書PDFを生成
func Generate(a bunkatsu.Agreement, opt Options) ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	st := bunkatsulayout.DefaultStyle()
	st.Era = opt.Era
	return render.ToPDFMulti(bunkatsulayout.Build(a, st))
}

// JSONを解釈しPDF生成(MCPツール等からの利用想定)
func GenerateFromJSON(data []byte, opt Options) ([]byte, error) {
	a, err := bunkatsuinput.Decode(data)
	if err != nil {
		return nil, err
	}
	return Generate(a, opt)
}
