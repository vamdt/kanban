package crawl

import (
	"testing"
	"time"

	. "./base"
)

func stringSlice2Tdatas(dates []string) Tdatas {
	tdatas := Tdatas{}
	fmt := "2006-01-02 15:04:05"
	tdatas.Data = make([]Tdata, len(dates))
	for i := len(tdatas.Data) - 1; i > -1; i-- {
		d, _ := time.Parse(fmt[:len(dates[i])], dates[i])
		tdatas.Data[i].Time = d
		tdatas.Data[i].Volume = 1
	}
	return tdatas
}

func TestTdatasAdd(t *testing.T) {
	dates := []string{
		"2000-01-01",
		"2000-01-02",
		"2000-01-03",
		"2000-01-04",
		"2000-01-06",
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
		cases := stringSlice2Tdatas(dates)
		old_len := len(cases.Data)

		d, _ := time.Parse("2006-01-02", tests[i].date)
		tdata := Tdata{Time: d, Volume: 1}
		index, _ := cases.Add(tdata)
		new_len := len(cases.Data)

		if tests[i].new && new_len-old_len != 1 {
			t.Error(
				"For", "case", i, tests[i].date,
				"expected", "length +1",
				"got", "newlen", new_len, "oldlen", old_len,
			)
		}

		if tests[i].index != index {
			t.Error(
				"For", "case", i,
				"expected", "insert at", tests[i].index,
				"got", "index", index,
			)
		}

		if !(tests[i].index < new_len) {
			t.Error(
				"For", "case", i,
				"expected", "index lt len",
				"got", "index >= len",
			)
		}

		if !cases.Data[tests[i].index].Time.Equal(d) {
			t.Error(
				"For", "case", i, tests[i].date,
				"expected", "Time eq [].Time",
				"got", "[].Time", cases.Data[tests[i].index].Time,
			)
		}
	}
}

func TestFirstLastdayData(t *testing.T) {
	TestDropLastdayData(t)
}

func TestDropLastdayData(t *testing.T) {
	type date_pair struct {
		date  []string
		index int
	}
	tests := []date_pair{
		{[]string{
			"2000-01-01 01:00:00",
			"2000-01-01 02:00:00",
		}, 0},
		{[]string{
			"1999-01-01 01:00:00",
			"2000-01-01 01:00:00",
			"2000-01-01 02:00:00",
		}, 1},
		{[]string{
			"1999-01-01 01:00:00",
			"2000-01-01 01:00:00",
			"2000-01-01 02:00:00",
			"2000-01-02 02:00:00",
		}, 3},
		{[]string{
			"1999-01-01 01:00:00",
			"2000-01-01 01:00:00",
			"2000-01-01 02:00:00",
			"2000-01-02 23:59:59",
		}, 3},
		{[]string{
			"1999-01-01 01:00:00",
			"2000-01-01 01:00:00",
			"2000-01-01 02:00:00",
			"2000-01-02 00:00:00",
		}, 3},
	}

	fmt := "2006-01-02 15:04:05"
	for i, l := 0, len(tests); i < l; i++ {
		cases := stringSlice2Tdatas(tests[i].date)
		index := cases.First_lastday_data()
		if index != tests[i].index {
			t.Error(
				"For", "case", i, tests[i],
				"expected", tests[i].index,
				"got", index,
			)
		}

		cases.Drop_lastday_data()
		if index != len(cases.Data) {
			t.Error(
				"For", "case", i, tests[i],
				"expected", "len(Data)==index",
				"got", "index", index, "len", len(cases.Data),
			)
		}
		if index < 1 {
			continue
		}
		ndate := cases.Data[index-1].Time.Format(fmt[:len(tests[i].date[index-1])])
		if ndate != tests[i].date[index-1] {
			t.Error(
				"For", "case", i, tests[i],
				"expected", "last Data.Time eq ", tests[i].date[index-1],
				"got", "date[", index-1, "]", ndate,
			)
		}
		t.Log("expected", tests[i].date[index-1], "got", ndate)
	}
}
