package family

import (
	"fmt"
	"time"
)

// 年月日のみの暦日(時刻・タイムゾーン非対応)
type Date struct {
	Year  int
	Month int
	Day   int
}

// 未設定(ゼロ値)判定
func (d Date) IsZero() bool { return d == Date{} }

// 実在する暦日か判定(2月30日等はfalse)
func (d Date) Valid() bool {
	if d.Month < 1 || d.Month > 12 || d.Day < 1 {
		return false
	}
	t := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
	return t.Year() == d.Year && int(t.Month()) == d.Month && t.Day() == d.Day
}

// dがotherより前の暦日か判定
func (d Date) Before(other Date) bool {
	switch {
	case d.Year != other.Year:
		return d.Year < other.Year
	case d.Month != other.Month:
		return d.Month < other.Month
	default:
		return d.Day < other.Day
	}
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}
