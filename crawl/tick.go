package crawl

import (
	"bytes"
	"fmt"
	"time"

	"github.com/golang/glog"

	. "./base"
	"./robot"
)

type Ticks struct {
	Data []Tick `json:"data"`
	play []Tick
}

func (p *Ticks) clean() {
	if len(p.play) > 0 {
		return
	}

	ldata := len(p.Data)
	if ldata < 1 {
		return
	}

	t := p.Data[ldata-1].Time.Truncate(time.Hour * 24)
	i, _ := (TickSlice(p.Data)).Search(t)
	if i > 0 {
		p.Data = p.Data[i:]
	}
}

func FixTickTime(ticks []Tick) {
	if ticks == nil {
		return
	}
	c := len(ticks)
	if c < 2 {
		return
	}
	for i := 1; i < c; i++ {
		t := ticks[i-1].Time
		for j := 1; i < c && t == ticks[i].Time; i++ {
			ticks[i].Time = t.Add(time.Duration(j) * time.Millisecond)
			j++
		}
	}
}

func FixTickData(ticks []Tick) {
	if ticks == nil {
		return
	}
	c := len(ticks)
	for i := 0; i < c; i++ {
		if i == 0 {
			ticks[i].Change = ticks[i].Price
		} else {
			ticks[i].Change = ticks[i].Price - ticks[i-1].Price
		}
		ticks[i].Turnover = ticks[i].Volume * ticks[i].Price / 100
		ticks[i].Volume = ticks[i].Volume / 100
	}
}

func Tick_get_today_date(id string) (time.Time, string, error) {
	body := robot.Tick_download_real_from_sina(id)
	if body == nil {
		return market_begin_day, "", fmt.Errorf("get realtime info fail")
	}

	lines := bytes.Split(body, []byte("\";"))
	if len(lines) < 1 {
		return market_begin_day, "", fmt.Errorf("get realtime info empty")
	}
	info := bytes.Split(lines[0], []byte("=\""))
	if len(info) != 2 {
		return market_begin_day, "", fmt.Errorf("get realtime info format error, donot found =\"")
	}

	infos := bytes.Split(info[1], []byte(","))
	if len(infos) < 33 {
		glog.Warningln("sina hq api, res format changed")
		return market_begin_day, "", fmt.Errorf("sina hq api, res format changed")
	}

	t, err := time.Parse("2006-01-02 15:04:05", string(infos[30])+" "+string(infos[31]))
	return t, string(infos[0]), err
}

func Tick_collection_name(id string) string {
	return fmt.Sprintf("%s.tick", id)
}

func Tick_sina_url(id string, t time.Time) string {
	return fmt.Sprintf("http://market.finance.sina.com.cn/downxls.php?date=%s&symbol=%s",
		t.Format("2006-01-02"), id)
}

func (p *Ticks) latest_time() time.Time {
	if len(p.Data) < 1 {
		return market_begin_day
	}
	return p.Data[len(p.Data)-1].Time
}

func (p *Ticks) hasTimeData(t time.Time) bool {
	i, ok := (TickSlice(p.Data)).Search(t)
	if ok {
		return true
	}
	if i >= len(p.Data) {
		return false
	}
	return t.Truncate(24 * time.Hour).Equal(p.Data[i].Time.Truncate(24 * time.Hour))
}

func (p *Ticks) Add(data Tick) {
	if data.Volume == 0 && data.Price == 0 {
		return
	}
	if data.Volume == 0 && data.Change == 0 && data.Turnover == 0 {
		return
	}
	if len(p.Data) < 1 {
		p.Data = []Tick{data}
	} else if data.Time.After(p.Data[len(p.Data)-1].Time) {
		p.Data = append(p.Data, data)
	} else if data.Time.Equal(p.Data[len(p.Data)-1].Time) {
		p.Data[len(p.Data)-1] = data
	} else {
		i, ok := (TickSlice(p.Data)).Search(data.Time)
		if ok {
			p.Data[i] = data
			return
		}

		if i < 1 {
			p.Data = append([]Tick{data}, p.Data...)
		} else {
			p.Data = append(p.Data, data)
			copy(p.Data[i+1:], p.Data[i:])
			p.Data[i] = data
		}
	}
}
