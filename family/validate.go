package family

import (
	"errors"
	"fmt"
)

// Documentの意味的整合性を検証し、問題を連結して返却(日付書式等の構文検証はinputjsonの責務)
func (d Document) Validate() error {
	var errs []error

	if d.Decedent.Name == "" {
		errs = append(errs, errors.New("被相続人: name is required"))
	}
	if d.Decedent.DeathDate.IsZero() {
		errs = append(errs, errors.New("被相続人: death date is required"))
	}

	checkDate := func(role, field string, dt Date) {
		if !dt.IsZero() && !dt.Valid() {
			errs = append(errs, fmt.Errorf("%s: invalid %s (%s)", role, field, dt))
		}
	}
	checkDate("被相続人", "birth date", d.Decedent.BirthDate)
	checkDate("被相続人", "death date", d.Decedent.DeathDate)
	if !d.Decedent.BirthDate.IsZero() && !d.Decedent.DeathDate.IsZero() &&
		d.Decedent.DeathDate.Before(d.Decedent.BirthDate) {
		errs = append(errs, errors.New("被相続人: death date is before birth date"))
	}

	d.eachPerson(func(role string, p Person) {
		if p.Name == "" {
			errs = append(errs, fmt.Errorf("%s: name is required", role))
		}
		checkDate(role, "birth date", p.BirthDate)
		if p.DeathDate != nil {
			checkDate(role, "death date", *p.DeathDate)
			if !p.BirthDate.IsZero() && p.DeathDate.Before(p.BirthDate) {
				errs = append(errs, fmt.Errorf("%s: death date is before birth date", role))
			}
		}
	})

	return errors.Join(errs...)
}

// 被相続人を除く全関係者を役割ラベル付き走査
func (d Document) eachPerson(fn func(role string, p Person)) {
	if d.Spouse != nil {
		fn("配偶者", *d.Spouse)
	}
	for _, p := range d.Ascendants {
		fn("直系尊属", p)
	}
	var walk func(role string, n *Node)
	walk = func(role string, n *Node) {
		fn(role, n.Person)
		if n.Spouse != nil {
			fn(role+"の配偶者", *n.Spouse)
		}
		for _, c := range n.Descendants {
			walk(role+"の卑属", c)
		}
	}
	for _, c := range d.Children {
		walk("子", c)
	}
	for _, s := range d.Siblings {
		walk("兄弟姉妹", s)
	}
}
