package family

import (
	"strings"
	"testing"
)

func TestValidate_OK(t *testing.T) {
	doc := Document{
		Decedent: Decedent{Name: "甲", BirthDate: Date{1950, 1, 1}, DeathDate: Date{2025, 1, 1}},
		Children: []*Node{{Person: Person{Name: "乙", BirthDate: Date{1980, 5, 5}}}},
	}
	if err := doc.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_Errors(t *testing.T) {
	cases := []struct {
		name string
		doc  Document
		want string
	}{
		{
			"氏名なし",
			Document{Decedent: Decedent{DeathDate: Date{2025, 1, 1}}},
			"被相続人: name is required",
		},
		{
			"死亡日なし",
			Document{Decedent: Decedent{Name: "甲"}},
			"被相続人: death date is required",
		},
		{
			"死亡が出生より前",
			Document{Decedent: Decedent{Name: "甲", BirthDate: Date{2025, 1, 1}, DeathDate: Date{2000, 1, 1}}},
			"death date is before birth date",
		},
		{
			"相続人の氏名なし",
			Document{Decedent: Decedent{Name: "甲", DeathDate: Date{2025, 1, 1}}, Children: []*Node{{}}},
			"子: name is required",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.doc.Validate()
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("err = %v, want contains %q", err, c.want)
			}
		})
	}
}
