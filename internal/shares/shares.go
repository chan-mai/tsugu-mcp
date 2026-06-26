// Package sharesは法定相続分を算定する(根拠: 民法900・901条ほか)
// 現行(昭和56年1月1日以降)の規定を既約分数で計算し、配分の助言はしない
package shares

import (
	"math/big"

	"tsugu-mcp/ymd"
)

// 血族相続人(子・直系尊属・兄弟姉妹のいずれか)
type Heir struct {
	Name         string
	Alive        bool   // 生存
	Renounced    bool   // 相続放棄(いなかったとみなす、代襲なし)
	Disqualified bool   // 欠格・廃除(本人脱落だが代襲あり)
	HalfBlood    bool   // 半血(兄弟姉妹のみ、相続分は全血の1/2)
	Substitutes  []Heir // 代襲者(子→孫、兄弟姉妹→甥姪)
}

// 算定入力(順位は子→直系尊属→兄弟姉妹の最先順位のみ相続)
type Input struct {
	DeathDate  ymd.Date
	SpouseName string // 配偶者の氏名(空なら配偶者なし)
	Children   []Heir // 第1順位
	Ascendants []Heir // 第2順位(直系尊属)
	Siblings   []Heir // 第3順位
}

// 1相続人の割当
type Share struct {
	Name     string
	Role     string // 配偶者 / 子 / 孫等(代襲) / 直系尊属 / 兄弟姉妹 / 甥姪(代襲)
	Fraction string // 既約分数(例 1/2)
}

// 算定結果
type Result struct {
	Shares   []Share
	Sum      string // 検算(相続人がいれば1)
	Warnings []string
	Note     string
}

const note = "これは法定相続分(法律上の原則的割合)。実際の取得割合は遺産分割協議で自由に定められる。特別受益・寄与分・遺留分・非嫡出子の時点区分は考慮しない。"

// 法定相続分を算定する
func Calculate(in Input) Result {
	var res Result
	res.Note = note
	if in.DeathDate.IsZero() {
		res.Warnings = append(res.Warnings, "death date not given; assumed current law (from 1981-01-01)")
	} else if before(in.DeathDate, 1981, 1, 1) {
		res.Warnings = append(res.Warnings, "death before 1981-01-01: old law applies and shares differ (expert review required)")
	}

	hasChildren := anyLives(in.Children)
	hasAscendants := anyLives(in.Ascendants)
	hasSiblings := anyLives(in.Siblings)
	spouse := in.SpouseName != ""

	var blood []Share
	var spouseFrac *big.Rat
	switch {
	case hasChildren:
		spouseFrac = ratOrNil(spouse, 1, 2)
		distribute(in.Children, bloodTotal(spouse, 1, 2), "子", "孫等(代襲)", weightOne, &blood)
	case hasAscendants:
		spouseFrac = ratOrNil(spouse, 2, 3)
		distribute(in.Ascendants, bloodTotal(spouse, 2, 3), "直系尊属", "直系尊属", weightOne, &blood)
	case hasSiblings:
		spouseFrac = ratOrNil(spouse, 3, 4)
		distribute(in.Siblings, bloodTotal(spouse, 3, 4), "兄弟姉妹", "甥姪(代襲)", weightSibling, &blood)
	default:
		if spouse {
			spouseFrac = big.NewRat(1, 1)
		} else {
			res.Warnings = append(res.Warnings, "no heir matches (check inputs)")
		}
	}

	if spouseFrac != nil {
		res.Shares = append(res.Shares, Share{Name: in.SpouseName, Role: "配偶者", Fraction: spouseFrac.RatString()})
	}
	res.Shares = append(res.Shares, blood...)

	sum := new(big.Rat)
	for _, s := range res.Shares {
		r := new(big.Rat)
		r.SetString(s.Fraction)
		sum.Add(sum, r)
	}
	res.Sum = sum.RatString()
	return res
}

// 血族側の合計割合(配偶者ありなら1-配偶者分、なしなら全部)
func bloodTotal(spouse bool, num, den int64) *big.Rat {
	if !spouse {
		return big.NewRat(1, 1)
	}
	return new(big.Rat).Sub(big.NewRat(1, 1), big.NewRat(num, den))
}

func ratOrNil(spouse bool, num, den int64) *big.Rat {
	if !spouse {
		return nil
	}
	return big.NewRat(num, den)
}

func weightOne(Heir) int64 { return 1 }

func weightSibling(h Heir) int64 {
	if h.HalfBlood {
		return 1
	}
	return 2
}

// groupの相続人へtotalを重み付き配分(代襲は被代襲者の取り分を代襲者で頭割り)
func distribute(heirs []Heir, total *big.Rat, role, subRole string, weight func(Heir) int64, out *[]Share) {
	var totalW int64
	for _, h := range heirs {
		if branchLives(h) {
			totalW += weight(h)
		}
	}
	if totalW == 0 {
		return
	}
	for _, h := range heirs {
		if !branchLives(h) {
			continue
		}
		share := new(big.Rat).Mul(total, big.NewRat(weight(h), totalW))
		if h.Alive && !h.Disqualified {
			*out = append(*out, Share{Name: h.Name, Role: role, Fraction: share.RatString()})
		} else {
			distribute(h.Substitutes, share, subRole, subRole, weightOne, out)
		}
	}
}

// 相続権が生きている枝か(放棄=死、欠格/廃除・死亡は代襲の生存で判定)
func branchLives(h Heir) bool {
	if h.Renounced {
		return false
	}
	if h.Alive && !h.Disqualified {
		return true
	}
	for _, s := range h.Substitutes {
		if branchLives(s) {
			return true
		}
	}
	return false
}

func anyLives(hs []Heir) bool {
	for _, h := range hs {
		if branchLives(h) {
			return true
		}
	}
	return false
}

func before(d ymd.Date, y, m, day int) bool {
	if d.Year != y {
		return d.Year < y
	}
	if d.Month != m {
		return d.Month < m
	}
	return d.Day < day
}
