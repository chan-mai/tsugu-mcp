package wareki

import "testing"

func TestEraYear(t *testing.T) {
	cases := []struct {
		y, m, d int
		era     string
		wy      int
	}{
		{2026, 1, 1, "令和", 8},
		{2019, 5, 1, "令和", 1},
		{2024, 6, 15, "令和", 6},
		{1989, 1, 8, "平成", 1},
		{1985, 4, 1, "昭和", 60},
	}
	for _, c := range cases {
		era, wy := EraYear(c.y, c.m, c.d)
		if era != c.era || wy != c.wy {
			t.Errorf("EraYear(%d,%d,%d) = %q,%d, want %q,%d", c.y, c.m, c.d, era, wy, c.era, c.wy)
		}
	}
}
