package crawl

import "github.com/golang/glog"

type hub_parser struct {
	typing_parser
}

func minInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (p *Tdatas) ParseHub(base *Tdatas) bool {
	line := p.Segment.Line
	if base != nil {
		line = base.Hub.Line
	}
	p.Hub.drop_last_5_data()
	hasnew := false
	start := 0
	if l := len(p.Hub.Data); l > 0 {
		end := p.Hub.Data[l-1].ETime
		for i := len(line) - 1; i > -1; i-- {
			if end.Before(line[i].ETime) {
				continue
			}
			start = i
			break
		}
	}

	for i, l := start, len(line); i+2 < l; i++ {
		a := &line[i]
		b := &line[i+1]
		c := &line[i+2]
		minHigh, maxLow := a.High, a.Low
		minHigh = minInt(minHigh, b.High)
		minHigh = minInt(minHigh, c.High)
		maxLow = maxInt(maxLow, b.Low)
		maxLow = maxInt(maxLow, c.Low)
		if minHigh-maxLow < p.min_hub_height {
			continue
		}
		hub := *a
		hub.High = minHigh
		hub.Low = maxLow
		hub.end = c.end
		hub.ETime = c.ETime
		p.Hub.Data = append(p.Hub.Data, hub)
		i += 2
		hasnew = true
	}
	glog.Infoln("hub len(line)=", len(line), hasnew)
	return hasnew
}

func begin_end(line []Typing, begin, end int) (int, int) {
	l := len(line)
	if l < 1 {
		glog.Fatalf("segment line len=%d should > 0", l)
	}
	rv_begin := 0
	rv_end := 0
	for i := 0; i < l; i++ {
		if begin == line[i].begin {
			rv_begin = i
		}
		if end == line[i].end {
			rv_end = i
		}
	}
	return rv_begin, rv_end
}

func GG(line []Typing, t Typing) int {
	l := len(line)
	begin, end := begin_end(line, t.begin, t.end)
	if begin < 0 || begin >= l {
		glog.Fatalf("segment line begin=%d should in range [0, %d)", begin, l)
	}
	if end < 0 || end >= l {
		glog.Fatalf("segment line end=%d should in range [0, %d)", end, l)
	}

	v := line[begin].High
	for i, end := begin+1, end+1; i < end; i++ {
		if v < line[i].High {
			v = line[i].High
		}
	}
	return v
}

func DD(line []Typing, t Typing) int {
	l := len(line)
	begin, end := begin_end(line, t.begin, t.end)
	if begin < 0 || begin >= l {
		glog.Fatalf("segment line begin=%d should in range [0, %d)", begin, l)
	}
	if end < 0 || end >= l {
		glog.Fatalf("segment line end=%d should in range [0, %d)", end, l)
	}

	v := line[begin].Low
	for i, end := begin+1, end+1; i < end; i++ {
		if v > line[i].Low {
			v = line[i].Low
		}
	}
	return v
}

func (p *Tdatas) LinkHub(next *Tdatas) {
	hub := p.Hub
	hub.drop_last_5_line()
	segline := p.Segment.Line
	start := 0
	ldata := len(hub.Data)
	line := hub.Line
	prev := Typing{}

	if l := len(line); l > 0 {
		end := line[l-1].end
		for i := ldata - 1; i > -1; i-- {
			if end < hub.Data[i].end {
				continue
			}
			start = i
			prev = line[i]
			break
		}
	}

	for i := start; i < ldata; i++ {
		t := hub.Data[i]
		t.High = GG(segline, t)
		t.Low = DD(segline, t)
		if prev.I == 0 {
			line = append(line, t)
		} else if LineContain(&prev, &t) {
			line = append(line, t)
		} else if prev.Type == UpTyping && prev.High < t.High {
			l := len(line)
			line[l-1].High = t.High
			line[l-1].end = t.end
			line[l-1].ETime = t.ETime
		} else if prev.Type == DownTyping && prev.Low > t.Low {
			l := len(line)
			line[l-1].Low = t.Low
			line[l-1].end = t.end
			line[l-1].ETime = t.ETime
		} else {
			line = append(line, t)
		}

		prev = t
	}

	hub.Line = line
	glog.Infoln("hub link len(line)=", len(line))
	if next != nil {
		next.Segment.Line = hub.Line
	}
}
