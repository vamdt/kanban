package crawl

import "time"

func (p *Stock) Ticks2M1s() {
	p.M1s.Drop_lastday_data()
	start_time := p.M1s.latest_time().AddDate(0, 0, 1).Truncate(time.Hour * 24)
	start, _ := (TickSlice(p.Ticks.Data)).Search(start_time)
	for i, c := start, len(p.Ticks.Data); i < c; {
		end := Minuteend(p.Ticks.Data[i].Time)
		hour, min, _ := end.Clock()
		if hour == 9 && min <= 30 {
			end = end.Truncate(30 * time.Minute).Add(31 * time.Minute)
		}
		tdata, j := MergeTickTil(&p.Ticks, i, end)
		p.M1s.Add(tdata)
		i += j
	}
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
	tdata.Time = tdata.Time.Truncate(time.Minute)
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

func (p *Stock) M1s2M5s() {
	p.M5s.Drop_lastday_data()
	start_time := p.M5s.latest_time().AddDate(0, 0, 1).Truncate(time.Hour * 24)
	start, _ := (TdataSlice(p.M1s.Data)).Search(start_time)
	for i, c := start, len(p.M1s.Data); i < c; {
		t := Minute5end(p.M1s.Data[i].Time)
		tdata, j := p.M1s.MergeTil(i, t)
		tdata.Time = t
		p.M5s.Add(tdata)
		i += j
	}
}

func (p *Stock) M1s2M30s() {
	p.M30s.Drop_lastday_data()
	start_time := p.M30s.latest_time().AddDate(0, 0, 1).Truncate(time.Hour * 24)
	start, _ := (TdataSlice(p.M1s.Data)).Search(start_time)
	for i, c := start, len(p.M1s.Data); i < c; {
		t := Minute30end(p.M1s.Data[i].Time)
		tdata, j := p.M1s.MergeTil(i, t)
		tdata.Time = t
		p.M30s.Add(tdata)
		i += j
	}
}

func (p *Stock) Days2Weeks() {
	p.Weeks.Drop_lastday_data()
	start_time := p.Weeks.latest_time().AddDate(0, 0, 1).Truncate(time.Hour * 24)
	start, _ := (TdataSlice(p.Days.Data)).Search(start_time)
	for i, c := start, len(p.Days.Data); i < c; {
		t := Weekend(p.Days.Data[i].Time)
		tdata, j := p.Days.MergeTil(i, t)
		tdata.Time = tdata.Time.Truncate(time.Hour * 24)
		p.Weeks.Add(tdata)
		i += j
	}
}

func (p *Stock) Days2Months() {
	p.Months.Drop_lastday_data()
	start_time := p.Months.latest_time().AddDate(0, 0, 1).Truncate(time.Hour * 24)
	start, _ := (TdataSlice(p.Days.Data)).Search(start_time)
	for i, c := start, len(p.Days.Data); i < c; {
		t := Monthend(p.Days.Data[i].Time)
		tdata, j := p.Days.MergeTil(i, t)
		tdata.Time = tdata.Time.Truncate(time.Hour * 24)
		p.Months.Add(tdata)
		i += j
	}
}
