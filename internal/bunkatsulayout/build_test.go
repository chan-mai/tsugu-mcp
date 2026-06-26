package bunkatsulayout

import (
	"strings"
	"testing"

	"tsugu-mcp/bunkatsu"
	"tsugu-mcp/touki"
	"tsugu-mcp/ymd"
)

func TestWrapText(t *testing.T) {
	em := 3.5
	lines := wrapText(strings.Repeat("あ", 100), 10*em, em) // 1行約10字
	if len(lines) < 8 {
		t.Errorf("too few wrapped lines: %d", len(lines))
	}
	for _, ln := range lines {
		if len([]rune(ln)) > 12 {
			t.Errorf("line too long: %q", ln)
		}
	}
}

func TestWrapText_Kinsoku(t *testing.T) {
	// 行頭禁則: 句点が行頭に来ない
	em := 3.5
	for _, ln := range wrapText("あいうえおかきくけこ。さしすせそ", 10*em, em) {
		if strings.HasPrefix(ln, "。") {
			t.Errorf("line starts with full stop: %q", ln)
		}
	}
}

func TestBuild_Pages(t *testing.T) {
	a := bunkatsu.Agreement{
		Decedent: bunkatsu.Decedent{Name: "甲", Address: "東京都千代田区一番町1番", DeathDate: ymd.Date{Year: 2024, Month: 6, Day: 15}},
		Heirs:    []bunkatsu.Heir{{Name: "乙", Address: "東京都"}, {Name: "丙", Address: "神奈川県"}},
		Allocations: []bunkatsu.Allocation{{
			Acquirers:  []bunkatsu.Acquirer{{Name: "乙", Share: "2分の1"}},
			Properties: []touki.Property{{Kind: touki.Land, Number: "1", Location: "東京都千代田区一番町1番", LotNumber: "1番", LandCategory: "宅地", Area: "100.00"}},
		}},
		AgreedDate: ymd.Date{Year: 2024, Month: 12, Day: 10},
	}
	pages := Build(a, DefaultStyle())
	if len(pages) == 0 {
		t.Fatal("ページが生成されない")
	}
	for _, p := range pages {
		for _, l := range p.Labels {
			if l.Y < 0 || l.Y > p.Height {
				t.Errorf("label out of page: y=%.1f %q", l.Y, l.Text)
			}
		}
	}
}
