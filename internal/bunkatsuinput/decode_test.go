package bunkatsuinput

import (
	"testing"

	"github.com/chan-mai/tsugu-mcp/touki"
	"github.com/chan-mai/tsugu-mcp/ymd"
)

func TestDecode_Valid(t *testing.T) {
	data := []byte(`{
		"decedent": { "name": "山田 太郎", "address": "東京都千代田区一番町1番1号", "deathDate": "2024-06-15" },
		"heirs": [
			{ "name": "山田 花子", "address": "東京都千代田区一番町1番1号" },
			{ "name": "山田 一郎", "address": "東京都新宿区西新宿2番2号" }
		],
		"allocations": [
			{ "acquirers": [ { "name": "山田 一郎" } ],
			  "properties": [ { "kind": "土地", "location": "東京都千代田区一番町1番", "lotNumber": "1番", "area": "123.45" } ] },
			{ "acquirers": [ { "name": "山田 花子" } ], "items": [ "みずほ銀行 普通預金 1234567" ] }
		],
		"agreedDate": "2024-12-10",
		"copies": 3
	}`)

	ag, err := Decode(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ag.AgreedDate != (ymd.Date{Year: 2024, Month: 12, Day: 10}) {
		t.Errorf("agreedDate=%v", ag.AgreedDate)
	}
	if len(ag.Heirs) != 2 {
		t.Errorf("heirs=%+v", ag.Heirs)
	}
	if len(ag.Allocations) != 2 {
		t.Fatalf("allocations=%+v", ag.Allocations)
	}
	if got := ag.Allocations[0].Properties[0].Kind; got != touki.Land {
		t.Errorf("property kind=%v", got)
	}
	if len(ag.Allocations[1].Items) != 1 {
		t.Errorf("items=%+v", ag.Allocations[1].Items)
	}
	if ag.Copies != 3 {
		t.Errorf("copies=%d", ag.Copies)
	}
}

func TestDecode_BadDate(t *testing.T) {
	_, err := Decode([]byte(`{"decedent":{"name":"x","deathDate":"2024-06-15"},"agreedDate":"2024-13-40"}`))
	if err == nil {
		t.Fatal("expected error for invalid agreedDate")
	}
}

func TestDecode_UnknownKind(t *testing.T) {
	_, err := Decode([]byte(`{"decedent":{"name":"x","deathDate":"2024-06-15"},"agreedDate":"2024-12-10","allocations":[{"acquirers":[{"name":"y"}],"properties":[{"kind":"船","location":"z"}]}]}`))
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
}

func TestDecodeCertificate_Valid(t *testing.T) {
	data := []byte(`{
		"decedent": { "name": "山田 太郎", "deathDate": "2024-06-15" },
		"acquirer": "山田 一郎",
		"signers": [ "山田 花子", "山田 一郎" ],
		"signDate": "2024-12-10"
	}`)

	c, err := DecodeCertificate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Decedent.Name != "山田 太郎" {
		t.Errorf("decedent name=%q", c.Decedent.Name)
	}
	if c.Acquirer != "山田 一郎" {
		t.Errorf("acquirer=%q", c.Acquirer)
	}
	if len(c.Signers) != 2 {
		t.Errorf("signers=%+v", c.Signers)
	}
	if c.SignDate != (ymd.Date{Year: 2024, Month: 12, Day: 10}) {
		t.Errorf("signDate=%v", c.SignDate)
	}
}

func TestDecodeCertificate_BadDate(t *testing.T) {
	_, err := DecodeCertificate([]byte(`{"decedent":{"name":"x","deathDate":"2024-06-15"},"acquirer":"y","signDate":"bad"}`))
	if err == nil {
		t.Fatal("expected error for invalid signDate")
	}
}
