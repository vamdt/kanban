package crawl

func (p *Tdatas) Macd() {
	macd(p, 12, 26, 9)
}

func macd(td *Tdatas, short, long, m int) {
	c := len(td.Data)
	if c < 1 {
		return
	}
	td.Data[0].emas = 0
	td.Data[0].emal = 0
	td.Data[0].DIFF = 0
	td.Data[0].DEA = 0
	td.Data[0].MACD = 0
	s1, s2 := short-1, short+1
	l1, l2 := long-1, long+1
	m1, m2 := m-1, m+1
	if c > 1 {
		td.Data[1].emas = (td.Data[0].Close*s1 + td.Data[1].Close*200) / s2
		td.Data[1].emal = (td.Data[0].Close*l1 + td.Data[1].Close*200) / l2
		td.Data[1].DIFF = td.Data[1].emas - td.Data[1].emal
		td.Data[1].DEA = (td.Data[0].DEA*m1 + td.Data[1].DIFF*2) / m2
		td.Data[1].MACD = 2 * (td.Data[1].DIFF - td.Data[1].DEA)
	}
	for i := 2; i < c; i++ {
		td.Data[i].emas = (td.Data[i-1].emas*s1 + td.Data[i].Close*200) / s2
		td.Data[i].emal = (td.Data[i-1].emal*l1 + td.Data[i].Close*200) / l2
		td.Data[i].DIFF = td.Data[i].emas - td.Data[i].emal
		td.Data[i].DEA = (td.Data[i-1].DEA*m1 + td.Data[i].DIFF*2) / m2
		td.Data[i].MACD = 2 * (td.Data[i].DIFF - td.Data[i].DEA)
	}
}
