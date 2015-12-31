package crawl

import (
	"sort"
	"strconv"
	"time"
)

const (
	lmt   = "2006-01-02 15:04:05"
	smt   = "2006-01-02"
	l_lmt = len(lmt)
	l_smt = len(smt)
)

type Tdata struct {
	Time   time.Time `json:"time"`
	Open   int       `json:"open"`
	Close  int       `json:"close"`
	High   int       `json:"high"`
	Low    int       `json:"low"`
	Volume int       `json:"volume"`
	emas   int
	emal   int
	DIFF   int
	DEA    int
	MACD   int
}

type TdataSlice []Tdata

func (p TdataSlice) Len() int           { return len(p) }
func (p TdataSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TdataSlice) Less(i, j int) bool { return p[i].Time.Before(p[j].Time) }

func SearchTdataSliceByTime(a TdataSlice, t time.Time) int {
	return sort.Search(len(a), func(i int) bool {
		// a[i].Time >= t
		return a[i].Time.After(t) || a[i].Time.Equal(t)
	})
}

func (p TdataSlice) SearchByTime(t time.Time) (int, bool) {
	i := SearchTdataSliceByTime(p, t)
	if i < p.Len() {
		return i, t.Equal(p[i].Time)
	}
	return i, false
}

type Tdatas struct {
	Data    []Tdata `json:"data"`
	Typing  typing_parser
	Segment segment_parser
	Hub     hub_parser

	min_hub_height int
}

func (p *Tdatas) Add(data Tdata) {
	if len(p.Data) < 1 {
		p.Data = append(p.Data, data)
	} else if data.Time.After(p.Data[len(p.Data)-1].Time) {
		p.Data = append(p.Data, data)
	} else if data.Time.Equal(p.Data[len(p.Data)-1].Time) {
		p.Data[len(p.Data)-1] = data
	} else {
		j := len(p.Data) - 1
		should_insert := true
		for i := j - 1; i > -1; i-- {
			if p.Data[i].Time.After(data.Time) {
				j = i
				continue
			} else if p.Data[i].Time.Equal(data.Time) {
				p.Data[i] = data
				should_insert = false
			}
			break
		}

		if should_insert {
			if j < 1 {
				p.Data = append([]Tdata{data}, p.Data...)
			} else {
				p.Data = append(p.Data, data)
				copy(p.Data[j+1:], p.Data[j:])
				p.Data[j] = data
			}
		}
	}
}

func (p *Tdatas) latest_time() time.Time {
	if len(p.Data) < 1 {
		return market_begin_day
	}
	return p.Data[len(p.Data)-1].Time
}

func (p *Tdata) FromBytes(timestr, open, high, cloze, low, volume []byte) {
	p.FromString(string(timestr), string(open), string(high), string(cloze),
		string(low), string(volume))
}

func (p *Tdata) FromString(timestr, open, high, cloze, low, volume string) {
	if len(timestr) == l_lmt {
		p.Time, _ = time.Parse(lmt, timestr)
	} else {
		p.Time, _ = time.Parse(smt, timestr)
	}
	p.Open = ParseCent(open)
	p.High = ParseCent(high)
	p.Low = ParseCent(low)
	p.Close = ParseCent(cloze)
	p.Volume, _ = strconv.Atoi(volume)
}
