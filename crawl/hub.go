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
		start = p.Hub.Data[l-1].end + 1
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
		hub.begin = i
		hub.end = i + 2
		hub.ETime = c.ETime
		p.Hub.Data = append(p.Hub.Data, hub)
		i += 2
		hasnew = true
	}
	glog.Infoln("hub len(line)=", len(line), hasnew)
	return hasnew
}

func GG(line []Typing, t Typing) int {
	l := len(line)
	begin, end := t.begin, t.end
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
	begin, end := t.begin, t.end
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
		start = line[l-1].end + 1
		if line[l-1].end < ldata {
			prev = hub.Data[line[l-1].end]
		} else {
			glog.Fatalln("line[].end >= len(hub.Data)")
		}
	}

	for i := start; i < ldata; i++ {
		t := hub.Data[i]
		t.High = GG(segline, t)
		t.Low = DD(segline, t)
		t.begin = i
		t.end = i
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
