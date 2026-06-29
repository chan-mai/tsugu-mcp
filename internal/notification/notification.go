// Package notificationは相続登記の申請義務期限と相続人申告登記を案内
package notification

import (
	"fmt"

	"github.com/chan-mai/tsugu-mcp/internal/wareki"
	"github.com/chan-mai/tsugu-mcp/ymd"
)

// 案内条件
type Input struct {
	DeathDate ymd.Date // 相続開始日
	KnownDate ymd.Date // 相続開始・所有権取得を知った日(空ならDeathDate)
}

// 案内結果
type Result struct {
	Deadline     string   // 申請義務の期限
	DeadlineNote string   // 期限の根拠
	Provisional  []string // 相続人申告登記の要点
	Penalty      string   // 過料
	Note         string
}

// 義務化施行日(令和6年4月1日)
var enforcedAt = ymd.Date{Year: 2024, Month: 4, Day: 1}

// 申請義務期限と相続人申告登記の案内を返す
func Guide(in Input) Result {
	base := in.KnownDate
	if base.IsZero() {
		base = in.DeathDate
	}

	var res Result
	switch {
	case base.IsZero():
		res.Deadline = "不明(相続開始・取得を知った日が必要)"
		res.DeadlineNote = "原則、知った日から3年以内。施行日前に知っていた場合は令和9年(2027)3月31日まで"
	case before(base, enforcedAt):
		res.Deadline = "2027-03-31(令和9年3月31日)"
		res.DeadlineNote = "施行日(令和6年4月1日)前に相続・取得を知ったケース。施行日と知った日のいずれか遅い日から3年(改正法附則5条6項)"
	default:
		d := addYears(base, 3)
		res.Deadline = fmt.Sprintf("%04d-%02d-%02d(%s)", d.Year, d.Month, d.Day, wareki.Format(d.Year, d.Month, d.Day, wareki.Wareki))
		res.DeadlineNote = "相続開始・所有権取得を知った日から3年以内(不登法76条の2)。遺産分割が成立したら成立日から3年以内に別途登記が必要"
	}

	res.Provisional = []string{
		"趣旨: 3年内の相続登記が難しい場合に義務を簡易に履行する仕組み(不登法76条の3)",
		"必要書類: 申出書 / 申出人が相続人と分かる戸籍証明書 / 申出人の住所証明 /(代理は)委任状。被相続人の出生〜死亡の網羅収集や遺産分割は不要",
		"効果: 申出人は相続登記の申請義務を履行したものとみなされ過料を免れる",
		"限界: 持分(法定相続分)は登記されない報告的登記。遺産分割成立後は成立日から3年以内に内容に応じた相続登記が別途必要",
		"登録免許税は非課税。オンライン(かんたん登記供託申請)で押印・電子署名なしに申出可",
	}
	res.Penalty = "正当な理由なく怠ると10万円以下の過料(不登法164条1項)。自動賦課ではなく、登記官の催告→裁判所通知を経る"
	res.Note = "相続人申告登記は暫定的な義務履行であり、最終的な権利の登記(所有権移転)が別途必要"
	return res
}

func before(a, b ymd.Date) bool {
	if a.Year != b.Year {
		return a.Year < b.Year
	}
	if a.Month != b.Month {
		return a.Month < b.Month
	}
	return a.Day < b.Day
}

func addYears(d ymd.Date, n int) ymd.Date {
	out := ymd.Date{Year: d.Year + n, Month: d.Month, Day: d.Day}
	if d.Month == 2 && d.Day == 29 { // 閏日は2/28へ寄せる
		out.Day = 28
	}
	return out
}
