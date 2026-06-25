// Package familyは相続関係説明図の入力ドメインモデル
// 描画・JSON・PDF非依存の純粋データ構造
package family

// 相続人の相続上の取扱い(図中の注記に対応)
type Outcome int

const (
	OutcomeNone             Outcome = iota // 注記なし
	OutcomeInherit                         // 相続
	OutcomeRenounce                        // 相続放棄
	OutcomeDivision                        // 分割(遺産分割協議により取得)
	OutcomeByRepresentation                // 代襲相続
)

// 図中の注記文字列(不要なら空文字)
func (o Outcome) Label() string {
	switch o {
	case OutcomeInherit:
		return "(相)"
	case OutcomeRenounce:
		return "(相続放棄)"
	case OutcomeDivision:
		return "(分割)"
	case OutcomeByRepresentation:
		return "(代襲)"
	default:
		return ""
	}
}

// 被相続人以外の関係者
type Person struct {
	Name         string
	Relationship string // 続柄(表示用、長男・妻・弟など)
	Address      string
	BirthDate    Date
	DeathDate    *Date // nil=存命
	Outcome      Outcome
	Applicant    bool // 申立人
}

// 子孫・傍系のツリーノード(配偶者と次世代を保持可)
type Node struct {
	Person
	Spouse      *Person // 配偶者
	Descendants []*Node // 子(代襲を任意世代まで)
}

// 被相続人
type Decedent struct {
	Name               string
	RegisteredDomicile string // 本籍
	LastAddress        string // 最後の住所
	RegistryAddress    string // 登記上の住所
	BirthDate          Date
	DeathDate          Date
}

// 作成者(住所・氏名)
type Preparer struct {
	Address string
	Name    string
}

// 相続関係説明図1枚分の入力一式
type Document struct {
	Decedent   Decedent
	Spouse     *Person  // 被相続人の配偶者
	Ascendants []Person // 直系尊属(父母)
	Children   []*Node  // 第1順位
	Siblings   []*Node  // 第3順位(兄弟姉妹)
	Preparer   Preparer
	PreparedAt Date
}
