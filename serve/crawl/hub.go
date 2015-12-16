package crawl

import "log"

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
	hasnew := false
	start := 0
	if l := len(p.Hub.Data); l > 0 {
		end := p.Hub.Data[l-1].End
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
		hub.ETime = c.ETime
		p.Hub.Data = append(p.Hub.Data, hub)
		hasnew = true
	}
	log.Println("hub len(line)=", len(line), hasnew)
	return hasnew
}

func (p *hub_parser) Link() bool {
	hasnew := false
	start := 0
	ldata := len(p.Data)
	if l := len(p.Line); l > 0 {
		t := p.Line[l-1]
		for i := ldata - 1; i > -1; i-- {
			if p.Data[i].End == t.End {
				start = i
				break
			}
		}
	}

	typing := Typing{}
	for i := start; i < ldata; i++ {
		t := p.Data[i]
		if typing.I == 0 {
			typing = t
			continue
		}

		if LineContain(&typing, &t) {
			typing.Type = DullTyping
			typing.High = maxInt(typing.High, t.High)
			typing.Low = minInt(typing.Low, t.Low)
		} else if typing.High < t.High {
			typing.Type = UpTyping
			typing.High = t.High
		} else if typing.Low > t.Low {
			typing.Type = DownTyping
			typing.Low = t.Low
		} else {
			log.Println("found unkonw typing of hub", typing, t)
		}

		typing.End = t.End
		typing.ETime = t.ETime
		if l := len(p.Line); l > 0 && p.Line[l-1].Type == typing.Type {
			if typing.Type == DullTyping {
				p.Line[l-1].High = maxInt(typing.High, p.Line[l-1].High)
				p.Line[l-1].Low = minInt(typing.Low, p.Line[l-1].Low)
			} else if typing.Type == UpTyping {
				p.Line[l-1].High = typing.High
			} else {
				p.Line[l-1].Low = typing.Low
			}
		} else {
			p.Line = append(p.Line, typing)
		}
		typing = t
		hasnew = true
	}

	log.Println("hub link len(line)=", len(p.Line), hasnew)
	return hasnew
}
