package crawl

import "time"

const (
	_ = iota
	TopTyping
	BottomTyping
)

type Typing struct {
	I     int
	Time  time.Time
	Price int
	Type  int
}

func (p *Stock) Chan() {
	p.M1s.Chan()
	p.M5s.Chan()
}

func (p *Tdatas) Chan() {
	start := 0
	typing := Typing{}
	if l := len(p.Typing); l > 0 {
		start = p.Typing[l-1].I
	}
	for i, c := start+2, len(p.Data); i < c; i++ {
		if IsTopTyping(&p.Data[i-2], &p.Data[i-1], &p.Data[i]) {
			typing.I = i - 1
			typing.Time = p.Data[i-1].Time
			typing.Price = p.Data[i-1].High
			typing.Type = TopTyping
			p.Typing = append(p.Typing, typing)
		} else if IsBottomTyping(&p.Data[i-2], &p.Data[i-1], &p.Data[i]) {
			typing.I = i - 1
			typing.Time = p.Data[i-1].Time
			typing.Price = p.Data[i-1].Low
			typing.Type = BottomTyping
			p.Typing = append(p.Typing, typing)
		}
	}
}

func IsTopTyping(a, b, c *Tdata) bool {
	return b.High > a.High && b.High > c.High && b.Low > a.Low && b.Low > c.Low
}

func IsBottomTyping(a, b, c *Tdata) bool {
	return b.High < a.High && b.High < c.High && b.Low < a.Low && b.Low < c.Low
}
