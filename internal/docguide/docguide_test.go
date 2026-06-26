package docguide

import (
	"strings"
	"testing"
)

func flat(r Result) string {
	var sb strings.Builder
	for _, c := range r.Categories {
		sb.WriteString(c.Name)
		for _, it := range c.Items {
			sb.WriteString(it)
		}
	}
	for _, n := range r.Notes {
		sb.WriteString(n)
	}
	return sb.String()
}

func TestRequiredDocuments_Siblings(t *testing.T) {
	s := flat(RequiredDocuments(Input{Method: "legal", HeirPattern: "siblings", ApplicantAtWindow: true}))
	if !strings.Contains(s, "父・母") {
		t.Error("siblings pattern must include parents' continuous register")
	}
	if !strings.Contains(s, "広域交付") {
		t.Error("siblings at window must include 広域交付 caveat")
	}
}

func TestRequiredDocuments_AgreementHasSeal(t *testing.T) {
	if !strings.Contains(flat(RequiredDocuments(Input{Method: "agreement", HeirPattern: "children"})), "印鑑証明書") {
		t.Error("agreement must include 印鑑証明書")
	}
}

func TestRequiredDocuments_LegalNoSeal(t *testing.T) {
	if strings.Contains(flat(RequiredDocuments(Input{Method: "legal", HeirPattern: "children"})), "印鑑証明書") {
		t.Error("legal must not include 印鑑証明書")
	}
}

func TestRequiredDocuments_Substitution(t *testing.T) {
	if !strings.Contains(flat(RequiredDocuments(Input{Method: "legal", HeirPattern: "substitution"})), "先に死亡した子") {
		t.Error("substitution must include deceased child's register")
	}
}

func TestRequiredDocuments_Identity(t *testing.T) {
	r := RequiredDocuments(Input{Method: "legal", HeirPattern: "children", RegistryAddressDiffersFromHonseki: true})
	found := false
	for _, c := range r.Categories {
		if c.Name == "同一性証明" {
			found = true
		}
	}
	if !found {
		t.Error("differing registry address must add 同一性証明 category")
	}
}
