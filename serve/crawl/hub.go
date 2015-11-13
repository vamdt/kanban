package crawl

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
		line = base.Hub
	}
	hasnew := false
	start := 0
	if l := len(p.Hub); l > 0 {
		start = p.Hub[l-1].I + 1
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
		hub := Typing{High: minHigh, Low: maxLow}
		hub.I = i + 2
		hub.begin = a.begin
		hub.End = c.End
		p.Hub = append(p.Hub, hub)
	}
	return hasnew
}
