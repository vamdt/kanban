package crawl

import "time"

func (p *Stock) Ticks2M1s() {
	for i, c := 0, len(p.Ticks.Data); i < c; {
		end := Minuteend(p.Ticks.Data[i].Time)
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
	tdata.Time = tdata.Time.Truncate(time.Minute).Add(time.Minute)
	return tdata, i - begin
}

func MergeTil(td *Days, begin int, end time.Time) (Tdata, int) {
	if begin < 0 {
		begin = 0
	}
	tdata := td.Data[begin]
	tdata.Volume = 0
	i := begin
	c := len(td.Data)
	for ; i < c; i++ {
		data := td.Data[i]
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
	tdata.Time = tdata.Time.Truncate(time.Hour * 24)
	return tdata, i - begin
}

func (p *Stock) Days2Weeks() {
	for i, c := 0, len(p.Days.Data); i < c; {
		t := Weekend(p.Days.Data[i].Time)
		tdata, j := MergeTil(&p.Days, i, t)
		p.Weeks.Add(tdata)
		i += j
	}
}

func (p *Stock) Days2Months() {
	for i, c := 0, len(p.Days.Data); i < c; {
		t := Monthend(p.Days.Data[i].Time)
		tdata, j := MergeTil(&p.Days, i, t)
		p.Months.Add(tdata)
		i += j
	}
}
