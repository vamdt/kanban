package crawl

import (
	"testing"
	"time"
)

type test_cent_pair struct {
	str string
	exp int
}

var cent_tests = []test_cent_pair{
	{"10.00", 1000},
	{"10.01", 1001},
	{"-10.00", -1000},
	{"-10.01", -1001},
	{"0.01", 1},
	{"-0.01", -1},
	{"1", 100},
	{"", 0},
}

func TestParseCent(t *testing.T) {
	for _, pair := range cent_tests {
		v := ParseCent(pair.str)
		if v != pair.exp {
			t.Error(
				"For", pair.str,
				"expected", pair.exp,
				"got", v,
			)
		}
	}
}

func TestMonthend(t *testing.T) {
	type date_pair struct {
		f, exp string
	}
	tests := []date_pair{
		date_pair{f: "2000-01-01", exp: "2000-02-01"},
		date_pair{f: "2000-01-31", exp: "2000-02-01"},
		date_pair{f: "2000-02-01", exp: "2000-03-01"},
		date_pair{f: "2000-02-29", exp: "2000-03-01"},
		date_pair{f: "2000-04-01", exp: "2000-05-01"},
	}
	for _, td := range tests {
		d1, _ := time.Parse("2006-01-02", td.f)
		d1e, _ := time.Parse("2006-01-02", td.exp)
		d2 := Monthend(d1)
		if !d2.Equal(d1e) {
			t.Error(
				"For", d1,
				"expected", d1e,
				"got", d2,
			)
		}
	}
}

func TestWeekend(t *testing.T) {
	type date_pair struct {
		f, exp string
	}
	tests := []date_pair{
		date_pair{f: "2000-01-01", exp: "2000-01-01"},
		date_pair{f: "2000-02-01", exp: "2000-02-05"},
		date_pair{f: "2001-02-01", exp: "2001-02-03"},
		date_pair{f: "2000-04-01", exp: "2000-04-01"},
	}
	for _, td := range tests {
		d1, _ := time.Parse("2006-01-02", td.f)
		d1e, _ := time.Parse("2006-01-02", td.exp)
		d2 := Weekend(d1)
		if !d2.Equal(d1e) {
			t.Error(
				"For", d1,
				"expected", d1e,
				"got", d2,
			)
		}
	}
}

func TestMinuteend(t *testing.T) {
	type date_pair struct {
		f, exp string
	}
	tests := []date_pair{
		date_pair{f: "01:01:00", exp: "01:02:00"},
		date_pair{f: "01:01:01", exp: "01:02:00"},
		date_pair{f: "01:01:59", exp: "01:02:00"},
		date_pair{f: "01:01:29", exp: "01:02:00"},
		date_pair{f: "01:01:31", exp: "01:02:00"},
		date_pair{f: "01:02:01", exp: "01:03:00"},
		date_pair{f: "01:02:29", exp: "01:03:00"},
		date_pair{f: "01:02:30", exp: "01:03:00"},
		date_pair{f: "01:02:31", exp: "01:03:00"},
	}
	for _, td := range tests {
		d1, _ := time.Parse("15:04:05", td.f)
		d1e, _ := time.Parse("15:04:05", td.exp)
		d2 := Minuteend(d1)
		if !d2.Equal(d1e) {
			t.Error(
				"For", d1,
				"expected", d1e,
				"got", d2,
			)
		}
	}
}

func TestMinute5end(t *testing.T) {
	type date_pair struct {
		f, exp string
	}
	tests := []date_pair{
		date_pair{f: "01:00:00", exp: "01:05:00"},
		date_pair{f: "01:00:01", exp: "01:05:00"},
		date_pair{f: "01:01:59", exp: "01:05:00"},
		date_pair{f: "01:04:29", exp: "01:05:00"},
		date_pair{f: "01:04:59", exp: "01:05:00"},
		date_pair{f: "01:05:00", exp: "01:10:00"},
		date_pair{f: "01:05:01", exp: "01:10:00"},
		date_pair{f: "01:05:29", exp: "01:10:00"},
		date_pair{f: "01:05:30", exp: "01:10:00"},
		date_pair{f: "01:05:31", exp: "01:10:00"},
	}
	for _, td := range tests {
		d1, _ := time.Parse("15:04:05", td.f)
		d1e, _ := time.Parse("15:04:05", td.exp)
		d2 := Minute5end(d1)
		if !d2.Equal(d1e) {
			t.Error(
				"For", d1,
				"expected", d1e,
				"got", d2,
			)
		}
	}
}

func TestMinute30end(t *testing.T) {
	type date_pair struct {
		f, exp string
	}
	tests := []date_pair{
		date_pair{f: "01:00:00", exp: "01:30:00"},
		date_pair{f: "01:00:01", exp: "01:30:00"},
		date_pair{f: "01:29:59", exp: "01:30:00"},
		date_pair{f: "01:30:00", exp: "02:00:00"},
		date_pair{f: "01:30:01", exp: "02:00:00"},
		date_pair{f: "01:30:29", exp: "02:00:00"},
		date_pair{f: "01:59:59", exp: "02:00:00"},
	}
	for _, td := range tests {
		d1, _ := time.Parse("15:04:05", td.f)
		d1e, _ := time.Parse("15:04:05", td.exp)
		d2 := Minute30end(d1)
		if !d2.Equal(d1e) {
			t.Error(
				"For", d1,
				"expected", d1e,
				"got", d2,
			)
		}
	}
}
