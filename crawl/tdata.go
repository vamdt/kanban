package crawl

import (
	"time"

	. "./base"
)

type Tdatas struct {
	Data    []Tdata `json:"data"`
	Typing  typing_parser
	Segment segment_parser
	Hub     hub_parser
	tag     string
	start   time.Time

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

func (p *Tdatas) Add(data Tdata) (int, bool) {
	if data.Volume == 0 && data.Open == 0 {
		return 0, false
	}

	l := len(p.Data)
	if l < 1 {
		p.Data = append(p.Data, data)
		return 0, true
	}

	if data.Time.After(p.Data[l-1].Time) {
		p.Data = append(p.Data, data)
		return l, true
	}

	if data.Time.Equal(p.Data[l-1].Time) {
		p.Data[l-1] = data
		return l - 1, false
	}

	i, ok := (TdataSlice(p.Data)).Search(data.Time)
	if ok {
		p.Data[i] = data
		return i, false
	}

	if i < 1 {
		p.Data = append([]Tdata{data}, p.Data...)
	} else {
		p.Data = append(p.Data, data)
		copy(p.Data[i+1:], p.Data[i:])
		p.Data[i] = data
	}
	return i, true
}

func (p *Tdatas) latest_time() time.Time {
	if len(p.Data) < 1 {
		return market_begin_day
	}
	return p.Data[len(p.Data)-1].Time
}
