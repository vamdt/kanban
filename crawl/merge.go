package crawl

import (
	"time"

	. "./base"
)

func (p *Stock) Ticks2M1s() int {
	p.M1s.Drop_lastday_data()
	index := len(p.M1s.Data)
	start_time := p.M1s.latest_time().AddDate(0, 0, 1).Truncate(time.Hour * 24)
	start, _ := (TickSlice(p.Ticks.Data)).Search(start_time)
	for i, c := start, len(p.Ticks.Data); i < c; {
		end := Minuteend(p.Ticks.Data[i].Time)
		t := end.Add(-1 * time.Minute)
		h, m, _ := t.Clock()
		if h == 9 && m < 30 { // < 9:30
			end = end.Truncate(time.Hour).Add(30 * time.Minute)
			t = end.Add(-1 * time.Minute)
		} else if h == 11 && m > 29 { // > 11:29 11:35
			end = end.Add(5 * time.Minute)
		} else if h == 15 { // > 14:59 15:05
			end = end.Add(5 * time.Minute)
		}
		tdata, j := MergeTickTil(&p.Ticks, i, end)
		i += j
		tdata.Time = t
		k, _ := p.M1s.Add(tdata)
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
		h, m, _ := t.Clock()
		if h == 11 && m == 30 { // > 11:24 11:35
			end = t.Add(5 * time.Minute)
		} else if h == 15 { // > 14:55 15:05
			end = t.Add(5 * time.Minute)
		}
		tdata, j := from.MergeTil(i, end)
		i += j

		if biglevel {
			tdata.Time = tdata.Time.Truncate(time.Hour * 24)
		} else {
			tdata.Time = t
		}
		k, _ := p.Add(tdata)
		if k < index {
			index = k
		}
	}
	return index
}
