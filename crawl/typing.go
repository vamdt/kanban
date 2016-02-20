package crawl

import (
	"flag"
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

type HL struct {
	High int
	Low  int
}

type Typing struct {
	begin int
	i     int
	end   int
	Time  time.Time
	Price int
	Type  int
	HL    `bson:",inline"`
	ETime time.Time
	Case1 bool
	b1    int
	e3    int
}

var strict_line bool = true

func init() {
	flag.BoolVar(&strict_line, "strict_line", true, "link typing with strict rule")
}

type TypingSlice []Typing

func (p TypingSlice) Len() int           { return len(p) }
func (p TypingSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TypingSlice) Less(i, j int) bool { return p[i].i < p[j].i }

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

func (p *typing_parser) drop_last_5_data() {
	l := len(p.Data)
	if l < 1 {
		return
	}

	if l > 5 {
		p.Data = p.Data[0 : l-5]
	} else {
		p.Data = []Typing{}
	}
}

func (p *typing_parser) drop_last_5_line() {
	l := len(p.Line)
	if l < 1 {
		return
	}

	if l > 5 {
		p.Line = p.Line[0 : l-5]
	} else {
		p.Line = []Typing{}
	}
}

func (p *typing_parser) parser_reset() {
	p.tp = []typing_parser_node{}
}

func (p *typing_parser) clean() {
	if len(p.tp) > 3 {
		p.tp = p.tp[len(p.tp)-3:]
	}
}

func (p *typing_parser) new_node(t typing_parser_node) {
	p.tp = append(p.tp, t)
}

func (p *typing_parser) parse_top_bottom() bool {
	if len(p.tp) < 3 {
		return false
	}
	typing := p.tp[len(p.tp)-2].t
	a := &p.tp[len(p.tp)-3].t
	b := &p.tp[len(p.tp)-2].t
	c := &p.tp[len(p.tp)-1].t
	if IsTopTyping(a.HL, b.HL, c.HL) {
		typing.Price = b.High
		typing.Type = TopTyping
	} else if IsBottomTyping(a.HL, b.HL, c.HL) {
		typing.Price = b.Low
		typing.Type = BottomTyping
	} else {
		return false
	}

	typing.b1 = p.tp[len(p.tp)-3].t.begin
	typing.e3 = p.tp[len(p.tp)-1].t.end
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

func (p *Typing) assertETimeMatchEnd(data TdataSlice, note string) int {
	i, ok := data.Search(p.ETime)
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

func (p *Tdatas) ReadContainedTdata(base HL, i int) (typing_parser_node, bool) {
	l := len(p.Data)
	n := typing_parser_node{}
	if i >= l {
		return n, false
	}

	a := p.Data[i]
	n.t.begin = i
	n.t.i = i
	n.t.end = i
	n.t.Time = a.Time
	n.t.HL = a.HL
	n.t.ETime = a.Time

	for i = i + 1; i < l; i++ {
		a := p.Data[i]
		if !Contain(n.t.HL, a.HL) {
			break
		}

		n.t.end = i
		n.t.ETime = a.Time
		if IsUpTyping(base, n.t.HL) {
			a.HL = UpContainMergeHL(n.t.HL, a.HL)
			if n.t.High < a.High {
				n.t.i = i
				n.t.Time = a.Time
			}
		} else {
			a.HL = DownContainMergeHL(n.t.HL, a.HL)
			if n.t.Low > a.Low {
				n.t.i = i
				n.t.Time = a.Time
			}
		}
		n.t.HL = a.HL
	}
	return n, true
}

// Lesson 62, 65, <b>77</b>
func (p *Tdatas) ParseTyping() {
	start := 0

	p.Typing.drop_last_5_data()
	p.Typing.parser_reset()

	base := HL{}
	if l := len(p.Typing.Data); l > 0 {
		typing := p.Typing.Data[l-1]
		start = typing.end + 1
		p.Typing.Data[l-1].assertETimeMatchEnd(p.Data, "ParseTyping start2")

		base = typing.HL
		if typing.Type == TopTyping {
			base.High = base.High - 1
			base.Low = base.Low - 1
		} else {
			base.High = base.High + 1
			base.Low = base.Low + 1
		}
		t, _ := p.ReadContainedTdata(base, typing.begin)
		p.Typing.new_node(t)
	} else {
		start = p.findChanTypingStart()
	}

	glog.Infof("start %d/%d", start, len(p.Data))
	for i, l := start, len(p.Data); i < l; i++ {
		t, ok := p.ReadContainedTdata(base, i)
		if !ok {
			break
		}
		i = t.t.end
		//t.assertETimeMatchEnd(p.Data, "ParseTyping Contain")
		p.Typing.new_node(t)
		base = t.t.HL
		p.Typing.parse_top_bottom()
		p.Typing.clean()
	}
}

func IsTopTyping(a, b, c HL) bool {
	return IsUpTyping(a, b) && IsDownTyping(b, c)
}

func IsBottomTyping(a, b, c HL) bool {
	return IsDownTyping(a, b) && IsUpTyping(b, c)
}

func IsUpTyping(a, b HL) bool {
	return !Contain(a, b) && b.High > a.High
}

func IsDownTyping(a, b HL) bool {
	return !Contain(a, b) && b.Low < a.Low
}

func Contain(a, b HL) bool {
	// Fuzzy Lesson 67 答疑 2007-08-02 16:19:25
	// 缠中说禅：只要有一端相同，那必然是包含，
	// 两端相同那更是了，
	// 所以如果不是包含关系的，都必然不需要考虑相等关系
	return a.High == b.High || a.Low == b.Low || (a.High > b.High && a.Low < b.Low) || (a.High < b.High && a.Low > b.Low)
}

func DownContainMergeHL(a, b HL) HL {
	if b.Low < a.Low {
		a.Low = b.Low
	}
	if b.High < a.High {
		a.High = b.High
	}
	return a
}

func UpContainMergeHL(a, b HL) HL {
	if b.High > a.High {
		a.High = b.High
	}
	if b.Low > a.Low {
		a.Low = b.Low
	}
	return a
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

// Lesson 65, 77
func (p *typing_parser) LinkTyping() {
	p.drop_last_5_line()

	start := 0
	if l := len(p.Line); l > 0 {
		if i, ok := TypingSlice(p.Data).SearchByETime(p.Line[l-1].ETime); ok {
			start = i
		}
	}

	end := len(p.Data)
	typing := Typing{}
	for i := start; i < end; i++ {
		t := p.Data[i]
		if typing.Type == UnknowTyping {
			typing = t
			continue
		}

		if typing.Type == t.Type {
			if t.Type == TopTyping {
				if t.High > typing.High {
					if l := len(p.Line); l > 0 {
						p.Line[l-1].High = t.High
						p.Line[l-1].end = t.end
						p.Line[l-1].ETime = t.ETime
					}
					typing = t
					continue
				}
			} else {
				if t.Low < typing.Low {
					if l := len(p.Line); l > 0 {
						p.Line[l-1].Low = t.Low
						p.Line[l-1].end = t.end
						p.Line[l-1].ETime = t.ETime
					}
					typing = t
					continue
				}
			}
			continue
		}

		// check typing
		// 笔定义 第1条 分型不共用k线
		if t.b1 <= typing.e3 {
			continue
		}

		// 旧笔定义 有一个独立k线
		if strict_line && t.b1-typing.e3 < 2 {
			continue
		}
		// 新笔定义 第2条 Lesson 81 答疑部分
		if t.i-typing.i < 4 {
			continue
		}

		// Lesson 77
		if t.Type == TopTyping && t.High <= typing.High {
			continue
		}

		if t.Type == BottomTyping && t.High >= typing.High {
			continue
		}
		// check typing end

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
		typing.begin = t.begin
	}
}
