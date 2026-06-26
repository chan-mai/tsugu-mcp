// Package reginputは登記申請書JSONとtoukiモデルの境界
// 日付書式・種別語彙の解釈を隔離
package reginput

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"tsugu-mcp/touki"
	"tsugu-mcp/ymd"
)

type document struct {
	Causes          []cause     `json:"causes"`
	Decedent        decedent    `json:"decedent"`
	Applicants      []applicant `json:"applicants"`
	Attachments     []string    `json:"attachments"`
	DeclineIDInfo   bool        `json:"declineIdInfo"`
	ApplicationDate string      `json:"applicationDate"`
	Registry        string      `json:"registry"`
	TaxValue        string      `json:"taxValue"`
	RegistrationTax string      `json:"registrationTax"`
	Properties      []property  `json:"properties"`
}

type cause struct {
	Date string `json:"date"`
	Text string `json:"text"`
}

type decedent struct {
	Name string `json:"name"`
}

type applicant struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Share     string `json:"share"`
	NameKana  string `json:"nameKana"`
	BirthDate string `json:"birthDate"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contact   bool   `json:"contact"`
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

// JSONをtouki.Applicationへ変換、書式エラーは項目名付きで連結返却(意味的検証はtouki.Validateの責務)
func Decode(data []byte) (touki.Application, error) {
	var d document
	if err := json.Unmarshal(data, &d); err != nil {
		return touki.Application{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var errs []error
	app := touki.Application{
		Decedent:        touki.Decedent{Name: d.Decedent.Name},
		Attachments:     d.Attachments,
		DeclineIDInfo:   d.DeclineIDInfo,
		ApplicationDate: parseDate("applicationDate", d.ApplicationDate, &errs),
		Registry:        d.Registry,
		TaxValue:        d.TaxValue,
		RegistrationTax: d.RegistrationTax,
	}
	for i, c := range d.Causes {
		app.Causes = append(app.Causes, touki.Cause{
			Date: parseDate(fmt.Sprintf("causes[%d].date", i), c.Date, &errs),
			Text: c.Text,
		})
	}
	for i, a := range d.Applicants {
		app.Applicants = append(app.Applicants, touki.Applicant{
			Name:      a.Name,
			Address:   a.Address,
			Share:     a.Share,
			NameKana:  a.NameKana,
			BirthDate: parseDate(fmt.Sprintf("applicants[%d].birthDate", i), a.BirthDate, &errs),
			Email:     a.Email,
			Phone:     a.Phone,
			Contact:   a.Contact,
		})
	}
	for i, p := range d.Properties {
		app.Properties = append(app.Properties, touki.Property{
			Kind:         parseKind(fmt.Sprintf("properties[%d].kind", i), p.Kind, &errs),
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

	return app, errors.Join(errs...)
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
