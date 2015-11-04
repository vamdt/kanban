package crawl

import "log"

type segment_parser struct {
	Data []Typing
	tp   []typing_parser_node

	break_index   int
	unsure_typing Typing
	need_sure     bool
	wait_3end     bool
}

func (p *segment_parser) add_typing(typing Typing, case1 bool) bool {
	typing.Case1 = case1
	if p.need_sure {
		p.need_sure = false
		if l := len(p.Data); l > 0 && p.unsure_typing.I == p.Data[l-1].I {
			log.Println("overwrite prev case2 segment", p.unsure_typing, len(p.Data))
			p.Data[l-1] = p.unsure_typing
		} else {
			log.Println("new ensure case2 segment", p.unsure_typing, len(p.Data))
			p.Data = append(p.Data, p.unsure_typing)
		}
	}

	if !case1 {
		p.wait_3end = true
		p.need_sure = true
		p.unsure_typing = typing
		log.Println("new unsure case2 segment", typing, len(p.Data))
		return true
	}

	p.Data = append(p.Data, typing)
	p.wait_3end = true
  p.need_sure = false
	log.Println("new case1 segment typing", typing.Type, len(p.Data))
	return true
}

func (p *segment_parser) is_unsure_typing_fail(a *Tdata) bool {
	if !p.need_sure {
		return false
	}

	switch p.unsure_typing.Type {
	case BottomTyping:
		if p.unsure_typing.Low > a.Low {
			return true
		}
	case TopTyping:
		if p.unsure_typing.High < a.High {
			return true
		}
	}

	return false
}

func (p *segment_parser) clean_fail_unsure_typing() int {
	start := p.unsure_typing.I
	if l := len(p.Data); l > 0 && !p.Data[l-1].Case1 {
		p.unsure_typing = p.Data[l-1]
		p.need_sure = true
	} else {
		p.need_sure = false
	}
	return start
}

