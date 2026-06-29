package reglayout

import (
	"strings"
	"testing"

	"github.com/chan-mai/tsugu-mcp/touki"
)

func TestProperty_Condominium(t *testing.T) {
	app := touki.Application{
		Decedent:   touki.Decedent{Name: "甲"},
		Applicants: []touki.Applicant{{Name: "乙"}},
		Properties: []touki.Property{{
			Kind:         touki.Condominium,
			Number:       "1234567890123",
			Location:     "東京都千代田区一番町23番地",
			BuildingName: "一番町マンション",
			HouseNumber:  "一番町23番の301",
			UnitName:     "301号",
			BuildingType: "居宅",
			Structure:    "鉄骨造1階建",
			FloorArea:    "3階部分 60.12平方メートル",
			LandRights: []touki.LandRight{{
				Symbol: "1", LocationLot: "東京都千代田区一番町23番", Category: "宅地",
				Area: "500.00", RightType: "所有権", RightShare: "1000分の35",
			}},
		}},
	}
	var texts []string
	for _, p := range Build(app, DefaultStyle()) {
		for _, l := range p.Labels {
			texts = append(texts, l.Text)
		}
	}
	joined := strings.Join(texts, "\n")
	for _, want := range []string{"一棟の建物の表示", "専有部分の建物の表示", "敷地権の表示", "一番町マンション", "301号", "1000分の35", "500.00平方メートル"} {
		if !strings.Contains(joined, want) {
			t.Errorf("condominium output missing %q", want)
		}
	}
}
