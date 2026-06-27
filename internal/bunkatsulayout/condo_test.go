package bunkatsulayout

import (
	"strings"
	"testing"

	"tsugu-mcp/bunkatsu"
	"tsugu-mcp/touki"
	"tsugu-mcp/ymd"
)

func TestPropertyBlock_Condominium(t *testing.T) {
	a := bunkatsu.Agreement{
		Decedent: bunkatsu.Decedent{Name: "甲", DeathDate: ymd.Date{Year: 2024, Month: 6, Day: 15}},
		Heirs:    []bunkatsu.Heir{{Name: "乙"}},
		Allocations: []bunkatsu.Allocation{{
			Acquirers: []bunkatsu.Acquirer{{Name: "乙"}},
			Properties: []touki.Property{{
				Kind:         touki.Condominium,
				Location:     "東京都千代田区一番町23番地",
				BuildingName: "一番町マンション",
				HouseNumber:  "一番町23番の301",
				UnitName:     "301号",
				FloorArea:    "3階部分 60.12平方メートル",
				LandRights:   []touki.LandRight{{Symbol: "1", LocationLot: "一番町23番", Area: "500.00", RightType: "所有権", RightShare: "1000分の35"}},
			}},
		}},
	}
	var texts []string
	for _, p := range Build(a, DefaultStyle()) {
		for _, l := range p.Labels {
			texts = append(texts, l.Text)
		}
	}
	joined := strings.Join(texts, "\n")
	for _, want := range []string{"一棟の建物の表示", "敷地権の表示", "一番町マンション", "1000分の35"} {
		if !strings.Contains(joined, want) {
			t.Errorf("condominium output missing %q", want)
		}
	}
}
