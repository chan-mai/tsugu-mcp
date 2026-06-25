// Package assets は同梱フォント等の埋め込みリソース
package assets

import _ "embed"

// 埋め込みIPAexゴシック(TrueType)フォントデータ
// PDFサブセット埋め込み用、ライセンスはfonts/IPA_Font_License_Agreement_v1.0.txt
//
//go:embed fonts/ipaexg.ttf
var IPAexGothic []byte
