package bunkatsu

import (
	"strings"
	"testing"

	"github.com/chan-mai/tsugu-mcp/ymd"
)

func okCertificate() Certificate {
	return Certificate{
		Decedent: Decedent{Name: "甲", DeathDate: ymd.Date{Year: 2024, Month: 6, Day: 15}},
		Acquirer: "乙",
		Signers:  []string{"乙", "丙"},
	}
}

func TestCertificate_Validate_OK(t *testing.T) {
	if err := okCertificate().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCertificate_Validate_Errors(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*Certificate)
		want string
	}{
		{"死亡日なし", func(c *Certificate) { c.Decedent.DeathDate = ymd.Date{} }, "death date is required"},
		{"取得者なし", func(c *Certificate) { c.Acquirer = "" }, "取得者"},
		{"署名者なし", func(c *Certificate) { c.Signers = nil }, "at least one signer"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cert := okCertificate()
			c.mut(&cert)
			err := cert.Validate()
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("err = %v, want contains %q", err, c.want)
			}
		})
	}
}
