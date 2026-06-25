package ymd

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		in      string
		want    Date
		wantErr bool
	}{
		{"1970-01-27", Date{1970, 1, 27}, false},
		{"  2024-06-15 ", Date{2024, 6, 15}, false},
		{"", Date{}, false},
		{"2025-13-40", Date{}, true},
		{"2025/01/01", Date{}, true},
		{"abc", Date{}, true},
	}
	for _, c := range cases {
		got, err := Parse(c.in)
		if (err != nil) != c.wantErr {
			t.Errorf("Parse(%q) err=%v, wantErr=%v", c.in, err, c.wantErr)
		}
		if got != c.want {
			t.Errorf("Parse(%q)=%v, want %v", c.in, got, c.want)
		}
	}
}

func TestValid(t *testing.T) {
	if (Date{2025, 2, 30}).Valid() {
		t.Error("2025-02-30 should be invalid")
	}
	if !(Date{2024, 2, 29}).Valid() {
		t.Error("2024-02-29 should be valid")
	}
}
