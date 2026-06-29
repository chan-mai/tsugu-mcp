package bunkatsuinput

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/chan-mai/tsugu-mcp/bunkatsu"
)

type certDocument struct {
	Decedent certDecedent `json:"decedent"`
	Acquirer string       `json:"acquirer"`
	Signers  []string     `json:"signers"`
	SignDate string       `json:"signDate"`
}

type certDecedent struct {
	Name      string `json:"name"`
	DeathDate string `json:"deathDate"`
}

// JSONをbunkatsu.Certificateへ変換、書式エラーは項目名付きで連結返却
func DecodeCertificate(data []byte) (bunkatsu.Certificate, error) {
	var d certDocument
	if err := json.Unmarshal(data, &d); err != nil {
		return bunkatsu.Certificate{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var errs []error
	c := bunkatsu.Certificate{
		Decedent: bunkatsu.Decedent{
			Name:      d.Decedent.Name,
			DeathDate: parseDate("decedent.deathDate", d.Decedent.DeathDate, &errs),
		},
		Acquirer: d.Acquirer,
		Signers:  d.Signers,
		SignDate: parseDate("signDate", d.SignDate, &errs),
	}
	return c, errors.Join(errs...)
}
