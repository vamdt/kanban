package crawl

import (
	"sort"
	"strconv"
	"time"
)

const (
	lmt    = "2006-01-02 15:04:05"
	smt    = "2006-01-02"
	qqmt   = "060102"
	jqkamt = "20060102"
	l_lmt  = len(lmt)
	l_smt  = len(smt)
	l_qqmt = len(qqmt)
	l_jqka = len(jqkamt)
)

type Tdata struct {
	Time   time.Time `json:"time"`
	Open   int       `json:"open"`
	Close  int       `json:"close"`
	Volume int       `json:"volume"`
	HL     `bson:",inline"`
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

func SearchTdataSlice(a TdataSlice, t time.Time) int {
	return sort.Search(len(a), func(i int) bool {
		// a[i].Time >= t
		return a[i].Time.After(t) || a[i].Time.Equal(t)
	})
}

func (p TdataSlice) Search(t time.Time) (int, bool) {
	i := SearchTdataSlice(p, t)
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
	tag     string

	min_hub_height int

	base, next *Tdatas
}

func (p *typing_parser) tail(s *typing_parser, tail int) {
	if l := len(p.Data); l > 0 {
		start := l - tail
		if start < 0 || start >= l {
			start = 0
		}
		s.Data = p.Data[start:]
	}

	if l := len(p.Line); l > 0 {
		start := l - tail
		if start < 0 || start >= l {
			start = 0
		}
		s.Line = p.Line[start:]
	}
}

func (p *Tdatas) tail(s *Tdatas, tail int) {
	l := len(p.Data)
	if l < 1 {
		return
	}
	start := l - tail
	if start < 0 || start >= l {
		start = 0
	}
	s.Data = p.Data[start:]
	p.Typing.tail(&s.Typing, tail)
	p.Segment.tail(&s.Segment.typing_parser, tail)
	p.Hub.tail(&s.Hub.typing_parser, tail)
}

func (p *Tdatas) Init(hub_height int, tag string, base, next *Tdatas) {
	p.min_hub_height = hub_height
	p.tag = tag
	p.Typing.tag = tag
	p.Segment.tag = tag
	p.Hub.tag = tag
	p.base = base
	p.next = next
}

func (p *Tdatas) First_lastday_data() int {
	ldata := len(p.Data)
	if ldata < 1 {
		return 0
	}

	start := ldata - 240 - 10
	if start < 0 {
		start = 0
	}
	t := p.Data[ldata-1].Time.Truncate(time.Hour * 24)
	i, _ := (TdataSlice(p.Data[start:])).Search(t)
	return i + start
}

func (p *Tdatas) Drop_lastday_data() {
	ldata := len(p.Data)
	if ldata < 1 {
		return
	}

	i := p.First_lastday_data()
	if i < 1 {
		p.Data = []Tdata{}
		return
	}
	p.Data = p.Data[:i]
}

func (p *Tdatas) Add(data Tdata) int {
	l := len(p.Data)
	if l < 1 {
		p.Data = append(p.Data, data)
		return 0
	}

	if data.Time.After(p.Data[l-1].Time) {
		p.Data = append(p.Data, data)
		return l
	}

	if data.Time.Equal(p.Data[l-1].Time) {
		p.Data[l-1] = data
		return l - 1
	}

	i, ok := (TdataSlice(p.Data)).Search(data.Time)
	if ok {
		p.Data[i] = data
		return i
	}

	if i < 1 {
		p.Data = append([]Tdata{data}, p.Data...)
	} else {
		p.Data = append(p.Data, data)
		copy(p.Data[i+1:], p.Data[i:])
		p.Data[i] = data
	}
	return i
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
	ltime := len(timestr)
	switch ltime {
	case l_lmt:
		p.Time, _ = time.Parse(lmt, timestr)
	case l_smt:
		p.Time, _ = time.Parse(smt, timestr)
	case l_qqmt:
		p.Time, _ = time.Parse(qqmt, timestr)
	case l_jqka:
		p.Time, _ = time.Parse(jqkamt, timestr)
	}
	p.Open = ParseCent(open)
	p.High = ParseCent(high)
	p.Low = ParseCent(low)
	p.Close = ParseCent(cloze)
	p.Volume, _ = strconv.Atoi(volume)
}
