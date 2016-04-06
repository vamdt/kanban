package crawl

import . "./base"

func (p *Tdatas) Macd(start int, prev0 *Tdata) {
	p.macd(12, 26, 9, start, prev0)
}

func (td *Tdatas) macd(short, long, m, start int, prev0 *Tdata) {
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

	s1, s2 := short-1, short+1
	l1, l2 := long-1, long+1
	m1, m2 := m-1, m+1

	if start == 0 && td.Data[0].Emas == 0 && prev0 != nil && prev0.Emas != 0 {
		for i := start; i < c; i++ {
			start = i
			td.Data[i].Emas = prev0.Emas
			td.Data[i].Emal = prev0.Emal
			td.Data[i].DEA = prev0.DEA

			if prev0.Time.Equal(td.Data[i].Time) {
				break
			} else if prev0.Time.Before(td.Data[i].Time) {
				td.Data[i].Emas = (prev0.Emas*s1 + td.Data[i].Close*20) / s2
				td.Data[i].Emal = (prev0.Emal*l1 + td.Data[i].Close*20) / l2
				td.Data[i].DEA = (prev0.DEA*m1 + td.Data[i].DIFF*2) / m2
				break
			}
		}
		td.Data[start].UpdateMacd()
		start = start + 1
	}

	if start < 1 {
		td.Data[0].Emas = td.Data[0].Close * 10
		td.Data[0].Emal = td.Data[0].Close * 10
		td.Data[0].DEA = 0
		td.Data[0].UpdateMacd()
		start = 1
	}

	for i := start; i < c; i++ {
		td.Data[i].Emas = (td.Data[i-1].Emas*s1 + td.Data[i].Close*20) / s2
		td.Data[i].Emal = (td.Data[i-1].Emal*l1 + td.Data[i].Close*20) / l2
		td.Data[i].DEA = (td.Data[i-1].DEA*m1 + td.Data[i].DIFF*2) / m2
		td.Data[i].UpdateMacd()
	}
}
