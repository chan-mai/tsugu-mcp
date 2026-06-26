// Package bunkatsuは遺産分割協議書の入力ドメインモデル
// 描画・JSON・PDF非依存の純粋データ構造
package bunkatsu

import (
	"tsugu-mcp/touki"
	"tsugu-mcp/ymd"
)

// 被相続人(協議書は死亡日+住所+氏名で特定)
type Decedent struct {
	Name      string
	Address   string // 最後の住所(本籍を併記する実務もある)
	DeathDate ymd.Date
}

// 共同相続人(冒頭の列挙と署名欄に使う)
type Heir struct {
	Name    string
	Address string
}

// 取得者と持分
type Acquirer struct {
	Name  string
	Share string // 持分(例 2分の1、単独なら空)
}

// 取得の対応(取得者→財産)
type Allocation struct {
	Acquirers  []Acquirer
	Properties []touki.Property // 不動産(touki.Propertyを再利用)
	Items      []string         // 不動産以外の財産(預貯金等、任意のフリー記述)
}

// 遺産分割協議書1件分の入力一式
type Agreement struct {
	Decedent    Decedent
	Heirs       []Heir // 共同相続人全員
	Allocations []Allocation
	AgreedDate  ymd.Date // 協議成立日
	Copies      int      // 作成通数(0なら相続人数)
}
