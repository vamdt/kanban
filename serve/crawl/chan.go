package crawl

import "time"

const (
	UnknowTyping = iota
	TopTyping
	BottomTyping
	UpTyping
	DownTyping
)

type Typing struct {
	I     int
	Time  time.Time
	Price int
	Type  int
}

func (p *Stock) Chan() {
	p.M1s.ParseTyping()
	p.M5s.ParseTyping()
}

func (p *Tdatas) ParseTyping() {
	start := 0
	typing := Typing{}
	if l := len(p.Typing); l > 0 {
		start = p.Typing[l-1].I
	}
	var pra *Tdata
	for i, l := start+2, len(p.Data); i < l; i++ {
		if i > 2 {
			pra = &p.Data[i-3]
		} else {
			pra = nil
		}
		a := &p.Data[i-2]
		b := &p.Data[i-1]
		c := &p.Data[i]

		for Contain(a, b) {
			if pra != nil {
				a = PraTypingMerge(pra, a, b)
			} else {
				a = PobTypingMerge(a, b, c)
			}
			i++
			if i >= l {
				return
			}
			b = c
			c = &p.Data[i]
		}

		pra = a
		if IsTopTyping(a, b, c) {
			typing.Price = b.High
			typing.Type = TopTyping
		} else if IsBottomTyping(a, b, c) {
			typing.Price = b.Low
			typing.Type = BottomTyping
		} else {
			continue
		}
		typing.I = i - 1
		if len(p.Typing) > 0 && typing.I-p.Typing[len(p.Typing)-1].I < 4 {
			continue
		}
		typing.Time = b.Time
		p.Typing = append(p.Typing, typing)
	}
}

func IsTopTyping(a, b, c *Tdata) bool {
	return b.High > a.High && b.High > c.High && b.Low > a.Low && b.Low > c.Low
}

func IsBottomTyping(a, b, c *Tdata) bool {
	return b.High < a.High && b.High < c.High && b.Low < a.Low && b.Low < c.Low
}

func IsUpTyping(a, b *Tdata) bool {
	return b.High > a.High && b.Low >= a.Low
}

func IsDownTyping(a, b *Tdata) bool {
	return b.Low < a.Low && b.High <= a.High
}

func Contain(a, b *Tdata) bool {
	return (a.High > b.High && a.Low < b.Low) || (a.High < b.High && a.Low > b.Low)
}

func PraTypingMerge(pra, a, b *Tdata) *Tdata {
	t := *a
	if IsUpTyping(pra, a) {
		if b.High > a.High {
			t.High = b.High
		}
		if b.Low > a.Low {
			t.Low = b.Low
		}
	} else if IsDownTyping(pra, a) {
		if b.Low < a.Low {
			t.Low = b.Low
		}
		if b.High < a.High {
			t.High = b.High
		}
	} else {
		return nil
	}
	return &t
}

func PobTypingMerge(a, b, pob *Tdata) *Tdata {
	t := *a
	if IsUpTyping(b, pob) {
		if b.High > a.High {
			t.High = b.High
		}
		if b.Low > a.Low {
			t.Low = b.Low
		}
	} else if IsDownTyping(b, pob) {
		if b.Low < a.Low {
			t.Low = b.Low
		}
		if b.High < a.High {
			t.High = b.High
		}
	} else {
		return nil
	}
	return &t
}
