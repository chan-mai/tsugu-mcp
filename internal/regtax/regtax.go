// Package regtaxは相続登記の登録免許税を計算する(根拠: docs/knowledge/04)
// 相続による所有権移転(税率0.4%)の端数処理順序と免税措置の文言を扱う
package regtax

import (
	"fmt"
	"sort"
)

const (
	rateNumInheritance = 4    // 相続: 1000分の4(0.4%)
	rateDen            = 1000 // 税率の分母
	smallValueLimit    = 1_000_000
	minTax             = 1000

	stmtIntermediate = "租税特別措置法第84条の2の2第1項により非課税"
	stmtSmallValue   = "租税特別措置法第84条の2の2第2項により非課税"

	deadlineNote = "免税措置の適用期限は令和9年(2027)3月31日。期限後は要再確認。免税は申請書への条文記載が必要(自動適用されない)。"
)

// PrivateRoadは私道(評価額0)の認定価額パラメータ
type PrivateRoad struct {
	NeighborUnitPrice int     // 近傍宅地の1㎡単価(円)
	Area              float64 // 私道の地積(㎡)
}

// Propertyは1不動産の入力
type Property struct {
	Kind        string       // land | building
	Value       int          // 固定資産評価額(円)、私道はPrivateRoadで認定
	ShareNum    int          // 持分の分子
	ShareDen    int          // 持分の分母(0は1/1)
	Exemption   string       // "" 自動 / none / small_value / intermediate
	PrivateRoad *PrivateRoad // 私道認定(任意)
}

// Inputは計算入力(原因は相続0.4%固定)
type Input struct {
	Properties []Property
}

// Lineは物件別の計算結果
type Line struct {
	Index     int
	TaxValue  int    // 課税価格(切捨前の円)
	Taxable   bool   // 課税対象か
	Exemption string // 適用免税 small_value|intermediate(空は課税)
}

// Resultは計算結果
type Result struct {
	Lines            []Line
	TaxableTotal     int      // 課税標準(合算後1000円未満切捨)
	Tax              int      // 登録免許税(円)
	ExemptStatements []string // 適用免税の申請書記載文言
	EligibilityNotes []string // 免税対象になり得る物件への注意(適用は要記載)
	Note             string   // 適用期限等
}

// 登録免許税を計算(免税は明示指定時のみ適用し自動適用しない)
func Calculate(in Input) Result {
	var res Result
	var taxableSum float64
	stmts := map[string]bool{}

	for i, p := range in.Properties {
		v := propValue(p)
		line := Line{Index: i, TaxValue: int(v)}

		ex := p.Exemption
		if ex == "" && p.Kind == "land" && v <= smallValueLimit {
			res.EligibilityNotes = append(res.EligibilityNotes,
				fmt.Sprintf("物件%d: 課税価格が100万円以下のため措置2(第84条の2の2第2項)の対象になり得る。適用するには申請書に条文記載が必要", i+1))
		}

		switch ex {
		case "small_value":
			line.Exemption = "small_value"
			stmts[stmtSmallValue] = true
		case "intermediate":
			line.Exemption = "intermediate"
			stmts[stmtIntermediate] = true
		default:
			line.Taxable = true
			taxableSum += v
		}
		res.Lines = append(res.Lines, line)
	}

	res.TaxableTotal = floorTo(int(taxableSum), 1000)
	if res.TaxableTotal > 0 {
		tax := res.TaxableTotal * rateNumInheritance / rateDen
		tax = floorTo(tax, 100)
		if tax < minTax {
			tax = minTax
		}
		res.Tax = tax
	}

	for s := range stmts {
		res.ExemptStatements = append(res.ExemptStatements, s)
	}
	sort.Strings(res.ExemptStatements)
	res.Note = deadlineNote
	return res
}

func propValue(p Property) float64 {
	share := 1.0
	if p.ShareDen != 0 {
		share = float64(p.ShareNum) / float64(p.ShareDen)
	}
	if p.PrivateRoad != nil {
		return float64(p.PrivateRoad.NeighborUnitPrice) * p.PrivateRoad.Area * 0.3 * share
	}
	return float64(p.Value) * share
}

func floorTo(v, unit int) int { return v / unit * unit }
