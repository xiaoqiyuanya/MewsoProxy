package order

import (
	"testing"

	"mewsoproxy/server/model"
)

func intp(v int) *int { return &v }

func TestPeriodPrice(t *testing.T) {
	p := &model.Plan{
		MonthPrice:    intp(1000),
		QuarterPrice:  intp(2700),
		YearPrice:     intp(9000),
		OnetimePrice:  intp(500),
	}
	cases := []struct {
		name     string
		period   string
		wantVal  int
		wantDur  int64
		wantOK   bool
	}{
		{"month", "month_price", 1000, 30 * 86400, true},
		{"quarter", "quarter_price", 2700, 90 * 86400, true},
		{"year", "year_price", 9000, 365 * 86400, true},
		{"onetime fallback", "onetime_price", 500, 30 * 86400, true},
		{"invalid", "week_price", 0, 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val, dur, ok := periodPrice(p, c.period)
			if ok != c.wantOK {
				t.Fatalf("ok = %v, want %v", ok, c.wantOK)
			}
			if ok && (val != c.wantVal || dur != c.wantDur) {
				t.Fatalf("got val=%d dur=%d, want val=%d dur=%d", val, dur, c.wantVal, c.wantDur)
			}
		})
	}
}

func TestMaxPeriodPrice(t *testing.T) {
	p := &model.Plan{
		MonthPrice:    intp(1000),
		YearPrice:     intp(9000),
		QuarterPrice:  intp(2700),
		ThreeYearPrice: intp(18000),
	}
	if got := maxPeriodPrice(p); got != 18000 {
		t.Fatalf("maxPeriodPrice = %d, want 18000", got)
	}
}
