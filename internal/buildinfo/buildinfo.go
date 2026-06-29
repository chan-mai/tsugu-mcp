// Package buildinfoはバージョン情報を一元管理する
package buildinfo

import (
	"runtime/debug"
	"time"
)

// fallbackはVCS情報が取れないときの基準バージョン
const fallback = "0.0.0"

var jst = time.FixedZone("JST", 9*60*60)

func Version() string {
	if d := commitDate(); d != "" {
		return d
	}
	return fallback
}

func String() string {
	if rev := revision(); rev != "" {
		return Version() + "+" + rev
	}
	return Version()
}

// commitDateはビルド時に埋め込まれたコミット日時をYYYY.MM.DDで返す
func commitDate() string {
	for _, s := range settings() {
		if s.Key == "vcs.time" {
			if t, err := time.Parse(time.RFC3339, s.Value); err == nil {
				return t.In(jst).Format("2006.01.02")
			}
		}
	}
	return ""
}

// revisionは短縮コミットハッシュ
func revision() string {
	for _, s := range settings() {
		if s.Key == "vcs.revision" {
			if len(s.Value) > 12 {
				return s.Value[:12]
			}
			return s.Value
		}
	}
	return ""
}

func settings() []debug.BuildSetting {
	if bi, ok := debug.ReadBuildInfo(); ok {
		return bi.Settings
	}
	return nil
}
