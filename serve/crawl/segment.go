package crawl

func (p *Tdatas) ParseLine() bool {
	hasnew := false
	start := 0
	if l := len(p.Line); l > 0 {
		for i := len(p.Typing) - 1; i > -1; i-- {
			if p.Typing[i].I == p.Line[l-1].I {
				start = i + 1
				break
			}
		}
	}

	for i, c := start, len(p.Typing); i < c; i++ {
		if len(p.Line) > 0 && p.Line[len(p.Line)-1].Type == p.Typing[i].Type {
			continue
		}
		p.Line = append(p.Line, p.Typing[i])
		hasnew = true
	}

	return hasnew
}

func (p *Tdatas) ParseSegment() bool {
	return false
}
