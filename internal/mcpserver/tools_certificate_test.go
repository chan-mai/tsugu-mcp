package mcpserver

import (
	"context"
	"testing"
)

func TestHandleCertificate(t *testing.T) {
	res, out, err := handleCertificate(context.Background(), nil, certToolInput{Document: certDoc{
		Decedent: certDecedent{Name: "山田太郎", DeathDate: "2024-06-15"},
		Acquirer: "山田一郎",
		Signers:  []string{"山田花子", "山田一郎", "佐藤良子"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error: %+v", res.Content)
	}
	if out.Path == "" {
		t.Error("no output path")
	}
}

func TestHandleCertificate_Invalid(t *testing.T) {
	res, _, _ := handleCertificate(context.Background(), nil, certToolInput{Document: certDoc{
		Decedent: certDecedent{Name: "山田太郎", DeathDate: "2024-06-15"},
	}})
	if !res.IsError {
		t.Error("missing acquirer/signers should error")
	}
}
