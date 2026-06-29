package bunkatsulayout

import (
	"strings"
	"testing"

	"github.com/chan-mai/tsugu-mcp/bunkatsu"
	"github.com/chan-mai/tsugu-mcp/ymd"
)

func certInput() bunkatsu.Certificate {
	return bunkatsu.Certificate{
		Decedent: bunkatsu.Decedent{Name: "山田太郎", DeathDate: ymd.Date{Year: 2024, Month: 6, Day: 15}},
		Acquirer: "山田一郎",
		Signers:  []string{"山田花子", "山田一郎", "佐藤良子"},
		SignDate: ymd.Date{Year: 2026, Month: 1, Day: 1},
	}
}

func TestBuildCertificate_PerSigner(t *testing.T) {
	c := certInput()
	pages := BuildCertificate(c, DefaultStyle())
	if len(pages) != len(c.Signers) {
		t.Fatalf("pages = %d, want %d", len(pages), len(c.Signers))
	}
	for i, name := range c.Signers {
		signLine := ""
		for _, l := range pages[i].Labels {
			if strings.HasPrefix(l.Text, "上記相続人") {
				signLine = l.Text
			}
			if l.Y < 0 || l.Y > pages[i].Height {
				t.Errorf("label out of page %d: %q", i, l.Text)
			}
		}
		if !strings.Contains(signLine, name) {
			t.Errorf("page %d signature = %q, want signer %q", i, signLine, name)
		}
	}
}

func TestBuildCertificate_EmptySigners(t *testing.T) {
	c := certInput()
	c.Signers = nil
	if got := len(BuildCertificate(c, DefaultStyle())); got != 1 {
		t.Errorf("empty signers should yield 1 page, got %d", got)
	}
}
