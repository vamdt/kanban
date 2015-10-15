package crawl

import (
	"log"
	"time"
)

const (
	UnknowTyping = iota
	WaitTyping
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
	High  int
	Low   int
	begin int
	End   int
}

type TypingSlice []Typing

func (p TypingSlice) Len() int           { return len(p) }
func (p TypingSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TypingSlice) Less(i, j int) bool { return p[i].I < p[j].I }

func (p TypingSlice) MergeTyping(t Typing) (int, bool) {
	pos := 0
	ok := false
	for i := len(p) - 1; i > -1; i-- {
		if t.Type == p[i].Type {
			if t.Type == TopTyping {
				if t.High > p[i].High {
					p[i] = t
					pos, ok = i, true
					continue
				}
			} else {
				if t.Low < p[i].Low {
					p[i] = t
					pos, ok = i, true
					continue
				}
			}
		}
		break
	}
	return pos, ok
}

type typing_parser struct {
	d Tdata
	t Typing
}

func (p *Tdatas) ParseTyping() {
	var prev *typing_parser
	start := 0
	if l := len(p.tp); l > 0 {
		start = p.tp[l-1].t.End + 1
		prev = &p.tp[l-1]
	} else {
		start = 0
	}

	for i, l := start, len(p.Data); i < l; i++ {
		a := &p.Data[i]

		if len(p.tp) < 1 {
			tp := typing_parser{}
			tp.t.begin = i
			tp.t.I = i
			tp.t.End = i
			tp.d = p.Data[i]
			tp.t.High = tp.d.High
			tp.t.Low = tp.d.Low
			tp.t.Time = tp.d.Time
			p.tp = append(p.tp, tp)
			prev = &p.tp[len(p.tp)-1]
		}

		if IsUpTyping(&prev.d, a) {
			prev.t.End = i - 1
			tp := typing_parser{}
			tp.t.begin = i
			tp.t.I = i
			tp.t.End = i
			tp.d = p.Data[i]
			tp.t.High = tp.d.High
			tp.t.Low = tp.d.Low
			tp.t.Time = tp.d.Time
			p.tp = append(p.tp, tp)
			prev = &p.tp[len(p.tp)-1]
		} else if IsDownTyping(&prev.d, a) {
			prev.t.End = i - 1
			tp := typing_parser{}
			tp.t.begin = i
			tp.t.I = i
			tp.t.End = i
			tp.d = *a
			tp.t.High = tp.d.High
			tp.t.Low = tp.d.Low
			tp.t.Time = tp.d.Time
			p.tp = append(p.tp, tp)
			prev = &p.tp[len(p.tp)-1]
		} else if Contain(&prev.d, a) {
			var base *Tdata
			if len(p.tp) > 1 {
				base = &p.tp[len(p.tp)-2].d
			} else {
				base = &Tdata{}
			}
			a = TypingMerge(base, &prev.d, a)
			if IsUpTyping(base, &prev.d) {
				if prev.d.High != a.High {
					prev.t.I = i
					prev.t.Time = p.Data[i].Time
				}
			} else if IsDownTyping(base, &prev.d) {
				if prev.d.Low != a.Low {
					prev.t.I = i
					prev.t.Time = p.Data[i].Time
				}
			}
			prev.d = *a
			prev.t.End = i
			prev.t.High = a.High
			prev.t.Low = a.Low
		} else {
			log.Println("UnknowTyping", a)
		}

		if len(p.tp) > 3 {
			var tmp []typing_parser
			tmp = append(tmp, p.tp[len(p.tp)-3:]...)
			p.tp = tmp
			prev = &p.tp[len(p.tp)-1]
		}

		if len(p.tp) > 2 {
			typing := p.tp[len(p.tp)-2].t
			a := &p.tp[len(p.tp)-3].d
			b := &p.tp[len(p.tp)-2].d
			c := &p.tp[len(p.tp)-1].d
			if IsTopTyping(a, b, c) {
				typing.Price = b.High
				typing.Type = TopTyping
			} else if IsBottomTyping(a, b, c) {
				typing.Price = b.Low
				typing.Type = BottomTyping
			} else {
				continue
			}

			typing.High = b.High
			typing.Low = b.Low

			if len(p.Typing) > 0 {
				if typing.I-p.Typing[len(p.Typing)-1].I < 4 {
					continue
				}

				if typing.Type == TopTyping && p.Typing[len(p.Typing)-1].Type == BottomTyping {
					if typing.High <= p.Typing[len(p.Typing)-1].High {
						continue
					}
				}

				if typing.Type == p.Typing[len(p.Typing)-1].Type {
					if pos, ok := TypingSlice(p.Typing).MergeTyping(typing); ok {
						if pos < len(p.Typing)-1 {
							p.Typing = p.Typing[:pos+1]
						}
						continue
					}
				}
			}
			p.Typing = append(p.Typing, typing)
		}
	}
}

func IsTopTyping(a, b, c *Tdata) bool {
	return IsUpTyping(a, b) && IsDownTyping(b, c)
}

func IsBottomTyping(a, b, c *Tdata) bool {
	return IsDownTyping(a, b) && IsUpTyping(b, c)
}

func IsUpTyping(a, b *Tdata) bool {
	return !Contain(a, b) && b.High >= a.High
}

func IsDownTyping(a, b *Tdata) bool {
	return !Contain(a, b) && b.Low <= a.Low
}

func Contain(a, b *Tdata) bool {
	//return (a.High >= b.High && a.Low <= b.Low) || (a.High <= b.High && a.Low >= b.Low)
	// Lesson 65
	return (a.High > b.High && a.Low < b.Low) || (a.High <= b.High && a.Low >= b.Low)
}

func TypingMerge(pra, a, b *Tdata) *Tdata {
	t := *a
	if IsUpTyping(pra, a) {
		if b.High > a.High {
			t.High = b.High
		}
		if b.Low > a.Low {
			t.Low = b.Low
		}
		return &t
	} else if IsDownTyping(pra, a) {
		if b.Low < a.Low {
			t.Low = b.Low
		}
		if b.High < a.High {
			t.High = b.High
		}
		return &t
	}
	return nil
}
