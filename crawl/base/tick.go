package base

import (
	"sort"
	"strconv"
	"time"
)

const (
	_ int = iota
	Buy_tick
	Sell_tick
	Eq_tick
)

type Tick struct {
	Time     time.Time
	Price    int
	Change   int
	Volume   int // 手
	Turnover int // 元
	Type     int
}

type RealtimeTick struct {
	Tick
	HL
	Buyone  int
	Sellone int
	Status  int
	Name    string
}

type TickSlice []Tick

func (p TickSlice) Len() int           { return len(p) }
func (p TickSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TickSlice) Less(i, j int) bool { return p[i].Time.Before(p[j].Time) }

func SearchTickSlice(a TickSlice, t time.Time) int {
	return sort.Search(len(a), func(i int) bool {
		// a[i].Time >= t
		return a[i].Time.After(t) || a[i].Time.Equal(t)
	})
}

func (p TickSlice) Search(t time.Time) (int, bool) {
	i := SearchTickSlice(p, t)
	if i < p.Len() {
		return i, t.Equal(p[i].Time)
	}
	return i, false
}

func (p *Tick) FromString(date time.Time, timestr, price, change, volume, turnover, typestr []byte) {
	p.Time, _ = time.Parse("15:04:05", string(timestr))
	p.Time = date.Add(time.Second * time.Duration(TheSeconds(p.Time)))

	p.Price = ParseCent(string(price))
	p.Change = ParseCent(string(change))

	p.Volume, _ = strconv.Atoi(string(volume))
	p.Turnover, _ = strconv.Atoi(string(turnover))

	switch string(typestr) {
	case "UP":
		fallthrough
	case "买盘":
		p.Type = Buy_tick
	case "DOWN":
		fallthrough
	case "卖盘":
		p.Type = Sell_tick
	case "EQUAL":
		fallthrough
	case "中性盘":
		p.Type = Eq_tick
	}
}

func TheSeconds(t time.Time) int {
	return t.Hour()*60*60 + t.Minute()*60 + t.Second()
}

func (p *RealtimeTick) SetStatus(s []byte) {
	//"00":"","01":"临停1H","02":"停牌","03":"停牌","04":"临停","05":"停1/2","07":"暂停","-1":"无记录","-2":"未上市","-3":"退市"
	p.Status, _ = strconv.Atoi(string(s))
	if p.Status == 3 {
		p.Status = 2
	}
}
