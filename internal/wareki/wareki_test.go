package wareki

import "testing"

func TestFormat(t *testing.T) {
	cases := []struct {
		name    string
		y, m, d int
		style   Style
		want    string
	}{
		{"昭和末日", 1989, 1, 7, Wareki, "昭和64年1月7日"},
		{"平成改元日", 1989, 1, 8, Wareki, "平成元年1月8日"},
		{"平成末日", 2019, 4, 30, Wareki, "平成31年4月30日"},
		{"令和改元日", 2019, 5, 1, Wareki, "令和元年5月1日"},
		{"令和通常年", 2025, 3, 15, Wareki, "令和7年3月15日"},
		{"併記", 2025, 3, 15, WarekiWithSeireki, "令和7年(2025年)3月15日"},
		{"西暦のみ", 2025, 3, 15, Seireki, "2025年3月15日"},
		{"大正改元日", 1912, 7, 30, Wareki, "大正元年7月30日"},
		{"明治より前は西暦", 1850, 1, 1, Wareki, "1850年1月1日"},
		{"未設定", 0, 0, 0, Wareki, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Format(c.y, c.m, c.d, c.style); got != c.want {
				t.Errorf("Format(%d,%d,%d) = %q, want %q", c.y, c.m, c.d, got, c.want)
			}
		})
	}
}
