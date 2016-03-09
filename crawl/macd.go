package crawl

func (p *Tdatas) Macd(start int) {
	macd(p, 12, 26, 9, start)
}

func macd(td *Tdatas, short, long, m, start int) {
	c := len(td.Data)
	if c < 1 {
		return
	}

	if start > c {
		start = c
	}
	for i := start - 1; i > -1; i-- {
		if td.Data[i].Emas > 0 {
			break
		}
		start = i
	}

	if start < 1 {
		td.Data[0].Emas = td.Data[0].Close * 10
		td.Data[0].Emal = td.Data[0].Close * 10
		td.Data[0].DIFF = 0
		td.Data[0].DEA = 0
		td.Data[0].MACD = 0
		start = 1
	}

	s1, s2 := short-1, short+1
	l1, l2 := long-1, long+1
	m1, m2 := m-1, m+1
	for i := start; i < c; i++ {
		td.Data[i].Emas = (td.Data[i-1].Emas*s1 + td.Data[i].Close*20) / s2
		td.Data[i].Emal = (td.Data[i-1].Emal*l1 + td.Data[i].Close*20) / l2
		td.Data[i].DIFF = td.Data[i].Emas - td.Data[i].Emal
		td.Data[i].DEA = (td.Data[i-1].DEA*m1 + td.Data[i].DIFF*2) / m2
		td.Data[i].MACD = 2 * (td.Data[i].DIFF - td.Data[i].DEA)
	}
}
