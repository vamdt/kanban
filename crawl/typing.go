package crawl

import (
	"sort"
	"time"

	"github.com/golang/glog"
)

// 走势分为趋势和盘整
// 趋势分为上涨和下跌
const (
	UnknowTyping int = iota
	WaitTyping
	TopTyping
	BottomTyping
	UpTyping
	DownTyping
	DullTyping
)

type Typing struct {
	I     int
	Time  time.Time
	Price int
	Type  int
	High  int
	Low   int
	begin int
	end   int
	ETime time.Time
	Case1 bool
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

func SearchTypingSliceByTime(a TypingSlice, t time.Time) int {
	return sort.Search(len(a), func(i int) bool {
		return a[i].Time.After(t) || a[i].Time.Equal(t)
	})
}

func (p TypingSlice) SearchByTime(t time.Time) (int, bool) {
	i := SearchTypingSliceByTime(p, t)
	if i < p.Len() {
		return i, t.Equal(p[i].Time)
	}
	return i, false
}

func SearchTypingSliceByETime(a TypingSlice, t time.Time) int {
	return sort.Search(len(a), func(i int) bool {
		return a[i].ETime.After(t) || a[i].ETime.Equal(t)
	})
}

func (p TypingSlice) SearchByETime(t time.Time) (int, bool) {
	i := SearchTypingSliceByETime(p, t)
	if i < p.Len() {
		return i, t.Equal(p[i].ETime)
	}
	return i, false
}

type typing_parser_node struct {
	d Tdata
	t Typing
}

type typing_parser struct {
	Data []Typing
	Line []Typing
	tp   []typing_parser_node
	tag  string
}

func (p *typing_parser) parser_reset() {
	p.tp = []typing_parser_node{}
}

func (p *typing_parser) clean() {
	if len(p.tp) > 3 {
		p.tp = p.tp[len(p.tp)-3:]
	}
}

func (p *typing_parser) new_node(i int, td *Tdatas) {
	if len(p.tp) > 0 {
		p.tp[len(p.tp)-1].t.end = i - 1
		p.tp[len(p.tp)-1].t.ETime = td.Data[i-1].Time
		p.tp[len(p.tp)-1].t.assertETimeMatchEnd(td.Data, "new_node prev end")
	}
	tp := typing_parser_node{}
	tp.t.begin = i
	tp.t.I = i
	tp.t.end = i
	tp.t.ETime = td.Data[i].Time
	tp.t.assertETimeMatchEnd(td.Data, "new_node")
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
	p.parser_reset()
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

func (p *Typing) assertETimeMatchEnd(data TdataSlice, note string) int {
	i, ok := data.SearchByTime(p.ETime)
	if ok {
		if p.end != i {
			glog.Fatalf("%s assert end/%d eq SearchByTime/%d", note, p.end, i)
		}
	} else {
		glog.Fatalln("%s not found with time", note, p.ETime, data[len(data)-20:])
	}
	return i
}

func (p *Typing) assertETimeMatchEndLine(data TypingSlice, note string) int {
	i, ok := data.SearchByETime(p.ETime)
	if ok {
		if p.end != i {
			glog.Fatalf("%s assert end/%d eq SearchByETime/%d", note, p.end, i)
		}
	} else {
		glog.Fatalln("not found with etime", note, p.ETime)
	}
	return i
}

func (p *Tdatas) ParseTyping() bool {
	hasnew := false
	start := 0
	if l := len(p.Typing.Data); l > 0 {
		p.Typing.Data = p.Typing.Data[:l-1]
	}
	if l := len(p.Typing.Data); l > 0 {
		start = p.Typing.Data[l-1].end + 1
		start = 1 + p.Typing.Data[l-1].assertETimeMatchEnd(p.Data, "ParseTyping start2")
	} else {
		start = p.findChanTypingStart()
	}

	p.Typing.parser_reset()
	for i, l := start, len(p.Data); i < l; i++ {
		a := &p.Data[i]

		ltp := len(p.Typing.tp)
		if ltp < 1 {
			p.Typing.new_node(i, p)
			continue
		}

		prev := &p.Typing.tp[ltp-1]
		if Contain(&prev.d, a) {
			var base *Tdata
			if ltp > 1 {
				base = &p.Typing.tp[ltp-2].d
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
			prev.t.end = i
			prev.t.ETime = p.Data[i].Time
			prev.t.assertETimeMatchEnd(p.Data, "ParseTyping Contain")
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

func LineContain(a, b *Typing) bool {
	ta := Tdata{High: a.High, Low: a.Low}
	tb := Tdata{High: b.High, Low: b.Low}
	return Contain(&ta, &tb)
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

	if l := len(p.Line); l > 0 {
		if l > 5 {
			p.Line = p.Line[0 : l-5]
		} else {
			p.Line = []Typing{}
		}
	}

	if l := len(p.Line); l > 0 {
		t := p.Line[l-1]
		for i := len(p.Data) - 1; i > -1; i-- {
			if p.Data[i].end == t.end {
				start = i
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

	typing := Typing{}
	for i := start; i < end; i++ {
		t := p.Data[i]
		if typing.Type == UnknowTyping {
			typing = t
			continue
		}

		if typing.Type == t.Type {
			continue
		}

		typing.end = t.end
		typing.ETime = t.ETime
		if typing.Type == TopTyping {
			typing.Low = t.Low
			typing.Type = DownTyping
		} else if typing.Type == BottomTyping {
			typing.High = t.High
			typing.Type = UpTyping
		} else {
			glog.Fatalf("%s typing.Type=%d should be %d or %d", p.tag, typing.Type, TopTyping, BottomTyping)
		}
		p.Line = append(p.Line, typing)
		typing = t
		hasnew = true
	}

	return hasnew
}