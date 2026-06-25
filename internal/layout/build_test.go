package layout

import (
	"strings"
	"testing"

	"tsugu-mcp/family"
	"tsugu-mcp/internal/scene"
)

func date(y, m, d int) family.Date { return family.Date{Year: y, Month: m, Day: d} }

func baseDecedent() family.Decedent {
	return family.Decedent{
		Name:               "甲野 太郎",
		RegisteredDomicile: "東京都千代田区",
		LastAddress:        "東京都千代田区",
		BirthDate:          date(1950, 1, 1),
		DeathDate:          date(2025, 1, 1),
	}
}

// HeuristicMeasurerで実フォント無しに決定的レイアウトを検証
func build(doc family.Document) scene.Scene { return Build(doc, DefaultStyle(), HeuristicMeasurer{}) }

func TestBuild_SpouseChildren(t *testing.T) {
	doc := family.Document{
		Decedent: baseDecedent(),
		Spouse:   &family.Person{Name: "甲野 花子", Relationship: "妻", Outcome: family.OutcomeInherit},
		Children: []*family.Node{
			{Person: family.Person{Name: "甲野 一郎", Relationship: "長男", Outcome: family.OutcomeInherit}},
			{Person: family.Person{Name: "乙川 良子", Relationship: "長女", Outcome: family.OutcomeDivision}},
		},
	}
	s := build(doc)

	if got := personBoxes(s); got != 4 {
		t.Fatalf("人物カード = %d, want 4 (被相続人・配偶者・子2)", got)
	}
	if got := countDouble(s); got != 1 {
		t.Errorf("婚姻二重線 = %d, want 1", got)
	}

	dec := findBox(t, s, "被")
	son := findBox(t, s, "長男")
	spouse := findBox(t, s, "妻")
	if son.X <= dec.X {
		t.Errorf("子は被相続人より右(次世代列)であるべき: 子 x=%.1f 被相続人 x=%.1f", son.X, dec.X)
	}
	if spouse.Y <= dec.Y {
		t.Errorf("配偶者は被相続人の下に積むべき: 配偶者 y=%.1f 被相続人 y=%.1f", spouse.Y, dec.Y)
	}
	assertWithinPage(t, s)
}

func TestBuild_Representation(t *testing.T) {
	dead := date(2020, 1, 1)
	doc := family.Document{
		Decedent: baseDecedent(),
		Children: []*family.Node{{
			Person: family.Person{Name: "甲野 二郎", Relationship: "二男", DeathDate: &dead},
			Spouse: &family.Person{Name: "甲野 月子", Relationship: "妻"},
			Descendants: []*family.Node{
				{Person: family.Person{Name: "甲野 孫一", Relationship: "孫", Outcome: family.OutcomeByRepresentation}},
			},
		}},
	}
	s := build(doc)

	child := findBox(t, s, "二男")
	grandchild := findBox(t, s, "孫")
	if grandchild.X <= child.X {
		t.Errorf("代襲の孫は親より右の世代であるべき: 孫 x=%.1f 親 x=%.1f", grandchild.X, child.X)
	}
	assertWithinPage(t, s)
}

func TestBuild_Ascendants(t *testing.T) {
	doc := family.Document{
		Decedent: baseDecedent(),
		Ascendants: []family.Person{
			{Name: "甲野 祖父", Relationship: "父"},
			{Name: "甲野 祖母", Relationship: "母"},
		},
	}
	s := build(doc)

	if got := personBoxes(s); got != 3 {
		t.Fatalf("人物カード = %d, want 3 (被相続人・父・母)", got)
	}
	father := findBox(t, s, "父")
	dec := findBox(t, s, "被")
	if dec.X <= father.X {
		t.Errorf("被相続人は尊属より右の世代であるべき: 被相続人 x=%.1f 父 x=%.1f", dec.X, father.X)
	}
	assertWithinPage(t, s)
}

func TestBuild_DecedentOnly(t *testing.T) {
	s := build(family.Document{Decedent: baseDecedent()})
	if got := personBoxes(s); got != 1 {
		t.Fatalf("人物カード = %d, want 1", got)
	}
	assertWithinPage(t, s)
}

func personBoxes(s scene.Scene) int {
	n := 0
	for _, b := range s.Boxes {
		if !b.Border { // 枠なし = 人物カード
			n++
		}
	}
	return n
}

func countDouble(s scene.Scene) int {
	n := 0
	for _, e := range s.Edges {
		if e.Double {
			n++
		}
	}
	return n
}

func findBox(t *testing.T, s scene.Scene, substr string) scene.Box {
	t.Helper()
	for _, b := range s.Boxes {
		for _, ln := range b.Lines {
			if strings.Contains(ln, substr) {
				return b
			}
		}
	}
	t.Fatalf("%q を含む欄が見つかりません", substr)
	return scene.Box{}
}

func assertWithinPage(t *testing.T, s scene.Scene) {
	t.Helper()
	const eps = 0.01
	for _, b := range s.Boxes {
		if b.X < -eps || b.Y < -eps || b.X+b.W > s.Width+eps || b.Y+b.H > s.Height+eps {
			t.Errorf("欄がページ外: %+v (page %.0fx%.0f)", b, s.Width, s.Height)
		}
	}
}
