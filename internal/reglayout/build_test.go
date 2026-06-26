package reglayout

import (
	"strings"
	"testing"

	"tsugu-mcp/internal/scene"
	"tsugu-mcp/touki"
	"tsugu-mcp/ymd"
)

func app(nProps int) touki.Application {
	a := touki.Application{
		Causes:          []touki.Cause{{Date: ymd.Date{Year: 2024, Month: 6, Day: 15}, Text: "相続"}},
		Decedent:        touki.Decedent{Name: "山田 一郎"},
		Applicants:      []touki.Applicant{{Name: "山田 花子", BirthDate: ymd.Date{Year: 1950, Month: 5, Day: 5}, Contact: true}},
		ApplicationDate: ymd.Date{Year: 2024, Month: 12, Day: 10},
		Registry:        "東京法務局",
	}
	for i := 0; i < nProps; i++ {
		a.Properties = append(a.Properties, touki.Property{
			Kind: touki.Land, Number: "01", Location: "東京都千代田区一番町1番", LotNumber: "1番", LandCategory: "宅地", Area: "100.00",
		})
	}
	return a
}

func TestBuild_SinglePage(t *testing.T) {
	pages := Build(app(1), DefaultStyle())
	if len(pages) != 1 {
		t.Fatalf("pages = %d, want 1", len(pages))
	}
	if dashedEdges(pages[0]) < 4 {
		t.Errorf("1枚目に受付破線枠(4辺)が無い: %d", dashedEdges(pages[0]))
	}
}

func TestBuild_Paginates(t *testing.T) {
	const n = 20
	pages := Build(app(n), DefaultStyle())
	if len(pages) < 2 {
		t.Fatalf("pages = %d, want >=2", len(pages))
	}
	// 受付破線枠は1枚目のみ
	if dashedEdges(pages[0]) < 4 {
		t.Errorf("1枚目に受付枠が無い")
	}
	for i := 1; i < len(pages); i++ {
		if dashedEdges(pages[i]) != 0 {
			t.Errorf("%d枚目に受付枠が残っている", i+1)
		}
	}
	// 全不動産が描画される(不動産番号ラベル数の合計)
	if got := countLabel(pages, "不動産番号"); got != n {
		t.Errorf("不動産番号ラベル = %d, want %d", got, n)
	}
	// 全要素がページ内
	assertWithinPage(t, pages)
}

func dashedEdges(s scene.Scene) int {
	n := 0
	for _, e := range s.Edges {
		if e.Dashed {
			n++
		}
	}
	return n
}

func countLabel(pages []scene.Scene, substr string) int {
	n := 0
	for _, p := range pages {
		for _, l := range p.Labels {
			if strings.HasPrefix(l.Text, substr) {
				n++
			}
		}
	}
	return n
}

func assertWithinPage(t *testing.T, pages []scene.Scene) {
	t.Helper()
	for _, p := range pages {
		for _, l := range p.Labels {
			if l.Y < 0 || l.Y > p.Height {
				t.Errorf("ラベルがページ外: y=%.1f text=%q", l.Y, l.Text)
			}
		}
	}
}
