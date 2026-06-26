package agreement_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"tsugu-mcp/agreement"
)

// testdataの遺産分割協議書サンプルが公開API経由でPDF化されるか確認(end-to-end)
func TestGenerateFromJSON_Samples(t *testing.T) {
	samples, err := filepath.Glob("../testdata/bunkatsu_*.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) == 0 {
		t.Fatal("bunkatsu サンプルが見つかりません")
	}
	for _, path := range samples {
		t.Run(filepath.Base(path), func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			pdf, err := agreement.GenerateFromJSON(data, agreement.DefaultOptions())
			if err != nil {
				t.Fatalf("生成に失敗: %v", err)
			}
			if !bytes.HasPrefix(pdf, []byte("%PDF")) {
				t.Errorf("not a PDF")
			}
		})
	}
}

func TestGenerateFromJSON_Invalid(t *testing.T) {
	_, err := agreement.GenerateFromJSON([]byte(`{"decedent":{"name":"甲"}}`), agreement.DefaultOptions())
	if err == nil {
		t.Fatal("検証エラーが返るべき")
	}
}
