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
		end := p.Hub[l-1].End
		for i := len(line) - 1; i > -1; i-- {
			if end == line[i].End {
				start = i
				break
			}
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
		hub.End = c.End
		p.Hub = append(p.Hub, hub)
	}
	return hasnew
}
