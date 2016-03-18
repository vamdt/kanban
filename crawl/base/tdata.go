package base

import (
	"sort"
	"strconv"
	"time"
)

const (
	lmt    = "2006-01-02 15:04:05"
	smt    = "2006-01-02"
	QQmt   = "060102"
	JQKAmt = "20060102"
	l_lmt  = len(lmt)
	l_smt  = len(smt)
	l_qqmt = len(QQmt)
	l_jqka = len(JQKAmt)
)

type Tdata struct {
	Time   time.Time
	Open   int       `json:"open"`
	Close  int       `json:"close"`
	Volume int       `json:"volume"`
	HL     `bson:",inline"`
	Emas   int `json:"-"`
	Emal   int `json:"-"`
	DIFF   int
	DEA    int
	MACD   int
}

type TdataSlice []Tdata

func (p TdataSlice) Len() int           { return len(p) }
func (p TdataSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TdataSlice) Less(i, j int) bool { return p[i].Time.Before(p[j].Time) }

func SearchTdataSlice(a TdataSlice, t time.Time) int {
	return sort.Search(len(a), func(i int) bool {
		// a[i].Time >= t
		return a[i].Time.After(t) || a[i].Time.Equal(t)
	})
}

func (p TdataSlice) Search(t time.Time) (int, bool) {
	i := SearchTdataSlice(p, t)
	if i < p.Len() {
		return i, t.Equal(p[i].Time)
	}
	return i, false
}

func (p *Tdata) FromBytes(timestr, open, high, cloze, low, volume []byte) {
	p.FromString(string(timestr), string(open), string(high), string(cloze),
		string(low), string(volume))
}

func (p *Tdata) FromString(timestr, open, high, cloze, low, volume string) {
	ltime := len(timestr)
	switch ltime {
	case l_lmt:
		p.Time, _ = time.Parse(lmt, timestr)
	case l_smt:
		p.Time, _ = time.Parse(smt, timestr)
	case l_qqmt:
		p.Time, _ = time.Parse(QQmt, timestr)
	case l_jqka:
		p.Time, _ = time.Parse(JQKAmt, timestr)
	}
	p.Open = ParseCent(open)
	p.High = ParseCent(high)
	p.Low = ParseCent(low)
	p.Close = ParseCent(cloze)
	p.Volume, _ = strconv.Atoi(volume)
}
