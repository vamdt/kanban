package crawl

import "github.com/golang/glog"

type segment_parser struct {
	typing_parser

	break_index int
	wait_3end   bool
}

func (p *segment_parser) need_sure() bool {
	l := len(p.Data)
	if l < 1 {
		return false
	}

	return !p.Data[l-1].Case1
}

func (p *segment_parser) add_typing(typing Typing, case1 bool) {
	l := len(p.Data)
	typing.Case1 = case1

	p.Data = append(p.Data, typing)
	p.wait_3end = true
	glog.V(SegmentD).Infof("new segment typing [%d] case1[%t] %+v", l, case1, typing)
}

func (p *segment_parser) is_unsure_typing_fail(a *Tdata) bool {
	l := len(p.Data)
	if l < 1 {
		return false
	}

	if !p.need_sure() {
		return false
	}

	t := p.Data[l-1]
	switch t.Type {
	case BottomTyping:
		if t.Low > a.Low {
			return true
		}
	case TopTyping:
		if t.High < a.High {
			return true
		}
	}

	return false
}

func (p *segment_parser) clean_fail_unsure_typing() int {
	l := len(p.Data)
	if l < 1 {
		panic("should not be here, len(Data) < 1")
	}
	start := p.Data[l-1].i
	p.Data = p.Data[:l-1]
	return start
}

func (p *segment_parser) new_node(i int, ptyping *typing_parser, isbreak bool) {
	line := ptyping.Line
	tp := typing_parser_node{}
	tp.t = line[i]
	tp.t.begin = i
	tp.t.i = i
	tp.t.end = i
	tp.t.assertETimeMatchEndLine(line, "new_node")

	if l := len(p.tp); l > 0 {
		p.tp[l-1].t.end = i - 2
		p.tp[l-1].t.ETime = line[i-2].ETime
		p.tp[l-1].t.assertETimeMatchEndLine(line, "new node prev")
	} else if l = len(p.Data); l > 0 && p.Data[l-1].end == i-1 {
		if t := p.Data[l-1]; t.begin < t.end {
			switch tp.t.Type {
			case UpTyping:
				tp.t.Low = minInt(t.Low, tp.t.High)
			case DownTyping:
				tp.t.High = maxInt(t.High, tp.t.Low)
			}
		}
	}

	tp.d.Time = tp.t.Time
	tp.d.High = tp.t.High
	tp.d.Low = tp.t.Low
	p.tp = append(p.tp, tp)
	if isbreak {
		p.break_index = len(p.tp) - 1
	} else {
		p.break_index--
	}
	glog.V(SegmentD).Infoln("new node len(tp)", len(p.tp), "line:", i, "len(data):", len(p.Data), "bindex", p.break_index, isbreak)
}

func (p *segment_parser) reset() {
	p.tp = []typing_parser_node{}
	p.break_index = -1
}

func (p *segment_parser) clean() {
	if len(p.tp) > 3 {
		p.tp = p.tp[len(p.tp)-3:]
	}
}

func (p *segment_parser) isLineBreak(lines []Typing, i int) bool {
	l := len(lines)
	if i+1 >= l {
		return false
	}
	ltp := len(p.tp)
	if ltp < 1 {
		return false
	}

	line, nline := lines[i], lines[i+1]
	if p.tp[0].t.Type == UpTyping {
		if line.High > p.tp[ltp-1].d.Low {
			return line.Low < p.tp[ltp-1].d.Low && line.Low < nline.Low
		}
	} else if p.tp[0].t.Type == DownTyping {
		if line.Low < p.tp[ltp-1].d.High {
			return line.High > nline.High && line.High > p.tp[ltp-1].d.High
		}
	}
	return false
}

func (p *segment_parser) handle_special_case1(i int, a *Tdata) bool {
	ltp := len(p.tp)
	if ltp < 2 {
		return false
	}
	//        |       |
	// check ||| or |||
	//        ||     ||
	//         |     |
	if p.break_index != ltp-1 {
		return false
	}
	pprev := &p.tp[ltp-2]
	prev := &p.tp[ltp-1]
	if !Contain(&pprev.d, &prev.d) {
		return false
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
		p.add_typing(typing, true)
		p.reset()
	}
	return case1_seg_ok
}

func merge_contain_node(prev *typing_parser_node, a *Tdata, i int, line *Typing) {
	if prev.t.Type == UpTyping {
		a = DownContainMerge(&prev.d, a)
		if prev.d.Low != a.Low {
			prev.t.i = i
		}
	} else {
		if prev.t.Type != DownTyping {
			glog.Fatalln("prev should be a DownTyping line %+v", prev)
		}
		a = UpContainMerge(&prev.d, a)
		if prev.d.High != a.High {
			prev.t.i = i
		}
	}
	glog.V(SegmentD).Infof("merge prev t %+v with line[%d] %+v", prev.t, i, line)
	prev.d = *a
	prev.t.High = prev.d.High
	prev.t.Low = prev.d.Low
	prev.t.end = i
	prev.t.ETime = line.ETime
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

func (p *Tdatas) need_wait_3end(i int, a *Tdata, line []Typing) (bool, int) {
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
					prev.t.i = i
				}
				prev.d = *a
				prev.t.end = i
				prev.t.ETime = line[i].ETime
				prev.t.assertETimeMatchEndLine(line, "need_wait_3end")
				return true, i
			}
		} else if prev.t.Type == DownTyping {
			if pprev.d.Low > a.Low {
				a = UpContainMerge(&prev.d, a)
				if prev.d.High != a.High {
					prev.t.i = i
				}
				prev.d = *a
				prev.t.end = i
				prev.t.ETime = line[i].ETime
				prev.t.assertETimeMatchEndLine(line, "need_wait_3end")
				return true, i
			}
		}
	}
	p.Segment.reset()

	i = pprev.t.end + 1
	i = 1 + pprev.t.assertETimeMatchEndLine(line, "need_wait_3end")
	p.Segment.new_node(i, &p.Typing, false)

	i = prev.t.end - 1
	p.Segment.wait_3end = false
	return false, i
}

