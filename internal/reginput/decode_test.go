package reginput

import (
	"testing"

	"tsugu-mcp/touki"
	"tsugu-mcp/ymd"
)

func TestDecode_Valid(t *testing.T) {
	data := []byte(`{
		"causes": [ { "date": "2024-06-15", "text": "相続" } ],
		"decedent": { "name": "山田 太郎", "address": "東京都千代田区一番町1番地" },
		"applicants": [ { "name": "山田 花子", "birthDate": "1945-05-05", "share": "2分の1", "contact": true } ],
		"applicationDate": "2024-12-10",
		"properties": [
			{ "kind": "土地", "location": "東京都千代田区一番町1番", "lotNumber": "1番", "area": "123.45" },
			{ "kind": "building", "location": "東京都千代田区一番町1番地1", "houseNumber": "1番1", "floorArea": "60.00" }
		]
	}`)

	app, err := Decode(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app.ApplicationDate != (ymd.Date{Year: 2024, Month: 12, Day: 10}) {
		t.Errorf("applicationDate=%v", app.ApplicationDate)
	}
	if len(app.Applicants) != 1 || !app.Applicants[0].Contact {
		t.Errorf("applicant 解釈に失敗: %+v", app.Applicants)
	}
	if len(app.Properties) != 2 || app.Properties[0].Kind != touki.Land || app.Properties[1].Kind != touki.Building {
		t.Fatalf("property kind 解釈に失敗: %+v", app.Properties)
	}
}

func TestDecode_BadDate(t *testing.T) {
	_, err := Decode([]byte(`{"applicants":[{"name":"x"}],"properties":[],"applicationDate":"2024-13-40"}`))
	if err == nil {
		t.Fatal("不正な日付でエラーになるべき")
	}
}

func TestDecode_UnknownKind(t *testing.T) {
	_, err := Decode([]byte(`{"properties":[{"kind":"船"}]}`))
	if err == nil {
		t.Fatal("不明なkindでエラーになるべき")
	}
}
