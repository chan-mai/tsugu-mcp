package registration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"tsugu-mcp/registration"
)

// testdataの登記申請書サンプルが公開API経由でPDF化されるか確認(end-to-end)
func TestGenerateFromJSON_Samples(t *testing.T) {
	samples, err := filepath.Glob("../testdata/touki_*.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) == 0 {
		t.Fatal("touki サンプルが見つかりません")
	}
	for _, path := range samples {
		t.Run(filepath.Base(path), func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			pdf, err := registration.GenerateFromJSON(data, registration.DefaultOptions())
			if err != nil {
				t.Fatalf("生成に失敗: %v", err)
			}
			if !bytes.HasPrefix(pdf, []byte("%PDF")) {
				t.Errorf("PDFになっていません")
			}
		})
	}
}

func TestGenerateFromJSON_Invalid(t *testing.T) {
	// 申請人・不動産なし → Validate でエラー
	_, err := registration.GenerateFromJSON([]byte(`{"registry":"東京法務局"}`), registration.DefaultOptions())
	if err == nil {
		t.Fatal("検証エラーが返るべき")
	}
}
