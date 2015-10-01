package crawl

import "time"

func (p *Stock) WeeksFromDays() {
	for i, c := 0, len(p.Days.Data); i < c; {
		t := p.Days.Data[i].Time
		wd := t.Weekday()
		if wd < time.Saturday {
			t = t.AddDate(0, 0, int(time.Saturday-wd))
		}
		tdata, j := MergeTil(&p.Days, i, t)
		p.Weeks.Add(tdata)
		i += j
	}
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
		if data.Time.After(end) {
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

func (p *Stock) MonthsFromDays() {
	for i, c := 0, len(p.Days.Data); i < c; {
		t := p.Days.Data[i].Time
		_, _, d := t.Date()
		t = t.AddDate(0, 1, -d)
		tdata, j := MergeTil(&p.Days, i, t)
		p.Months.Add(tdata)
		i += j
	}
}
