package crawl

import (
	"testing"
	"time"
)

func TestTicksAdd(t *testing.T) {
	var base_ticks_data = func() Ticks {
		dates := []string{
			"2000-01-01",
			"2000-01-02",
			"2000-01-03",
			"2000-01-04",
			"2000-01-06",
		}
		ticks := Ticks{}
		ticks.Data = make([]Tick, len(dates))
		for i, l := 0, len(ticks.Data); i < l; i++ {
			d, _ := time.Parse("2006-01-02", dates[i])
			ticks.Data[i].Time = d
		}
		return ticks
	}

	type date_pair struct {
		date  string
		index int
		new   bool
	}
	tests := []date_pair{
		{"1999-01-01", 0, true},
		{"2000-01-01", 0, false},
		{"2000-01-02", 1, false},
		{"2000-01-04", 3, false},
		{"2000-01-05", 4, true},
		{"2000-01-06", 4, false},
		{"2000-01-07", 5, true},
	}

	for i, l := 0, len(tests); i < l; i++ {
		ticks := base_ticks_data()
		old_len := len(ticks.Data)

		d, _ := time.Parse("2006-01-02", tests[i].date)
		tick := Tick{Time: d}
		ticks.Add(tick)
		new_len := len(ticks.Data)

		if tests[i].new && new_len-old_len != 1 {
			t.Error(
				"For", "case", i, tests[i].date,
				"expected", "length +1",
				"got", "newlen", new_len, "oldlen", old_len,
			)
		}

		if !(tests[i].index < new_len) {
			t.Error(
				"For", "case", i,
				"expected", "index lt len",
				"got", "index >= len",
			)
		}

		if !ticks.Data[tests[i].index].Time.Equal(d) {
			t.Error(
				"For", "case", i, tests[i].date,
				"expected", "Time eq [].Time",
				"got", "[].Time", ticks.Data[tests[i].index].Time,
			)
		}
	}
}
