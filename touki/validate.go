package touki

import (
	"errors"
	"fmt"
)

// Applicationの意味的整合性を検証し、問題を連結して返却(日付書式等の構文検証はreginputの責務)
func (a Application) Validate() error {
	var errs []error

	if len(a.Applicants) == 0 {
		errs = append(errs, errors.New("申請人: at least one required"))
	}
	for i, ap := range a.Applicants {
		if ap.Name == "" {
			errs = append(errs, fmt.Errorf("申請人[%d]: name is required", i))
		}
		if !ap.BirthDate.IsZero() && !ap.BirthDate.Valid() {
			errs = append(errs, fmt.Errorf("申請人[%d]: invalid birth date (%s)", i, ap.BirthDate))
		}
	}

	if len(a.Properties) == 0 {
		errs = append(errs, errors.New("不動産: at least one required"))
	}
	for i, p := range a.Properties {
		switch p.Kind {
		case Land:
			if p.Location == "" || p.LotNumber == "" || p.Area == "" {
				errs = append(errs, fmt.Errorf("不動産[%d](土地): location, lotNumber and area are required", i))
			}
		case Building:
			if p.Location == "" || p.HouseNumber == "" || p.FloorArea == "" {
				errs = append(errs, fmt.Errorf("不動産[%d](建物): location, houseNumber and floorArea are required", i))
			}
		}
	}

	for i, c := range a.Causes {
		if !c.Date.IsZero() && !c.Date.Valid() {
			errs = append(errs, fmt.Errorf("原因[%d]: invalid date (%s)", i, c.Date))
		}
	}
	if !a.ApplicationDate.IsZero() && !a.ApplicationDate.Valid() {
		errs = append(errs, fmt.Errorf("申請日: invalid date (%s)", a.ApplicationDate))
	}

	return errors.Join(errs...)
}
