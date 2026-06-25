// Package inputjson は外部JSONとfamilyモデルの境界
// 日付書式・列挙語彙の解釈を隔離しfamily層を外部表現から保護
package inputjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"tsugu-mcp/family"
)

type document struct {
	Decedent   decedent `json:"decedent"`
	Spouse     *person  `json:"spouse"`
	Ascendants []person `json:"ascendants"`
	Children   []node   `json:"children"`
	Siblings   []node   `json:"siblings"`
	Preparer   preparer `json:"preparer"`
	PreparedAt string   `json:"preparedAt"`
}

type preparer struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type decedent struct {
	Name               string `json:"name"`
	RegisteredDomicile string `json:"registeredDomicile"`
	LastAddress        string `json:"lastAddress"`
	RegistryAddress    string `json:"registryAddress"`
	BirthDate          string `json:"birthDate"`
	DeathDate          string `json:"deathDate"`
}

type person struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Address      string `json:"address"`
	BirthDate    string `json:"birthDate"`
	DeathDate    string `json:"deathDate"`
	Outcome      string `json:"outcome"`
	Applicant    bool   `json:"applicant"`
}

type node struct {
	person
	Spouse      *person `json:"spouse"`
	Descendants []node  `json:"descendants"`
}

// JSONをfamily.Documentへ変換、書式エラーは項目名付きで連結返却
// 意味的検証はfamily.Document.Validateの責務
func Decode(data []byte) (family.Document, error) {
	var d document
	if err := json.Unmarshal(data, &d); err != nil {
		return family.Document{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var errs []error
	birth := parseDateInto("decedent.birthDate", d.Decedent.BirthDate, &errs)
	death := parseDateInto("decedent.deathDate", d.Decedent.DeathDate, &errs)
	prepared := parseDateInto("preparedAt", d.PreparedAt, &errs)

	doc := family.Document{
		Decedent: family.Decedent{
			Name:               d.Decedent.Name,
			RegisteredDomicile: d.Decedent.RegisteredDomicile,
			LastAddress:        d.Decedent.LastAddress,
			RegistryAddress:    d.Decedent.RegistryAddress,
			BirthDate:          birth,
			DeathDate:          death,
		},
		Preparer:   family.Preparer{Address: d.Preparer.Address, Name: d.Preparer.Name},
		PreparedAt: prepared,
	}
	if d.Spouse != nil {
		p := buildPerson("spouse", *d.Spouse, &errs)
		doc.Spouse = &p
	}
	for i, a := range d.Ascendants {
		doc.Ascendants = append(doc.Ascendants, buildPerson(fmt.Sprintf("ascendants[%d]", i), a, &errs))
	}
	for i, c := range d.Children {
		doc.Children = append(doc.Children, buildNode(fmt.Sprintf("children[%d]", i), c, &errs))
	}
	for i, s := range d.Siblings {
		doc.Siblings = append(doc.Siblings, buildNode(fmt.Sprintf("siblings[%d]", i), s, &errs))
	}

	return doc, errors.Join(errs...)
}

func buildNode(field string, n node, errs *[]error) *family.Node {
	out := &family.Node{Person: buildPerson(field, n.person, errs)}
	if n.Spouse != nil {
		sp := buildPerson(field+".spouse", *n.Spouse, errs)
		out.Spouse = &sp
	}
	for i, c := range n.Descendants {
		out.Descendants = append(out.Descendants, buildNode(fmt.Sprintf("%s.descendants[%d]", field, i), c, errs))
	}
	return out
}

func buildPerson(field string, p person, errs *[]error) family.Person {
	out := family.Person{
		Name:         p.Name,
		Relationship: p.Relationship,
		Address:      p.Address,
		BirthDate:    parseDateInto(field+".birthDate", p.BirthDate, errs),
		Outcome:      parseOutcome(field+".outcome", p.Outcome, errs),
		Applicant:    p.Applicant,
	}
	if strings.TrimSpace(p.DeathDate) != "" {
		death := parseDateInto(field+".deathDate", p.DeathDate, errs)
		out.DeathDate = &death
	}
	return out
}

func parseDateInto(field, s string, errs *[]error) family.Date {
	s = strings.TrimSpace(s)
	if s == "" {
		return family.Date{}
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s: date must be in YYYY-MM-DD format: %q", field, s))
		return family.Date{}
	}
	return family.Date{Year: t.Year(), Month: int(t.Month()), Day: t.Day()}
}

func parseOutcome(field, s string, errs *[]error) family.Outcome {
	switch strings.TrimSpace(s) {
	case "", "none":
		return family.OutcomeNone
	case "inherit", "相続":
		return family.OutcomeInherit
	case "renounce", "相続放棄", "放棄":
		return family.OutcomeRenounce
	case "division", "分割":
		return family.OutcomeDivision
	case "by_representation", "代襲", "代襲相続":
		return family.OutcomeByRepresentation
	default:
		*errs = append(*errs, fmt.Errorf("%s: unknown outcome: %q", field, s))
		return family.OutcomeNone
	}
}
