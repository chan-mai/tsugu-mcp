// Package agreementは遺産分割協議書PDF生成の公開API
// bunkatsu.AgreementまたはJSONからA4縦の複数ページPDFバイト列を返す
package agreement

import (
	"github.com/chan-mai/tsugu-mcp/bunkatsu"
	"github.com/chan-mai/tsugu-mcp/internal/bunkatsuinput"
	"github.com/chan-mai/tsugu-mcp/internal/bunkatsulayout"
	"github.com/chan-mai/tsugu-mcp/internal/render"
	"github.com/chan-mai/tsugu-mcp/internal/wareki"
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

// bunkatsu.Certificateを検証し遺産分割協議証明書PDFを生成(相続人ごとに1ページ)
func GenerateCertificate(c bunkatsu.Certificate, opt Options) ([]byte, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	st := bunkatsulayout.DefaultStyle()
	st.Era = opt.Era
	return render.ToPDFMulti(bunkatsulayout.BuildCertificate(c, st))
}

// JSONを解釈し協議証明書PDFを生成
func GenerateCertificateFromJSON(data []byte, opt Options) ([]byte, error) {
	c, err := bunkatsuinput.DecodeCertificate(data)
	if err != nil {
		return nil, err
	}
	return GenerateCertificate(c, opt)
}