func (p *segment_parser) new_node(i int, ptyping *typing_parser, isbreak bool) {
	if len(p.tp) > 0 {
		p.tp[len(p.tp)-1].t.End = i - 2
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
	if isbreak {
		p.break_index = len(p.tp) - 1
	} else {
		p.break_index--
	}
	log.Println("new node len(tp)", len(p.tp), "line:", i, "len(data):", len(p.Data), "bindex", p.break_index, isbreak)
}

func (p *segment_parser) clear() {
	p.tp = []typing_parser_node{}
	p.break_index = -1
}

func (p *segment_parser) clean() {
	if len(p.tp) > 3 {
		var tmp []typing_parser_node
		tmp = append(tmp, p.tp[len(p.tp)-3:]...)
		p.tp = tmp
	}
}

func (p *segment_parser) isLineBreak(line, nline *Typing) bool {
	ltp := len(p.tp)
	if ltp < 1 {
		return false
	}

	if p.tp[0].t.Type == UpTyping {
		if p.tp[ltp-1].d.Low < line.High {
			if LineContain(&p.tp[ltp-1].t, line) {
				if p.break_index != ltp-1 {
					if nline.Low < line.Low {
						return false
					}
				}
			}
			return true
		}
	} else if p.tp[0].t.Type == DownTyping {
		if p.tp[ltp-1].d.High > line.Low {
			if LineContain(&p.tp[ltp-1].t, line) {
				if p.break_index != ltp-1 {
					if nline.High > line.High {
						return false
					}
				}
			}
			return true
		}
	}
	return false
}

func (p *segment_parser) handle_special_case1(i int, a *Tdata) (bool, int) {
	ltp := len(p.tp)
	if ltp < 2 {
		return false, i
	}
	//        |       |
	// check ||| or |||
	//        ||     ||
	//         |     |
	if p.break_index != ltp-1 {
		return false, i
	}
	pprev := &p.tp[ltp-2]
	prev := &p.tp[ltp-1]
	if !Contain(&pprev.d, &prev.d) {
		return false, i
	}

	typing := prev.t
	typing.High = prev.d.High
	typing.Low = prev.d.Low
	typing.Time = prev.d.Time
	case1_seg_ok := false
	if prev.t.Type == DownTyping && prev.d.Low > a.Low && prev.d.High > a.High {
		// TopTyping yes
		case1_seg_ok = true
		typing.Type = TopTyping
		typing.Price = typing.High
	} else if prev.t.Type == UpTyping && prev.d.High < a.High && prev.d.Low < a.Low {
		// BottomTyping yes
		case1_seg_ok = true
		typing.Type = BottomTyping
		typing.Price = typing.Low
	}

	if case1_seg_ok {
		typing.End = i - 2
		prev.t.End = typing.End
		p.add_typing(typing, true)
		p.clear()
	}
	return case1_seg_ok, typing.End - 1
}

func need_skip_line(prev *typing_parser_node, a *Tdata) bool {
	if prev.t.Type == DownTyping && a.High < prev.d.High && a.Low < prev.d.Low {
		return true
	}
	if prev.t.Type == UpTyping && a.Low > prev.d.Low && a.High > prev.d.High {
		return true
	}
	return false
}

func (p *Tdatas) need_wait_3end(i int, a *Tdata) (bool, int) {
	ltp := len(p.Segment.tp)
	if ltp < 3 || !p.Segment.wait_3end {
		return false, i
	}
	prev := &p.Segment.tp[ltp-1]
	pprev := &p.Segment.tp[ltp-2]
	if Contain(&prev.d, a) {
		if prev.t.Type == UpTyping {
			if pprev.d.High < a.High {
				a = DownContainMerge(&prev.d, a)
				if prev.d.Low != a.Low {
					prev.t.I = i
				}
				prev.d = *a
				prev.t.End = i
				return true, i
			}
		} else if prev.t.Type == DownTyping {
			if pprev.d.Low > a.Low {
				a = UpContainMerge(&prev.d, a)
				if prev.d.High != a.High {
					prev.t.I = i
				}
				prev.d = *a
				prev.t.End = i
				return true, i
			}
		}
	}
	p.Segment.clear()
	i = pprev.t.End + 1
	p.Segment.new_node(i, &p.Typing, false)
	i = prev.t.End - 1
	p.Segment.wait_3end = false
	return false, i
}

func (p *Tdatas) ParseSegment() bool {
	hasnew := false
	start := 0

	l := len(p.Typing.Line)
	if l > 0 && p.Typing.Line[l-1].Type != UpTyping && p.Typing.Line[l-1].Type != DownTyping {
		l--
	}

	if x := len(p.Segment.tp); x > 0 {
		start = p.Segment.tp[x-1].t.End + 1
	} else if y := len(p.Segment.Data); y > 0 {
		start = p.Segment.Data[y-1].End + 1
	} else if l > 100 {
		for i := 0; i < l; i++ {
			if i+2 > l {
				return hasnew
			}

			if p.Typing.Line[i].Type == UpTyping && p.Typing.Line[i+2].High > p.Typing.Line[i].High && p.Typing.Line[i+2].Low > p.Typing.Line[i].Low {
				// Up yes
				start = i + 1
				break
			} else if p.Typing.Line[i].Type == DownTyping && p.Typing.Line[i+2].High < p.Typing.Line[i].High && p.Typing.Line[i+2].Low < p.Typing.Line[i].Low {
				// Down yes
				start = i + 1
				break
			}
		}
	} else {
		start = 1
	}

	log.Println("start", start)
	ltp := len(p.Segment.tp)
	for i := start; i < l; i += 2 {

		ltp = len(p.Segment.tp)
		if ltp < 1 {
			p.Segment.new_node(i, &p.Typing, false)
			continue
		}

		prev := &p.Segment.tp[ltp-1]

		a := &Tdata{}
		a.High = p.Typing.Line[i].High
		a.Low = p.Typing.Line[i].Low
		a.Time = p.Typing.Line[i].Time

		if p.Segment.need_sure && p.Segment.is_unsure_typing_fail(a) {
			log.Println("found unsure typing fail", i, p.Segment.unsure_typing, a)
			i = p.Segment.clean_fail_unsure_typing()
			log.Println("new start", i, "need_sure", p.Segment.need_sure)
			p.Segment.clear()
			continue
		}

		//if ok, j := p.need_wait_3end(i, a); ok {
		//i = j
		//continue
		//}

		if Contain(&prev.d, a) {
			if !p.Segment.need_sure {
				if p.Segment.break_index == ltp-1 {
					// case   |  or |
					//       ||     |||
					//      |||      ||
					//      |         |
					if prev.t.Type == DownTyping && prev.d.High < a.High {
						p.Segment.new_node(i, &p.Typing, false)
						continue
					}
					if prev.t.Type == UpTyping && prev.d.Low > a.Low {
						p.Segment.new_node(i, &p.Typing, false)
						continue
					}
				} else if p.Segment.break_index < 0 {
					if i+1 < l {
						//       |
						// case |||
						//       |
						if p.Segment.isLineBreak(&p.Typing.Line[i], &p.Typing.Line[i+1]) {
							p.Segment.new_node(i, &p.Typing, true)
							continue
						}
					} else {
						return hasnew
					}
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
			prev.t.High = prev.d.High
			prev.t.Low = prev.d.Low
			prev.t.End = i
		} else {
			if ltp > 1 {
				if ok, j := p.Segment.handle_special_case1(i, a); ok {
					i = j
					hasnew = true
					continue
				}
			}

			if ltp < 2 {
				if need_skip_line(prev, a) {
					continue
				}
			}
			isbreak := false
			if p.Segment.break_index < 0 {
				isbreak = p.Segment.isLineBreak(&p.Typing.Line[i], nil)
			}
			p.Segment.new_node(i, &p.Typing, isbreak)
		}

		p.Segment.clean()
		if p.Segment.parse_top_bottom() {
			hasnew = true
			i = p.Segment.tp[len(p.Segment.tp)-2].t.End - 1
			p.Segment.clear()
		}
	}
	return hasnew
}

func hasGap(a, b *Tdata) bool {
	return a.Low > b.High || a.High < b.Low
}

func (p *segment_parser) parse_top_bottom() bool {
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

	dlen := len(p.Data)
	if dlen > 0 {
		if typing.Type == TopTyping && p.Data[dlen-1].Type == BottomTyping {
			if typing.High <= p.Data[dlen-1].High {
				log.Println("find a bottom high then top")
			}
		}
	}

	p.add_typing(typing, !hasGap(a, b))
	return true
}
