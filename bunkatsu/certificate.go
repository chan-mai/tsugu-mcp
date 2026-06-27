package bunkatsu

import (
	"errors"
	"fmt"

	"tsugu-mcp/ymd"
)

// 遺産分割協議証明書(相続人ごとに各自が取得を証明する個別書面型)
type Certificate struct {
	Decedent Decedent // 被相続人(Name+DeathDate、Addressは未使用)
	Acquirer string   // 不動産を取得する相続人の氏名
	Signers  []string // 共同相続人の氏名(各自1ページに氏名を印字)
	SignDate ymd.Date // 署名日(任意、年のみ印字で月日は空欄)
}

// Certificateの意味的整合性を検証
func (c Certificate) Validate() error {
	var errs []error

	if c.Decedent.Name == "" {
		errs = append(errs, errors.New("被相続人: name is required"))
	}
	if c.Decedent.DeathDate.IsZero() {
		errs = append(errs, errors.New("被相続人: death date is required"))
	} else if !c.Decedent.DeathDate.Valid() {
		errs = append(errs, fmt.Errorf("被相続人: invalid death date (%s)", c.Decedent.DeathDate))
	}
	if c.Acquirer == "" {
		errs = append(errs, errors.New("取得者: name is required"))
	}
	if len(c.Signers) == 0 {
		errs = append(errs, errors.New("共同相続人: at least one signer required"))
	}
	for i, s := range c.Signers {
		if s == "" {
			errs = append(errs, fmt.Errorf("共同相続人[%d]: name is required", i))
		}
	}
	if !c.SignDate.IsZero() && !c.SignDate.Valid() {
		errs = append(errs, fmt.Errorf("署名日: invalid date (%s)", c.SignDate))
	}

	return errors.Join(errs...)
}
