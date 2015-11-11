package crawl

func (p *Tdatas) ParseHub(base *Tdatas) {
	hasnew := false
	start := 0
	if l := len(p.Hub); l > 0 {
		start = p.Hub[l-1].I + 1
	}

	for i, l := start, len(p.Segment.Line); i < l; i++ {
	}
	return hasnew
}
