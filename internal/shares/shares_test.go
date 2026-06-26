package shares

import (
	"testing"

	"tsugu-mcp/ymd"
)

func got(r Result) map[string]string {
	m := map[string]string{}
	for _, s := range r.Shares {
		m[s.Name] = s.Fraction
	}
	return m
}

func eq(t *testing.T, r Result, want map[string]string) {
	t.Helper()
	g := got(r)
	if len(g) != len(want) {
		t.Fatalf("shares = %v, want %v", g, want)
	}
	for k, v := range want {
		if g[k] != v {
			t.Errorf("%s = %q, want %q (all: %v)", k, g[k], v, g)
		}
	}
	if r.Sum != "1" {
		t.Errorf("sum = %q, want 1", r.Sum)
	}
}

// 民法900・901条/docs/knowledge/07の計算例に一致すること
func TestCalculate(t *testing.T) {
	t.Run("配偶者+子3人", func(t *testing.T) {
		eq(t, Calculate(Input{SpouseName: "配", Children: []Heir{
			{Name: "子1", Alive: true}, {Name: "子2", Alive: true}, {Name: "子3", Alive: true},
		}}), map[string]string{"配": "1/2", "子1": "1/6", "子2": "1/6", "子3": "1/6"})
	})

	t.Run("配偶者+父母", func(t *testing.T) {
		eq(t, Calculate(Input{SpouseName: "配", Ascendants: []Heir{
			{Name: "父", Alive: true}, {Name: "母", Alive: true},
		}}), map[string]string{"配": "2/3", "父": "1/6", "母": "1/6"})
	})

	t.Run("子のみ2人", func(t *testing.T) {
		eq(t, Calculate(Input{Children: []Heir{
			{Name: "子1", Alive: true}, {Name: "子2", Alive: true},
		}}), map[string]string{"子1": "1/2", "子2": "1/2"})
	})

	t.Run("配偶者+全血兄2+半血弟1", func(t *testing.T) {
		eq(t, Calculate(Input{SpouseName: "配", Siblings: []Heir{
			{Name: "兄1", Alive: true}, {Name: "兄2", Alive: true}, {Name: "弟", Alive: true, HalfBlood: true},
		}}), map[string]string{"配": "3/4", "兄1": "1/10", "兄2": "1/10", "弟": "1/20"})
	})

	t.Run("兄弟のみ全血2半血1", func(t *testing.T) {
		eq(t, Calculate(Input{Siblings: []Heir{
			{Name: "兄1", Alive: true}, {Name: "兄2", Alive: true}, {Name: "弟", Alive: true, HalfBlood: true},
		}}), map[string]string{"兄1": "2/5", "兄2": "2/5", "弟": "1/5"})
	})

	t.Run("代襲(子B先死亡で孫D・E)", func(t *testing.T) {
		eq(t, Calculate(Input{SpouseName: "配", Children: []Heir{
			{Name: "A", Alive: true},
			{Name: "B", Alive: false, Substitutes: []Heir{{Name: "D", Alive: true}, {Name: "E", Alive: true}}},
		}}), map[string]string{"配": "1/2", "A": "1/4", "D": "1/8", "E": "1/8"})
	})

	t.Run("配偶者のみ", func(t *testing.T) {
		eq(t, Calculate(Input{SpouseName: "配"}), map[string]string{"配": "1"})
	})

	t.Run("子の放棄は代襲せず他の子へ", func(t *testing.T) {
		eq(t, Calculate(Input{SpouseName: "配", Children: []Heir{
			{Name: "A", Renounced: true}, {Name: "B", Alive: true},
		}}), map[string]string{"配": "1/2", "B": "1/2"})
	})

	t.Run("子全員放棄で第2順位へ", func(t *testing.T) {
		eq(t, Calculate(Input{SpouseName: "配",
			Children:   []Heir{{Name: "A", Renounced: true}},
			Ascendants: []Heir{{Name: "父", Alive: true}, {Name: "母", Alive: true}},
		}), map[string]string{"配": "2/3", "父": "1/6", "母": "1/6"})
	})
}

func TestCalculate_OldLawWarning(t *testing.T) {
	r := Calculate(Input{DeathDate: ymd.Date{Year: 1980, Month: 5, Day: 1}, Children: []Heir{{Name: "子", Alive: true}}})
	if len(r.Warnings) == 0 {
		t.Error("pre-1981 death should warn")
	}
}