func isLineUp(line []Typing, length, begin int) bool {
	i, l := begin, length
	if i+2 < l && line[i+2].High > line[i].High && line[i+2].Low > line[i].Low {
		return true
	}

	if i+4 < l && line[i+4].High > line[i].High && line[i+4].Low > line[i].Low {
		return true
	}
	return false
}

func isLineDown(line []Typing, length, begin int) bool {
	i, l := begin, length
	if i+2 < l && line[i+2].High < line[i].High && line[i+2].Low < line[i].Low {
		return true
	}
	if i+4 < l && line[i+4].High < line[i].High && line[i+4].Low < line[i].Low {
		return true
	}
	return false
}

func findLineDir(line []Typing, l int) int {
	for i := 0; i < l; i++ {
		if line[i].Type == UpTyping {
			// Up yes
			if isLineUp(line, l, i) {
				return i + 1
			}
		} else if line[i].Type == DownTyping {
			// Down yes
			if isLineDown(line, l, i) {
				return i + 1
			}
		}
	}
	return -1
}

func (p *Tdatas) ParseSegment() bool {
	hasnew := false
	start := 0

	p.Segment.drop_last_5_data()

	l := len(p.Typing.Line)
	if l > 0 && p.Typing.Line[l-1].Type != UpTyping && p.Typing.Line[l-1].Type != DownTyping {
		l--
	}

	if y := len(p.Segment.Data); y > 0 {
		start = p.Segment.Data[y-1].end + 1
		p.Segment.Data[y-1].assertETimeMatchEndLine(p.Typing.Line, "ParseSegment start2")
	} else {
		i := findLineDir(p.Typing.Line, l)
		if i == -1 {
			return hasnew
		}
		start = i
	}

	glog.V(SegmentV).Infof("start-%d lines-%d", start, l)
	p.Segment.reset()
	for i := start; i < l; i += 2 {

		ltp := len(p.Segment.tp)
		if ltp < 1 {
			p.Segment.new_node(i, &p.Typing, false)
			continue
		}

		prev := &p.Segment.tp[ltp-1]

		a := &Tdata{}
		a.High = p.Typing.Line[i].High
		a.Low = p.Typing.Line[i].Low
		a.Time = p.Typing.Line[i].Time

		if p.Segment.need_sure() && p.Segment.is_unsure_typing_fail(a) {
			glog.V(SegmentV).Infoln("found unsure typing fail", i, a, p.Segment.Data[len(p.Segment.Data)-1])
			i = p.Segment.clean_fail_unsure_typing() - 2
			glog.V(SegmentV).Infoln("new start", i, "need_sure", p.Segment.need_sure)
			p.Segment.reset()
			continue
		}

		//if ok, j := p.need_wait_3end(i, a, p.Typing.Line); ok {
		//i = j
		//continue
		//}

		if Contain(&prev.d, a) {
			if !p.Segment.need_sure() {
				if p.Segment.break_index == ltp-1 {
					// case   |  or |
					//       ||     |||
					//      |||      ||
					//      |         |
					if prev.t.Type == DownTyping && prev.d.High < a.High {
						isBreak := p.Segment.isLineBreak(p.Typing.Line, i)
						p.Segment.new_node(i, &p.Typing, isBreak)
						continue
					}
					if prev.t.Type == UpTyping && prev.d.Low > a.Low {
						isBreak := p.Segment.isLineBreak(p.Typing.Line, i)
						p.Segment.new_node(i, &p.Typing, isBreak)
						continue
					}
				} else if p.Segment.break_index < 0 {
					isBreak := p.Segment.isLineBreak(p.Typing.Line, i)
					if isBreak {
						//       |
						// case |||
						//       |
						p.Segment.new_node(i, &p.Typing, true)
						continue
					}
				}
			}

			merge_contain_node(prev, a, i, &p.Typing.Line[i])
			prev.t.assertETimeMatchEndLine(p.Typing.Line, "ParseSegment Contain")
			continue
		} else {
			if ltp > 1 {
				if ok := p.Segment.handle_special_case1(i, a); ok {
					i = i - 2 - 1
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
				isbreak = p.Segment.isLineBreak(p.Typing.Line, i)
			}
			p.Segment.new_node(i, &p.Typing, isbreak)
		}

		p.Segment.clean()
		if p.Segment.parse_top_bottom() {
			hasnew = true
			i = p.Segment.tp[len(p.Segment.tp)-2].t.end - 1
			p.Segment.reset()
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
				glog.Infoln("find a bottom high then top")
			}
		}
	}

	p.add_typing(typing, !hasGap(a, b))
	return true
}
