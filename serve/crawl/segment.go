package crawl

func (p *Tdatas) ParseLine() {
	start := 0
	if l := len(p.Line); l > 0 {
		start = p.Line[l-1].I
	}
	if start > 0 {
		for i := len(p.Typing) - 1; i > -1; i-- {
			if p.Typing[i].I == start {
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
	}
}

func (p *Tdatas) ParseSegment() {
}
