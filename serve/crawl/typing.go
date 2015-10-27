package crawl

import "time"

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

type typing_parser_node struct {
	d Tdata
	t Typing
}

type typing_parser struct {
	Data []Typing
	Line []Typing
	tp   []typing_parser_node
}

func (p *typing_parser) clear() {
	p.tp = []typing_parser_node{}
}

func (p *typing_parser) clean() {
	if len(p.tp) > 3 {
		var tmp []typing_parser_node
		tmp = append(tmp, p.tp[len(p.tp)-3:]...)
		p.tp = tmp
	}
}

func (p *typing_parser) new_node(i int, td *Tdatas) {
	if len(p.tp) > 0 {
		p.tp[len(p.tp)-1].t.End = i - 1
	}
	tp := typing_parser_node{}
	tp.t.begin = i
	tp.t.I = i
	tp.t.End = i
	tp.d = td.Data[i]
	p.tp = append(p.tp, tp)
}

func (p *typing_parser) parse_top_bottom() bool {
	if len(p.tp) < 3 {
		return false
	}
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
		return false
	}

	typing.High = b.High
	typing.Low = b.Low
	typing.Time = b.Time

	if len(p.Data) > 0 {
		if typing.I-p.Data[len(p.Data)-1].I < 4 {
			return false
		}

		if typing.Type == TopTyping && p.Data[len(p.Data)-1].Type == BottomTyping {
			if typing.High <= p.Data[len(p.Data)-1].High {
				return false
			}
		}

		if typing.Type == BottomTyping && p.Data[len(p.Data)-1].Type == TopTyping {
			if typing.High >= p.Data[len(p.Data)-1].High {
				return false
			}
		}

		if typing.Type == p.Data[len(p.Data)-1].Type {
			if pos, ok := TypingSlice(p.Data).MergeTyping(typing); ok {
				if pos < len(p.Data)-1 {
					p.Data = p.Data[:pos+1]
				}
				return true
			}
		}
	}
	p.Data = append(p.Data, typing)
	return true
}

func (p *Tdatas) findChanTypingStart() int {
	l := len(p.Data)
	if l < 240 {
		return 0
	}
	l = 240
	li, hi := 0, 0
	for i := 1; i < l; i++ {
		if p.Data[li].Low > p.Data[i].Low {
			li = i
		}
		if p.Data[hi].High < p.Data[i].High {
			hi = i
		}
	}

	for i := li - 1; i > -1 && p.Data[li].Low <= p.Data[i].Low; i-- {
		li = i
	}

	for i := hi - 1; i > -1 && p.Data[hi].High >= p.Data[i].High; i-- {
		hi = i
	}

	if hi > li {
		return li
	}
	return hi
}

func (p *Tdatas) ParseTyping() bool {
	hasnew := false
	start := 0
	if l := len(p.Typing.tp); l > 0 {
		start = p.Typing.tp[l-1].t.End + 1
	} else {
		start = p.findChanTypingStart()
	}

	for i, l := start, len(p.Data); i < l; i++ {
		a := &p.Data[i]

		if len(p.Typing.tp) < 1 {
			p.Typing.new_node(i, p)
			continue
		}

		prev := &p.Typing.tp[len(p.Typing.tp)-1]
		if Contain(&prev.d, a) {
			var base *Tdata
			if len(p.Typing.tp) > 1 {
				base = &p.Typing.tp[len(p.Typing.tp)-2].d
			} else {
				base = &Tdata{}
			}
			a = ContainMerge(base, &prev.d, a)
			if IsUpTyping(base, &prev.d) {
				if prev.d.High != a.High {
					prev.t.I = i
				}
			} else if IsDownTyping(base, &prev.d) {
				if prev.d.Low != a.Low {
					prev.t.I = i
				}
			}
			prev.d = *a
			prev.t.End = i
		} else {
			p.Typing.new_node(i, p)
		}

		p.Typing.clean()
		if p.Typing.parse_top_bottom() {
			hasnew = true
		}
	}
	return hasnew
}

func IsTopTyping(a, b, c *Tdata) bool {
	return IsUpTyping(a, b) && IsDownTyping(b, c)
}

func IsBottomTyping(a, b, c *Tdata) bool {
	return IsDownTyping(a, b) && IsUpTyping(b, c)
}

func IsUpTyping(a, b *Tdata) bool {
	return !Contain(a, b) && b.High > a.High
}

func IsDownTyping(a, b *Tdata) bool {
	return !Contain(a, b) && b.Low < a.Low
}

func Contain(a, b *Tdata) bool {
	// Fuzzy Lesson 67 答疑 2007-08-02 16:19:25
	// 缠中说禅：只要有一端相同，那必然是包含，
	// 两端相同那更是了，
	// 所以如果不是包含关系的，都必然不需要考虑相等关系
	return a.High == b.High || a.Low == b.Low || (a.High > b.High && a.Low < b.Low) || (a.High < b.High && a.Low > b.Low)
}

func DownContainMerge(a, b *Tdata) *Tdata {
	t := *a
	if b.Low < a.Low {
		t.Low = b.Low
		t.Time = b.Time
	}
	if b.High < a.High {
		t.High = b.High
	}
	return &t
}

func UpContainMerge(a, b *Tdata) *Tdata {
	t := *a
	if b.High > a.High {
		t.High = b.High
		t.Time = b.Time
	}
	if b.Low > a.Low {
		t.Low = b.Low
	}
	return &t
}

func ContainMerge(pra, a, b *Tdata) *Tdata {
	if IsUpTyping(pra, a) {
		return UpContainMerge(a, b)
	} else if IsDownTyping(pra, a) {
		return DownContainMerge(a, b)
	}
	return nil
}

func (p *typing_parser) LinkTyping() bool {
	hasnew := false
	start := 0
	typing := Typing{}
	if l := len(p.Line); l > 0 {
		typing = p.Line[l-1]
		if typing.Type == DownTyping {
			typing.Type = BottomTyping
		} else if typing.Type == UpTyping {
			typing.Type = TopTyping
		}
		for i := len(p.Data) - 1; i > -1; i-- {
			if p.Data[i].I == typing.I {
				start = i + 1
				break
			}
		}
	}

	end := len(p.Data) - 1
	for i := end - 1; i > -1; i-- {
		if p.Data[i].Type != p.Data[i+1].Type {
			end = i + 1
			break
		}
	}

	for i := start; i < end; i++ {
		t := p.Data[i]
		if typing.Type == UnknowTyping {
			typing = t
			continue
		}

		if typing.Type == t.Type {
			continue
		}

		typing.End = t.End
		if typing.Type == TopTyping {
			typing.Low = t.Low
			typing.Type = DownTyping
		} else if typing.Type == BottomTyping {
			typing.High = t.High
			typing.Type = UpTyping
		}
		p.Line = append(p.Line, typing)
		typing = t
		hasnew = true
	}

	return hasnew
}
