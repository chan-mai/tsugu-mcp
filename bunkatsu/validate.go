package bunkatsu

import (
	"errors"
	"fmt"
)

// Agreementの意味的整合性を検証し、問題を連結して返却
func (a Agreement) Validate() error {
	var errs []error

	if a.Decedent.Name == "" {
		errs = append(errs, errors.New("被相続人: name is required"))
	}
	if a.Decedent.DeathDate.IsZero() {
		errs = append(errs, errors.New("被相続人: death date is required"))
	} else if !a.Decedent.DeathDate.Valid() {
		errs = append(errs, fmt.Errorf("被相続人: invalid death date (%s)", a.Decedent.DeathDate))
	}

	if len(a.Heirs) == 0 {
		errs = append(errs, errors.New("共同相続人: at least one required"))
	}
	for i, h := range a.Heirs {
		if h.Name == "" {
			errs = append(errs, fmt.Errorf("共同相続人[%d]: name is required", i))
		}
	}

	if len(a.Allocations) == 0 {
		errs = append(errs, errors.New("取得の対応: at least one allocation required"))
	}
	for i, al := range a.Allocations {
		if len(al.Acquirers) == 0 {
			errs = append(errs, fmt.Errorf("取得の対応[%d]: at least one acquirer required", i))
		}
		for j, ac := range al.Acquirers {
			if ac.Name == "" {
				errs = append(errs, fmt.Errorf("取得の対応[%d].取得者[%d]: name is required", i, j))
			}
		}
		if len(al.Properties) == 0 && len(al.Items) == 0 {
			errs = append(errs, fmt.Errorf("取得の対応[%d]: at least one property or item required", i))
		}
	}

	if !a.AgreedDate.IsZero() && !a.AgreedDate.Valid() {
		errs = append(errs, fmt.Errorf("協議成立日: invalid date (%s)", a.AgreedDate))
	}

	return errors.Join(errs...)
}

// 作成通数(指定なしは相続人数)
func (a Agreement) CopyCount() int {
	if a.Copies > 0 {
		return a.Copies
	}
	if n := len(a.Heirs); n > 0 {
		return n
	}
	return 1
}
