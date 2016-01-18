package crawl

import "time"

func (p *Stock) Ticks2M1s() int {
	p.M1s.Drop_lastday_data()
	index := len(p.M1s.Data)
	start_time := p.M1s.latest_time().AddDate(0, 0, 1).Truncate(time.Hour * 24)
	start, _ := (TickSlice(p.Ticks.Data)).Search(start_time)
	for i, c := start, len(p.Ticks.Data); i < c; {
		t := Minuteend(p.Ticks.Data[i].Time)
		end := t
		t = t.Truncate(time.Minute)
		h, m, _ := t.UTC().Clock()
		if h == 1 && m < 26 { // < 9:26
			t = t.Truncate(time.Hour).Add(25 * time.Minute)
			end = end.Truncate(time.Hour).Add(26 * time.Minute)
		} else if h == 1 && m < 30 {
			t = t.Truncate(time.Hour).Add(30 * time.Minute)
			end = end.Truncate(time.Hour).Add(31 * time.Minute)
		} else if h == 3 && m > 28 { // >= 11:29 11:35
			end = end.Truncate(time.Hour).Add(35 * time.Minute)
		} else if h == 6 && m > 58 { // >= 14:59 15:05
			end = end.Truncate(time.Hour).Add(65 * time.Minute)
		}
		tdata, j := MergeTickTil(&p.Ticks, i, end)
		i += j
		tdata.Time = t
		k := p.M1s.Add(tdata)
		if k < index {
			index = k
		}
	}
	return index
}

func MergeTickTil(td *Ticks, begin int, end time.Time) (Tdata, int) {
	if begin < 0 {
		begin = 0
	}
	tdata := Tdata{}
	tdata.Open = td.Data[begin].Price
	tdata.High = tdata.Open
	tdata.Low = tdata.Open
	tdata.Volume = 0
	i := begin
	c := len(td.Data)
	for ; i < c; i++ {
		t := td.Data[i]
		if !t.Time.Before(end) {
			break
		}
		tdata.Time = t.Time
		tdata.Close = t.Price
		if t.Price > tdata.High {
			tdata.High = t.Price
		}
		if t.Price < tdata.Low {
			tdata.Low = t.Price
		}
		tdata.Volume += t.Volume
	}
	return tdata, i - begin
}

func (p *Tdatas) MergeTil(begin int, end time.Time) (Tdata, int) {
	if begin < 0 {
		begin = 0
	}
	tdata := p.Data[begin]
	tdata.Volume = 0
	i := begin
	c := len(p.Data)
	for ; i < c; i++ {
		data := p.Data[i]
		if !data.Time.Before(end) {
			break
		}
		tdata.Time = data.Time
		tdata.Close = data.Close
		if data.High > tdata.High {
			tdata.High = data.High
		}
		if data.Low < tdata.Low {
			tdata.Low = data.Low
		}
		tdata.Volume += data.Volume
	}
	return tdata, i - begin
}

func (p *Tdatas) MergeFrom(from *Tdatas, biglevel bool, endtime func(t time.Time) time.Time) int {
	p.Drop_lastday_data()
	index := len(p.Data)
	start_time := p.latest_time().AddDate(0, 0, 1).Truncate(time.Hour * 24)
	start, _ := (TdataSlice(from.Data)).Search(start_time)
	for i, c := start, len(from.Data); i < c; {
		t := endtime(from.Data[i].Time)
		end := t
		h, m, _ := t.UTC().Clock()
		if h == 3 && m == 30 { // > 11:24 11:35
			end = t.Add(5 * time.Minute)
		} else if h == 7 && m == 00 { // > 14:55 15:05
			end = t.Add(5 * time.Minute)
		}
		tdata, j := from.MergeTil(i, end)
		i += j

		if biglevel {
			tdata.Time = tdata.Time.Truncate(time.Hour * 24)
		} else {
			tdata.Time = t
		}
		k := p.Add(tdata)
		if k < index {
			index = k
		}
	}
	return index
}
