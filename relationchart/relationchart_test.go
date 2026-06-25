package relationchart_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"tsugu-mcp/relationchart"
)

// testdata各サンプルが公開API経由でPDF化されるか確認(end-to-end)
func TestGenerateFromJSON_Samples(t *testing.T) {
	samples, err := filepath.Glob("../testdata/sample_*.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) == 0 {
		t.Fatal("サンプルJSONが見つかりません")
	}
	for _, path := range samples {
		t.Run(filepath.Base(path), func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			pdf, err := relationchart.GenerateFromJSON(data, relationchart.DefaultOptions())
			if err != nil {
				t.Fatalf("生成に失敗: %v", err)
			}
			if !bytes.HasPrefix(pdf, []byte("%PDF")) {
				t.Errorf("PDFになっていません")
			}
		})
	}
}

func TestGenerateFromJSON_InvalidReturnsError(t *testing.T) {
	// 死亡日のない被相続人→Validateでエラー
	_, err := relationchart.GenerateFromJSON([]byte(`{"decedent":{"name":"甲"}}`), relationchart.DefaultOptions())
	if err == nil {
		t.Fatal("検証エラーが返るべき")
	}
}
