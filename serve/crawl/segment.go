package crawl

import "log"

type segment_parser struct {
	Data []Typing
	tp   []typing_parser_node
}

func (p *segment_parser) new_node(i int, ptyping *typing_parser) {
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

func (p *segment_parser) clear() {
	p.tp = []typing_parser_node{}
}

func (p *segment_parser) clean() {
	if len(p.tp) > 3 {
		var tmp []typing_parser_node
		tmp = append(tmp, p.tp[len(p.tp)-3:]...)
		p.tp = tmp
	}
}

func (p *Tdatas) IsLineBreakSegment(li int) bool {
	cl := len(p.Typing.Line)
	if cl-1 < li {
		return false
	}

	// find segment start
	si := 0
	for i := len(p.Segment.Data) - 1; i > -1; i-- {
		if p.Segment.Data[i].I <= li {
			si = p.Segment.Data[i].I
			break
		}
	}
	if si < 0 {
		si = 0
	}

	// type check
	if p.Typing.Line[li].Type != UpTyping && p.Typing.Line[li].Type != DownTyping {
		return false
	}

	if p.Typing.Line[si].Type != UpTyping && p.Typing.Line[si].Type != DownTyping {
		return false
	}

	if p.Typing.Line[li].Type == p.Typing.Line[si].Type {
		return false
	}

	// check break point
	if p.Typing.Line[si].Type == UpTyping {
		for i := li-2; i > si; i-- {
			if p.Typing.Line[i].High >= p.Typing.Line[li].Low {
				return true
			}
		}
	} else if p.Typing.Line[si].Type == DownTyping {
		for i := li-2; i > si; i-- {
			if p.Typing.Line[i].Low <= p.Typing.Line[li].High {
				return true
			}
		}
	}
	return false
}

func (p *Tdatas) feat_normalized(start, end int) {
	l := len(p.Typing.Line)
	if l > 0 && p.Typing.Line[l-1].Type != UpTyping && p.Typing.Line[l-1].Type != DownTyping {
		l--
	}
	if l < end {
		end = l
	}
	if start < 0 {
		log.Println("start < 0")
		return
	}

	p.Segment.tp = []typing_parser_node{}
	for i := start + 1; i < end; i += 2 {

		if len(p.Segment.tp) < 1 {
			p.Segment.new_node(i, &p.Typing)
			continue
		}

		prev := &p.Segment.tp[len(p.Segment.tp)-1]
		a := &Tdata{}
		a.High = p.Typing.Line[i].High
		a.Low = p.Typing.Line[i].Low
		a.Time = p.Typing.Line[i].Time

		if Contain(&prev.d, a) {
			if prev.t.Type == UpTyping {
				a = DownContainMerge(&prev.d, a)
				if prev.d.Low != a.Low {
					prev.t.I = i
				}
			} else {
				if prev.t.Type != DownTyping {
					log.Panicf("prev should be a DownTyping line %+v", prev)
				}
				a = UpContainMerge(&prev.d, a)
				if prev.d.High != a.High {
					prev.t.I = i
				}
			}
			prev.d = *a
			prev.t.End = i
		} else {
			p.Segment.new_node(i, &p.Typing)
		}

		p.Segment.clean()
	}
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
				p.Segment.new_node(i, &p.Typing)
				continue
			}

			if p.Typing.Line[i].Type == UpTyping && p.Typing.Line[i+2].Low < p.Typing.Line[i].High {
				// Up yes
				i++
				p.Segment.new_node(i, &p.Typing)
			} else if p.Typing.Line[i].Type == DownTyping && p.Typing.Line[i+2].High > p.Typing.Line[i].Low {
				// Down yes
				i++
				p.Segment.new_node(i, &p.Typing)
			} else {
				i--
			}
			continue
		}

		prev := &p.Segment.tp[len(p.Segment.tp)-1]
		a := &Tdata{}
		a.High = p.Typing.Line[i].High
		a.Low = p.Typing.Line[i].Low
		a.Time = p.Typing.Line[i].Time

		if dlen := len(p.Segment.Data); dlen > 0 {
			need_override_prev_typing := false
			if prev.t.Type == UpTyping && a.High >= p.Segment.Data[dlen-1].High {
				need_override_prev_typing = true
			}
			if prev.t.Type == DownTyping && a.Low <= p.Segment.Data[dlen-1].Low {
				need_override_prev_typing = true
			}
			if need_override_prev_typing {
				start := 0
				if dlen > 1 {
					start = p.Segment.Data[dlen-2].I
				} else {
					start = p.Segment.Data[dlen-1].I - 3
				}
				p.feat_normalized(start, i)
				p.Segment.Data = p.Segment.Data[:dlen-1]
				i--
				continue
			}
		}

		if Contain(&prev.d, a) {
      if len(p.Segment.tp) > 1 {
        firstIsBreak := false
        if len(p.Segment.tp) > 2 {
          pprev := &p.Segment.tp[len(p.Segment.tp)-2]
          if p.IsLineBreakSegment(pprev.t.I) {
            firstIsBreak = true
          }
        }
        if !firstIsBreak && p.IsLineBreakSegment(prev.t.I) {
          p.Segment.new_node(i, &p.Typing)
          continue
        }
      }

			if prev.t.Type == UpTyping {
				a = DownContainMerge(&prev.d, a)
				if prev.d.Low != a.Low {
					prev.t.I = i
				}
			} else {
				if prev.t.Type != DownTyping {
					log.Panicf("prev should be a DownTyping line %+v", prev)
				}
				a = UpContainMerge(&prev.d, a)
				if prev.d.High != a.High {
					prev.t.I = i
				}
			}
			prev.d = *a
			prev.t.End = i
		} else {
			if len(p.Segment.tp) < 3 && p.IsLineBreakSegment(i) {
				if prev.t.Type == DownTyping && a.High < prev.d.High {
					continue
				}
				if prev.t.Type == UpTyping && a.Low > prev.d.Low {
					continue
				}
			}
			p.Segment.new_node(i, &p.Typing)
		}

		p.Segment.clean()
		rv := p.Segment.parse_top_bottom(p)
		if rv == 1 {
			hasnew = true
			i = p.Segment.tp[len(p.Segment.tp)-2].t.I + 1
			p.Segment.clear()
			p.Segment.new_node(i, &p.Typing)
		} else if rv == 2 {
			hasnew = true
			i--
		}
	}
	return hasnew
}

func hasGap(a, b *Tdata) bool {
	return a.Low > b.High || a.High < b.Low
}

func (p *segment_parser) parse_top_bottom(tds *Tdatas) int {
	if len(p.tp) < 3 {
		return 0
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
		return 0
	}

	typing.High = b.High
	typing.Low = b.Low
	typing.Time = b.Time

	dlen := len(p.Data)
	if dlen > 0 {
		if typing.Type == TopTyping && p.Data[dlen-1].Type == BottomTyping {
			if typing.High <= p.Data[dlen-1].High {
				log.Println("find a bottom high then top")
			}
		}

		need_override_prev_typing := false
		if typing.Type == BottomTyping && c.High >= p.Data[dlen-1].High {
			need_override_prev_typing = true
		}
		if typing.Type == TopTyping && c.Low <= p.Data[dlen-1].Low {
			need_override_prev_typing = true
		}
		if need_override_prev_typing {
			start := 0
			end := typing.I
			if dlen > 1 {
				start = p.Data[dlen-2].I
			} else {
				start = p.Data[dlen-1].I - 3
			}
			tds.feat_normalized(start, end)
			p.Data = p.Data[:dlen-1]
			return 2
		}
	}
	p.Data = append(p.Data, typing)
	return 1
}
