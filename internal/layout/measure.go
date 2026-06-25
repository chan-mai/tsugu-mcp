package layout

// フォントサイズ(pt)でのテキスト描画幅(mm)を返す抽象。レイアウトを描画ライブラリから分離
type Measurer interface {
	Measure(text string, sizePt float64) float64
}

// 全角=1em・半角=0.5emで概算する既定実装(実フォント無しのレイアウト単体テスト用)
type HeuristicMeasurer struct{}

func (HeuristicMeasurer) Measure(text string, sizePt float64) float64 {
	em := sizePt * ptToMM
	var w float64
	for _, r := range text {
		if isHalfWidth(r) {
			w += em * 0.5
		} else {
			w += em
		}
	}
	return w
}

// 半角(Latin-1以下)概算判定
func isHalfWidth(r rune) bool { return r < 0x0100 }
