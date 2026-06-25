// Package wareki は西暦年月日を和暦表記へ変換
// 元号テーブル駆動の純粋関数、外部依存なし
package wareki

import "fmt"

// 日付の出力形式
type Style int

const (
	WarekiWithSeireki Style = iota // 令和7年(2025年)3月15日
	Wareki                         // 令和7年3月15日
	Seireki                        // 2025年3月15日
)

type era struct {
	name           string
	year, mon, day int // 改元日(グレゴリオ暦)
}

// 新しい元号順に並べ、境界判定で最初に該当したものを採用
var eras = []era{
	{"令和", 2019, 5, 1},
	{"平成", 1989, 1, 8},
	{"昭和", 1926, 12, 25},
	{"大正", 1912, 7, 30},
	{"明治", 1868, 1, 25},
}

// year/month/dayを指定形式の文字列へ
// year==0は空文字、明治より前は元号判定不能のため西暦表記
func Format(year, month, day int, style Style) string {
	if year == 0 {
		return ""
	}
	mmdd := fmt.Sprintf("%d月%d日", month, day)
	seireki := fmt.Sprintf("%d年%s", year, mmdd)
	if style == Seireki {
		return seireki
	}

	e, ok := lookup(year, month, day)
	if !ok {
		return seireki
	}
	eraYear := "元年"
	if y := year - e.year + 1; y != 1 {
		eraYear = fmt.Sprintf("%d年", y)
	}
	if style == Wareki {
		return fmt.Sprintf("%s%s%s", e.name, eraYear, mmdd)
	}
	return fmt.Sprintf("%s%s(%d年)%s", e.name, eraYear, year, mmdd)
}

func lookup(year, month, day int) (era, bool) {
	key := year*10000 + month*100 + day
	for _, e := range eras {
		if key >= e.year*10000+e.mon*100+e.day {
			return e, true
		}
	}
	return era{}, false
}
