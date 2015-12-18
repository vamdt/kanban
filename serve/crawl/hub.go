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
		hasnew = true
	}
	glog.Infoln("hub len(line)=", len(line), hasnew)
	return hasnew
}

func (p *hub_parser) Link() bool {
	hasnew := false
	start := 0
	ldata := len(p.Data)
	if l := len(p.Line); l > 0 {
		t := p.Line[l-1]
		for i := ldata - 1; i > -1; i-- {
			if p.Data[i].ETime.Equal(t.ETime) {
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
			glog.Infoln("found unkonw typing of hub", typing, t)
		}

		typing.end = t.end
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

	glog.Infoln("hub link len(line)=", len(p.Line), hasnew)
	return hasnew
}
