package inputjson

import (
	"testing"

	"tsugu-mcp/family"
)

func TestDecode_Valid(t *testing.T) {
	data := []byte(`{
		"decedent": {"name": "甲", "deathDate": "2025-03-15", "birthDate": "1950-01-02"},
		"spouse": {"name": "乙", "relationship": "妻", "outcome": "inherit"},
		"children": [
			{"name": "丙", "relationship": "長男", "outcome": "相続",
			 "descendants": [{"name": "丁", "relationship": "孫", "outcome": "by_representation"}]}
		],
		"preparedAt": "2026-06-26"
	}`)

	doc, err := Decode(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc.Decedent.DeathDate != (family.Date{Year: 2025, Month: 3, Day: 15}) {
		t.Errorf("deathDate = %v", doc.Decedent.DeathDate)
	}
	if doc.Spouse == nil || doc.Spouse.Outcome != family.OutcomeInherit {
		t.Errorf("spouse outcome 解釈に失敗: %+v", doc.Spouse)
	}
	if len(doc.Children) != 1 || len(doc.Children[0].Descendants) != 1 {
		t.Fatalf("子孫ツリーの構築に失敗: %+v", doc.Children)
	}
	if doc.Children[0].Outcome != family.OutcomeInherit { // "相続" の別名
		t.Errorf("和文 outcome 解釈に失敗: %v", doc.Children[0].Outcome)
	}
}

func TestDecode_BadDate(t *testing.T) {
	_, err := Decode([]byte(`{"decedent": {"name": "甲", "deathDate": "2025-13-40"}}`))
	if err == nil {
		t.Fatal("不正な日付でエラーになるべき")
	}
}

func TestDecode_UnknownOutcome(t *testing.T) {
	_, err := Decode([]byte(`{"decedent": {"name": "甲", "deathDate": "2025-01-01"}, "spouse": {"name": "乙", "outcome": "???"}}`))
	if err == nil {
		t.Fatal("不明な outcome でエラーになるべき")
	}
}

func TestDecode_EmptyDeathDateIsAlive(t *testing.T) {
	doc, err := Decode([]byte(`{"decedent": {"name": "甲", "deathDate": "2025-01-01"}, "children": [{"name": "丙", "relationship": "長男"}]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc.Children[0].DeathDate != nil {
		t.Errorf("死亡日未指定は存命(nil)であるべき")
	}
}
