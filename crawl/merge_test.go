package crawl

import (
	"testing"
	"time"

	. "./base"
)

func _d(s string) time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	d, _ := time.ParseInLocation("2006-01-02 15:04:05", s, loc)
	return d
}

func TestM1s2M5s(t *testing.T) {
	type Case struct {
		m1 []Tdata
		m5 []Tdata
	}

	tests := []Case{
		// 09:35 [09:30, 09:35)
		{[]Tdata{
			{Time: _d("2000-01-01 09:30:00"), Volume: 1},
			{Time: _d("2000-01-01 09:31:00"), Volume: 1},
			{Time: _d("2000-01-01 09:32:00"), Volume: 1},
			{Time: _d("2000-01-01 09:34:59"), Volume: 1},
			{Time: _d("2000-01-01 09:35:00"), Volume: 1},
		}, []Tdata{
			{Time: _d("2000-01-01 09:35:00"), Volume: 4},
			{Time: _d("2000-01-01 09:40:00"), Volume: 1},
		}},
		// 11:30 [11:25, 11:35)
		{[]Tdata{
			{Time: _d("2000-01-01 11:25:00"), Volume: 1},
			{Time: _d("2000-01-01 11:26:00"), Volume: 1},
			{Time: _d("2000-01-01 11:27:00"), Volume: 1},
			{Time: _d("2000-01-01 11:30:00"), Volume: 1},
			{Time: _d("2000-01-01 11:31:00"), Volume: 1},
		}, []Tdata{
			{Time: _d("2000-01-01 11:30:00"), Volume: 5},
		}},
		// 13:05 [13:00, 13:05)
		{[]Tdata{
			{Time: _d("2000-01-01 13:00:00"), Volume: 1},
			{Time: _d("2000-01-01 13:01:00"), Volume: 1},
			{Time: _d("2000-01-01 13:02:00"), Volume: 1},
			{Time: _d("2000-01-01 13:04:59"), Volume: 1},
			{Time: _d("2000-01-01 13:05:00"), Volume: 1},
		}, []Tdata{
			{Time: _d("2000-01-01 13:05:00"), Volume: 4},
			{Time: _d("2000-01-01 13:10:00"), Volume: 1},
		}},
		// 15:00 [14:55, 15:05)
		{[]Tdata{
			{Time: _d("2000-01-01 14:55:00"), Volume: 1},
			{Time: _d("2000-01-01 14:59:00"), Volume: 1},
			{Time: _d("2000-01-01 15:00:00"), Volume: 1},
			{Time: _d("2000-01-01 15:04:00"), Volume: 1},
		}, []Tdata{
			{Time: _d("2000-01-01 15:00:00"), Volume: 4},
		}},
	}
	for i, l := 0, len(tests); i < l; i++ {
		m1s := Tdatas{Data: tests[i].m1}
		m5s := Tdatas{}
		m5s.MergeFrom(&m1s, false, Minute5end)
		if len(m5s.Data) != len(tests[i].m5) {
			t.Error(
				"For", "case", i,
				"expected", "len eq",
				"got", "neq", m5s.Data,
			)
			continue
		}
		for k := 0; k < len(m5s.Data); k++ {
			if m5s.Data[k].Volume != tests[i].m5[k].Volume {
				t.Error(
					"For", "case", i, "data[", k, "]",
					"expected", "voleume eq",
					"got", "neq", m5s.Data[k],
				)
			}
		}
	}
}

func TestM1s2M30s(t *testing.T) {
	type Case struct {
		m1  []Tdata
		m30 []Tdata
	}

	tests := []Case{
		// 10:00 [09:00, 10:00)
		{[]Tdata{
			{Time: _d("2000-01-01 09:30:00"), Volume: 1},
			{Time: _d("2000-01-01 09:31:00"), Volume: 1},
			{Time: _d("2000-01-01 09:32:00"), Volume: 1},
			{Time: _d("2000-01-01 09:59:59"), Volume: 1},
			{Time: _d("2000-01-01 10:00:00"), Volume: 1},
		}, []Tdata{
			{Time: _d("2000-01-01 10:00:00"), Volume: 4},
			{Time: _d("2000-01-01 10:30:00"), Volume: 1},
		}},
		// 11:30 [11:00, 11:35)
		{[]Tdata{
			{Time: _d("2000-01-01 11:00:00"), Volume: 1},
			{Time: _d("2000-01-01 11:29:00"), Volume: 1},
			{Time: _d("2000-01-01 11:30:00"), Volume: 1},
			{Time: _d("2000-01-01 11:31:00"), Volume: 1},
		}, []Tdata{
			{Time: _d("2000-01-01 11:30:00"), Volume: 4},
		}},
		// 15:00 [14:30, 15:05)
		{[]Tdata{
			{Time: _d("2000-01-01 14:30:00"), Volume: 1},
			{Time: _d("2000-01-01 14:59:00"), Volume: 1},
			{Time: _d("2000-01-01 15:00:00"), Volume: 1},
			{Time: _d("2000-01-01 15:04:00"), Volume: 1},
		}, []Tdata{
			{Time: _d("2000-01-01 15:00:00"), Volume: 4},
		}},
	}
	for i, l := 0, len(tests); i < l; i++ {
		m1s := Tdatas{Data: tests[i].m1}
		m30s := Tdatas{}
		m30s.MergeFrom(&m1s, false, Minute30end)
		if len(m30s.Data) != len(tests[i].m30) {
			t.Error(
				"For", "case", i,
				"expected", "len eq",
				"got", "neq", m30s.Data,
			)
			continue
		}
		for k := 0; k < len(m30s.Data); k++ {
			if m30s.Data[k].Volume != tests[i].m30[k].Volume {
				t.Error(
					"For", "case", i, "data[", k, "]",
					"expected", "voleume eq",
					"got", "neq", m30s.Data[k],
				)
			}
		}
	}
}

func TestTicks2M1s(t *testing.T) {
	type date_pair struct {
		date []string
		len  int
	}
	tests := []date_pair{
		{[]string{
			"2000-01-01 10:00:00",
			"2000-01-02 10:00:00",
			"2000-01-02 10:01:00",
		}, 3},
	}

	fmt := "2006-01-02 15:04:05"
	for i, l := 0, len(tests); i < l; i++ {
		stock := Stock{}
		stock.Ticks = stringSlice2Ticks(tests[i].date)
		for j := 10; j > 0; j-- {
			stock.Ticks2M1s()
			if len(stock.M1s.Data) != tests[i].len {
				t.Error(
					"For", "case", i, tests[i],
					"expected", "m1s.data len=", tests[i].len,
					"got", len(stock.M1s.Data), stock.M1s.Data,
				)
			}
			for k := 0; k < len(tests[i].date); k++ {
				ts := stock.M1s.Data[k].Time.Format(fmt[:len(tests[i].date[k])])
				if ts != tests[i].date[k] {
					t.Error(
						"For", "case", i, "date[", k, "]",
						"expected", tests[i].date[k],
						"got", ts,
					)
				}
			}
		}
	}
}
