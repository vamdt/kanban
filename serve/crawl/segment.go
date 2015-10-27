package crawl

import "log"

type segment_parser struct {
	Data          []Typing
	tp            []typing_parser_node
	r1st_is_case1 bool
	r2rd_is_case1 bool
}

func (p *segment_parser) add_typing(typing Typing, case1 bool) {
	p.Data = append(p.Data, typing)
	p.r2rd_is_case1 = p.r1st_is_case1
	p.r1st_is_case1 = case1
	log.Println("new segment typing", typing.Type, p.r2rd_is_case1, p.r1st_is_case1, len(p.Data))
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
	log.Println("new node len(tp)", len(p.tp), "line:", i, "len(data):", len(p.Data))
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
		for i := li - 2; i > si; i-- {
			if p.Typing.Line[i].High >= p.Typing.Line[li].Low {
				return true
			}
		}
	} else if p.Typing.Line[si].Type == DownTyping {
		for i := li - 2; i > si; i-- {
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

	if false && len(p.Segment.tp) < 1 && len(p.Segment.Data) < 1 {
		for i := start; i < l; i++ {
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
    start = start + 1
  }

	log.Println("start", start)
	for i := start; i < l; i += 2 {

		if len(p.Segment.tp) < 1 {
			if i+2 > l {
				return hasnew
			}

			p.Segment.new_node(i, &p.Typing)
			continue
		}

		prev := &p.Segment.tp[len(p.Segment.tp)-1]

		if dlen := len(p.Segment.Data); dlen > 1 && !p.Segment.r1st_is_case1 {
			a := p.Segment.Data[dlen-1]
			b := p.Segment.tp[len(p.Segment.tp)-1].t
			not_complete := false
			if a.Type == TopTyping && b.High > a.High {
				not_complete = true
			}
			if a.Type == BottomTyping && b.Low < a.Low {
				not_complete = true
				log.Println("found a V VV")
			}

			if not_complete {
				start := 0
				if dlen > 1 {
					start = p.Segment.Data[dlen-2].I
				} else {
					start = p.Segment.Data[dlen-1].I - 3
				}
				p.feat_normalized(start, i)
				p.Segment.Data = p.Segment.Data[:dlen-1]
				i--
				hasnew = true
				continue
			}
		}

		if dlen := len(p.Segment.Data); dlen > 1 && !p.Segment.r2rd_is_case1 {
			a := p.Segment.Data[dlen-2]
			b := p.Segment.tp[len(p.Segment.tp)-1].t
			not_complete := false
			if a.Type == TopTyping && b.High > a.High {
				not_complete = true
			}
			if a.Type == BottomTyping && b.Low < a.Low {
				not_complete = true
			}

			if not_complete {
				start := 0
				if dlen > 1 {
					start = p.Segment.Data[dlen-2].I
				} else {
					start = p.Segment.Data[dlen-1].I - 3
				}
				p.feat_normalized(start, i)
				p.Segment.Data = p.Segment.Data[:dlen-2]
				continue
			}
		}

		a := &Tdata{}
		a.High = p.Typing.Line[i].High
		a.Low = p.Typing.Line[i].Low
		a.Time = p.Typing.Line[i].Time

		if Contain(&prev.d, a) {
			if len(p.Segment.tp) > 1 {
				firstIsBreak := false
				if len(p.Segment.tp) > 2 {
					pprev := &p.Segment.tp[len(p.Segment.tp)-2]
					if p.IsLineBreakSegment(pprev.t.I) {
						firstIsBreak = true
					}
				}
				if !firstIsBreak {
					if p.IsLineBreakSegment(prev.t.I) {
						if len(p.Segment.tp) > 2 && !hasGap(&prev.d, &p.Segment.tp[len(p.Segment.tp)-2].d) {
							p.Segment.new_node(i, &p.Typing)
							continue
						}
					} else if p.IsLineBreakSegment(i) {
						p.Segment.new_node(i, &p.Typing)
						continue
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
			prev.t.End = i
		} else {
			if len(p.Segment.tp) > 2 {
				pprev := &p.Segment.tp[len(p.Segment.tp)-2]
				typing := pprev.t
				if p.IsLineBreakSegment(pprev.t.I) && !hasGap(&pprev.d, &p.Segment.tp[len(p.Segment.tp)-3].d) {
					if Contain(&pprev.d, &prev.d) {
						need_wait := true
						case1_seg_ok := false
						for j := i; j < l; j += 2 {
							high := p.Typing.Line[j].High
							low := p.Typing.Line[j].Low
							if pprev.t.Type == DownTyping {
								if high > pprev.d.High {
									// not segment
									need_wait = false
									break
								} else if low < pprev.d.Low {
									// is segment
									need_wait = false

									typing.Type = TopTyping
									typing.High = pprev.d.High
									typing.Low = pprev.d.Low
									typing.Price = typing.High
									typing.Time = pprev.d.Time
									p.Segment.add_typing(typing, true)
									i = pprev.t.I - 1
									p.Segment.clear()
									hasnew = true
									case1_seg_ok = true
									break
								}
							} else if pprev.t.Type == UpTyping {
								if low < pprev.d.Low {
									// not segment
									need_wait = false
									break
								} else if high > pprev.d.High {
									// is segment
									need_wait = false

									typing.Type = BottomTyping
									typing.High = pprev.d.High
									typing.Low = pprev.d.Low
									typing.Price = pprev.d.Low
									typing.Time = pprev.d.Time
									p.Segment.add_typing(typing, true)
									i = pprev.t.I - 1
									p.Segment.clear()
									hasnew = true
									case1_seg_ok = true
									break
								}
							}
						}
						if need_wait {
							p.Segment.clear()
							return hasnew
						}
						if case1_seg_ok {
							continue
						}
					}
				}
			}

			if len(p.Segment.tp) < 2 && p.IsLineBreakSegment(i) {
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
		if p.Segment.parse_top_bottom() {
			hasnew = true
			i = p.Segment.tp[len(p.Segment.tp)-2].t.I - 1
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
