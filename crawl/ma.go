package crawl

func (p *Tdatas) Factor() int {
	ma := []int{5, 13, 21, 34, 55, 89, 144, 233}
	l := len(p.Data)
	if l < 1 {
		return 0
	}

	factor := len(ma) + 1
	last := p.Data[l-1].Close
	for i, m := range ma {
		if last <= p.Ma(m) {
			factor = i + 1
			break
		}
	}
	return factor
}

func (p *Tdatas) Ma(m int) int {
	v := 0

	end := len(p.Data)
	if end < 1 {
		return v
	}

	start := end - m
	if start < 0 {
		start = 0
	}

	m = end - start
	if m < 1 {
		m = 1
	}

	for i := start; i < end; i++ {
		v += p.Data[i].Close
	}
	return v / m
}
