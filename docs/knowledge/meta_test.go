package knowledge

import (
	"strings"
	"testing"
)

// 調査基準日のGo定数とmarkdown記載の乖離を防ぐ
func TestAsOfMatchesIndex(t *testing.T) {
	data, err := FS.ReadFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), AsOf) {
		t.Fatalf("docs/knowledge/README.md should contain as-of date %s", AsOf)
	}
}
