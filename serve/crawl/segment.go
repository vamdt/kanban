package crawl

import "log"

func (p *typing_parser) new_segment_node(i int, ptyping *typing_parser) {
	if len(p.tp) > 0 {
		p.tp[len(p.tp)-1].t.End = i - 1
	}
	tp := typing_parser_node{}
	tp.t = ptyping.Line[i]
	tp.t.begin = i
	tp.t.I = i
	tp.t.End = i
	tp.d.Time = tp.t.Time
	tp.d.High = tp.t.High
	tp.d.Low = tp.t.Low
	p.tp = append(p.tp, tp)
}

func (p *Tdatas) ParseSegment() bool {
	hasnew := false
	start := 0
	if l := len(p.Segment.tp); l > 0 {
		start = p.Segment.tp[l-1].t.End + 1
	} else {
		start = 0
	}

	l := len(p.Typing.Line)
	if l > 0 && p.Typing.Line[l-1].Type != UpTyping && p.Typing.Line[l-1].Type != DownTyping {
		l--
	}

	for i := start; i+1 < l; i += 2 {

		if len(p.Segment.tp) < 1 {
			if i+2 > l {
				return hasnew
			}

			if len(p.Segment.Data) > 0 {
				p.Segment.new_segment_node(i, &p.Typing)
				continue
			}

			if p.Typing.Line[i].Type == UpTyping && p.Typing.Line[i+2].Low < p.Typing.Line[i].High {
				// Up yes
				i++
				p.Segment.new_segment_node(i, &p.Typing)
			} else if p.Typing.Line[i].Type == DownTyping && p.Typing.Line[i+2].High > p.Typing.Line[i].Low {
				// Down yes
				i++
				p.Segment.new_segment_node(i, &p.Typing)
			} else {
				i--
			}
			continue
		}

		prev := &p.Segment.tp[len(p.Segment.tp)-1]
		if i == prev.t.End+1 {
			continue
		}
		a := &Tdata{}
		a.High = p.Typing.Line[i].High
		a.Low = p.Typing.Line[i].Low
		a.Time = p.Typing.Line[i].Time
		if IsUpTyping(&prev.d, a) {
			p.Segment.new_segment_node(i, &p.Typing)
		} else if IsDownTyping(&prev.d, a) {
			p.Segment.new_segment_node(i, &p.Typing)
		} else if Contain(&prev.d, a) {
			var base *Tdata
			if len(p.Segment.tp) > 1 {
				base = &p.Segment.tp[len(p.Segment.tp)-2].d
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
			log.Println("UnknowTyping", a)
		}

		p.Segment.clean()
		if p.Segment.parse_segment_top_bottom() {
			hasnew = true
			i--
			p.Segment.clear()
		}
	}
	return hasnew
}

func (p *typing_parser) parse_segment_top_bottom() bool {
	if len(p.tp) < 3 {
		return false
	}
	typing := p.tp[len(p.tp)-2].t
	a := &p.tp[len(p.tp)-3].d
	b := &p.tp[len(p.tp)-2].d
	c := &p.tp[len(p.tp)-1].d
	if typing.Type == UpTyping && IsBottomTyping(a, b, c) {
		typing.Price = b.Low
		typing.Type = BottomTyping
	} else if typing.Type == DownTyping && IsTopTyping(a, b, c) {
		typing.Price = b.High
		typing.Type = TopTyping
	} else {
		return false
	}

	typing.High = b.High
	typing.Low = b.Low
	typing.Time = b.Time

	if len(p.Data) > 0 {
		if typing.Type == TopTyping && p.Data[len(p.Data)-1].Type == BottomTyping {
			if typing.High <= p.Data[len(p.Data)-1].High {
				log.Println("find a bottom high then top")
			}
		}
		if false && typing.Type == p.Data[len(p.Data)-1].Type {
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
