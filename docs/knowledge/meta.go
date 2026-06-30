package knowledge

import (
	"github.com/chan-mai/tsugu-mcp/internal/wareki"
	"github.com/chan-mai/tsugu-mcp/ymd"
)

// 知識ベースの調査基準日(西暦YYYY-MM-DD), 法改正反映時に手動更新
const AsOf = "2026-06-27"

var AsOfJP = toWareki(AsOf)

func toWareki(s string) string {
	d, err := ymd.Parse(s)
	if err != nil {
		return s
	}
	return wareki.Format(d.Year, d.Month, d.Day, wareki.Wareki)
}
