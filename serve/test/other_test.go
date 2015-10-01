package main

import (
	"strconv"
	"testing"
	"time"
)

func TestAtoi(t *testing.T) {
	i, _ := strconv.Atoi("-05")
	if i != -5 {
		t.Error(
			"For", "-05",
			"expected", -5,
			"got", i,
		)
	}
}

func TestByteString(t *testing.T) {
	var b []byte
	s := string(b)
	if s != "" {
		t.Error(
			"For", "",
			"expected", "",
			"got", s,
		)
	}
}

func TestAddDate(t *testing.T) {
	type date_pair struct {
		f, exp string
	}
	tests := []date_pair{
		date_pair{f: "2000-01-01", exp: "2000-01-31"},
		date_pair{f: "2000-02-01", exp: "2000-02-29"},
		date_pair{f: "2001-02-01", exp: "2001-02-28"},
		date_pair{f: "2000-04-01", exp: "2000-04-30"},
	}
	for _, td := range tests {
		d1, _ := time.Parse("2006-01-02", td.f)
		d1e, _ := time.Parse("2006-01-02", td.exp)
		_, _, d := d1.Date()
		d2 := d1.AddDate(0, 1, -d)
		if !d2.Equal(d1e) {
			t.Error(
				"For", d1,
				"expected", d1e,
				"got", d2,
			)
		}
	}
}
