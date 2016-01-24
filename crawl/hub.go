package crawl

import "github.com/golang/glog"

type hub_parser struct {
	typing_parser
}

func (p *Tdatas) ParseHub() {
	line := p.Segment.Line
	p.Hub.drop_last_5_data()
	start := 0
	if l := len(p.Hub.Data); l > 0 {
		start = p.Hub.Data[l-1].end + 1
	} else {
		start = findLineDir(line, len(line))
	}
	if start == -1 {
		return
	}
	glog.Infoln(p.tag, "for hub start", start)

	for i, l := start, len(line); i+2 < l; i++ {
		lhub := len(p.Hub.Data)
		if lhub > 0 {
			hub := &p.Hub.Data[lhub-1]
			if hub.end == i {
				gn := g(line[i+2])
				dn := d(line[i+2])
				zg := hub.High
				zd := hub.Low

				// [dn, gn] # [ZD, ZG]
				if dn > zg || gn < zd {
				} else {
					hub.end = i + 2
					hub.ETime = line[i+2].ETime
					i++
				}
				continue
			} else if hub.end > i {
				continue
			}
		}
		zg := ZG(line[i], line[i+1], line[i+2])
		zd := ZD(line[i], line[i+1], line[i+2])
		if zg-zd < p.min_hub_height {
			i++
			continue
		}
		hub := line[i]
		hub.High = zg
		hub.Low = zd
		hub.begin = i
		hub.end = i + 2
		hub.ETime = line[i+2].ETime
		p.Hub.Data = append(p.Hub.Data, hub)
		i++
	}
}

func g(a Typing) int { return a.High }

func d(a Typing) int { return a.Low }

// ZG=min(g1, g2)
func ZG(a, b, c Typing) int {
	return minInt(a.High, c.High)
}

// ZD=max(d1, d2)
func ZD(a, b, c Typing) int {
	return maxInt(a.Low, c.Low)
}

// G=min(gn)
func G(line []Typing, t Typing) int {
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
		if v > line[i].High {
			v = line[i].High
		}
	}
	return v
}

// GG=max(gn)
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

// D=max(dn)
func D(line []Typing, t Typing) int {
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
		if v < line[i].Low {
			v = line[i].Low
		}
	}
	return v
}

// DD=min(dn)
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
	var prev *Typing

	if l := len(line); l > 0 {
		start = line[l-1].end + 1
		prev = &line[l-1]
	}

	for i := start; i < ldata; i++ {
		t := hub.Data[i]
		t.High = GG(segline, t)
		t.Low = DD(segline, t)
		t.begin = i
		t.end = i
		if prev == nil {
			line = append(line, t)
		} else if t.Low > prev.High {
			// UpTyping
			l := len(line)
			line[l-1].High = t.High
			line[l-1].end = t.end
			line[l-1].ETime = t.ETime
			line[l-1].Type = UpTyping
			t = line[l-1]
		} else if t.High < prev.Low {
			// DownTyping
			l := len(line)
			line[l-1].Low = t.Low
			line[l-1].end = t.end
			line[l-1].ETime = t.ETime
			line[l-1].Type = DownTyping
			t = line[l-1]
		} else {
			line = append(line, t)
		}

		prev = &t
	}

	hub.Line = line
	glog.Infoln(hub.tag, "hub link len(line)=", len(line))
	if next != nil {
		next.Segment.Line = hub.Line
	}
}
