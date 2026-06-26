// Package toukiは相続登記申請書の入力ドメインモデル
// 描画・JSON・PDF非依存の純粋データ構造
package touki

import "tsugu-mcp/ymd"

// 不動産の種別
type PropertyKind int

const (
	Land     PropertyKind = iota // 土地
	Building                     // 建物
)

// 原因(数次相続で複数)
type Cause struct {
	Date ymd.Date
	Text string // 例: 山田一郎相続 / 相続
}

// 相続人欄の被相続人
type Decedent struct {
	Name string
}

// 申請人(複数+持分)
type Applicant struct {
	Name      string
	Address   string
	Share     string // 持分(単独なら空)
	NameKana  string // 氏名ふりがな
	BirthDate ymd.Date
	Email     string
	Phone     string
	Contact   bool // 連絡先電話を表示する代表者
}

// 不動産の表示(土地または建物)
type Property struct {
	Kind     PropertyKind
	Number   string // 不動産番号
	Location string // 所在

	// 土地
	LotNumber    string // 地番
	LandCategory string // 地目
	Area         string // 地積

	// 建物
	HouseNumber  string // 家屋番号
	BuildingType string // 種類
	Structure    string // 構造
	FloorArea    string // 床面積
}

// 登記申請書1件分の入力一式(登記の目的=所有権移転は固定のため持たない)
type Application struct {
	Causes          []Cause
	Decedent        Decedent
	Applicants      []Applicant
	Attachments     []string // 添付情報
	DeclineIDInfo   bool     // 登記識別情報の通知を希望しない欄のチェック(falseでも空欄□を描画)
	ApplicationDate ymd.Date
	Registry        string // 法務局
	TaxValue        string // 課税価格(表示文字列)
	RegistrationTax string // 登録免許税(表示文字列)
	Properties      []Property
}
