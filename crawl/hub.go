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

func (p *Tdatas) ParseHub() bool {
	line := p.Segment.Line
	p.Hub.drop_last_5_data()
	hasnew := false
	start := 0
	if l := len(p.Hub.Data); l > 0 {
		start = p.Hub.Data[l-1].end + 1
	}

	for i, l := start, len(line); i+2 < l; i++ {
		zg := ZG(line[i], line[i+1], line[i+2])
		zd := ZD(line[i], line[i+1], line[i+2])
		if zg-zd < p.min_hub_height {
			continue
		}
		hub := line[i]
		hub.High = zg
		hub.Low = zd
		hub.begin = i
		hub.end = i + 2
		hub.ETime = line[i+2].ETime
		p.Hub.Data = append(p.Hub.Data, hub)
		i += 2
		hasnew = true
	}
	glog.Infoln(p.Hub.tag, "hub len(line)=", len(line), hasnew)
	return hasnew
}

func ZG(a, b, c Typing) int {
	return minInt(a.High, c.High)
}

func ZD(a, b, c Typing) int {
	return maxInt(a.Low, c.Low)
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
		} else if t.Low > prev.High {
			// UpTyping
			l := len(line)
			line[l-1].High = t.High
			line[l-1].end = t.end
			line[l-1].ETime = t.ETime
		} else if t.High < prev.Low {
			// DownTyping
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
	glog.Infoln(hub.tag, "hub link len(line)=", len(line))
	if next != nil {
		next.Segment.Line = hub.Line
	}
}
