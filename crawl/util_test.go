package crawl

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func TestIsTradeTime(t *testing.T) {
	type date_pair struct {
		t  string
		is bool
	}
	tests := []date_pair{
		date_pair{t: "2000-01-01 01:00:00", is: false},
		date_pair{t: "2000-05-15 09:25:00", is: true},
		date_pair{t: "2000-05-15 11:30:00", is: true},
		date_pair{t: "2000-05-15 11:35:00", is: true},
		date_pair{t: "2000-05-15 11:36:00", is: false},
		date_pair{t: "2000-05-15 12:35:01", is: false},
		date_pair{t: "2000-05-15 13:00:00", is: true},
		date_pair{t: "2000-05-15 15:00:00", is: true},
		date_pair{t: "2000-05-15 15:05:00", is: true},
		date_pair{t: "2000-05-15 15:06:00", is: false},
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for _, td := range tests {
		time, _ := time.ParseInLocation("2006-01-02 15:04:05", td.t, loc)
		if td.is != IsTradeTime(time) {
			t.Error(
				"For", time,
				"expected", td.is,
				"got", !td.is,
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

func TestTimeForMinuteend(t *testing.T) {
	f := "01:01:00"
	exp := "01:01:00"
	tf, _ := time.Parse("15:04:05", f)
	texp, _ := time.Parse("15:04:05", exp)
	for i := 0; i < 10; i++ {
		tf = tf.Truncate(time.Minute)
		if !tf.Equal(texp) {
			t.Error(
				"For", f,
				"expected", exp,
				"got", tf,
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

type desc_text_pair struct {
	Desc string
	Text string
	File string
}

func load_test_desc_text_files(pattern string) []desc_text_pair {
	if len(pattern) < 1 {
		return nil
	}
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}

	var sets = []desc_text_pair{}
	for _, f := range files {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			continue
		}
		content = bytes.TrimSpace(content)
		infos := bytes.SplitN(content, []byte("\n\n"), 3)
		if infos == nil || len(infos) < 2 {
			continue
		}
		t := desc_text_pair{Desc: string(infos[0]), Text: string(infos[1])}
		t.File = f
		sets = append(sets, t)
	}

	if len(sets) < 1 {
		return nil
	}
	return sets
}
