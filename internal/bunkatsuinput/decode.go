// Package bunkatsuinputは遺産分割協議書JSONとbunkatsuモデルの境界
package bunkatsuinput

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"tsugu-mcp/bunkatsu"
	"tsugu-mcp/touki"
	"tsugu-mcp/ymd"
)

type document struct {
	Decedent    decedent     `json:"decedent"`
	Heirs       []heir       `json:"heirs"`
	Allocations []allocation `json:"allocations"`
	AgreedDate  string       `json:"agreedDate"`
	Copies      int          `json:"copies"`
}

type decedent struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	DeathDate string `json:"deathDate"`
}

type heir struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type acquirer struct {
	Name  string `json:"name"`
	Share string `json:"share"`
}

type property struct {
	Kind         string `json:"kind"`
	Number       string `json:"number"`
	Location     string `json:"location"`
	LotNumber    string `json:"lotNumber"`
	LandCategory string `json:"landCategory"`
	Area         string `json:"area"`
	HouseNumber  string `json:"houseNumber"`
	BuildingType string `json:"buildingType"`
	Structure    string `json:"structure"`
	FloorArea    string `json:"floorArea"`
}

type allocation struct {
	Acquirers  []acquirer `json:"acquirers"`
	Properties []property `json:"properties"`
	Items      []string   `json:"items"`
}

// JSONをbunkatsu.Agreementへ変換、書式エラーは項目名付きで連結返却
func Decode(data []byte) (bunkatsu.Agreement, error) {
	var d document
	if err := json.Unmarshal(data, &d); err != nil {
		return bunkatsu.Agreement{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var errs []error
	ag := bunkatsu.Agreement{
		Decedent: bunkatsu.Decedent{
			Name:      d.Decedent.Name,
			Address:   d.Decedent.Address,
			DeathDate: parseDate("decedent.deathDate", d.Decedent.DeathDate, &errs),
		},
		AgreedDate: parseDate("agreedDate", d.AgreedDate, &errs),
		Copies:     d.Copies,
	}
	for _, h := range d.Heirs {
		ag.Heirs = append(ag.Heirs, bunkatsu.Heir{Name: h.Name, Address: h.Address})
	}
	for i, al := range d.Allocations {
		out := bunkatsu.Allocation{Items: al.Items}
		for _, ac := range al.Acquirers {
			out.Acquirers = append(out.Acquirers, bunkatsu.Acquirer{Name: ac.Name, Share: ac.Share})
		}
		for j, p := range al.Properties {
			out.Properties = append(out.Properties, touki.Property{
				Kind:         parseKind(fmt.Sprintf("allocations[%d].properties[%d].kind", i, j), p.Kind, &errs),
				Number:       p.Number,
				Location:     p.Location,
				LotNumber:    p.LotNumber,
				LandCategory: p.LandCategory,
				Area:         p.Area,
				HouseNumber:  p.HouseNumber,
				BuildingType: p.BuildingType,
				Structure:    p.Structure,
				FloorArea:    p.FloorArea,
			})
		}
		ag.Allocations = append(ag.Allocations, out)
	}

	return ag, errors.Join(errs...)
}

func parseDate(field, s string, errs *[]error) ymd.Date {
	d, err := ymd.Parse(s)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s: %w", field, err))
	}
	return d
}

func parseKind(field, s string, errs *[]error) touki.PropertyKind {
	switch strings.TrimSpace(s) {
	case "", "land", "土地":
		return touki.Land
	case "building", "建物":
		return touki.Building
	default:
		*errs = append(*errs, fmt.Errorf("%s: unknown kind: %q (land|building)", field, s))
		return touki.Land
	}
}
