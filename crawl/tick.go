package crawl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"time"
)

const (
	_ int = iota
	buy_tick
	sell_tick
	eq_tick
)

type Tick struct {
	Time     time.Time
	Price    int
	Change   int
	Volume   int // 手
	Turnover int // 元
	Type     int
}

type Ticks struct {
	Data []Tick `json:"data"`
	play []Tick
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

type RealtimeTick struct {
	Tick
	buyone  int
	sellone int
	status  int
}

func (p *RealtimeTick) set_status(s []byte) {
	//"00":"","01":"临停1H","02":"停牌","03":"停牌","04":"临停","05":"停1/2","07":"暂停","-1":"无记录","-2":"未上市","-3":"退市"
	p.status, _ = strconv.Atoi(string(s))
	if p.status == 3 {
		p.status = 2
	}
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
		p.Type = buy_tick
	case "DOWN":
		fallthrough
	case "卖盘":
		p.Type = sell_tick
	case "EQUAL":
		fallthrough
	case "中性盘":
		p.Type = eq_tick
	}
}

func raw_cache_filename(id string, t time.Time) string {
	return path.Join(os.Getenv("HOME"), "cache", t.Format("2006/0102"), id)
}

func tick_read_raw_cache(id string, t time.Time) ([]byte, error) {
	return ioutil.ReadFile(raw_cache_filename(id, t))
}

func tick_write_raw_cache(c []byte, id string, t time.Time) {
	if len(c) < 1 {
		return
	}
	f := raw_cache_filename(id, t)
	os.MkdirAll(path.Dir(f), 0755)
	ioutil.WriteFile(f, c, 0644)
}

func TheSeconds(t time.Time) int {
	return t.Hour()*60*60 + t.Minute()*60 + t.Second()
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

func Tick_download_from_sina(id string, t time.Time) []byte {
	body, err := tick_read_raw_cache(id, t)
	if err == nil {
		return body
	}

	body, err = Http_get_gbk(Tick_sina_url(id, t), nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	tick_write_raw_cache(body, id, t)
	return body
}

func Tick_download_today_from_sina(id string) []byte {
	url := fmt.Sprintf("http://vip.stock.finance.sina.com.cn/quotes_service/view/CN_TransListV2.php?num=9000&symbol=%s&rn=%d",
		id, time.Now().UnixNano()/int64(time.Millisecond))
	body, err := Http_get_gbk(url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	return body
}

func Tick_download_real_from_sina(id string) []byte {
	if len(id) < 1 {
		return nil
	}
	url := fmt.Sprintf("http://hq.sinajs.cn/rn=%d&list=%s",
		time.Now().UnixNano()/int64(time.Millisecond), id)
	body, err := Http_get_gbk(url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	return body
}

func Tick_get_today_date(id string) (time.Time, error) {
	body := Tick_download_real_from_sina(id)
	if body == nil {
		return market_begin_day, fmt.Errorf("get realtime info fail")
	}

	lines := bytes.Split(body, []byte("\";"))
	if len(lines) < 1 {
		return market_begin_day, fmt.Errorf("get realtime info empty")
	}
	info := bytes.Split(lines[0], []byte("=\""))
	if len(info) != 2 {
		return market_begin_day, fmt.Errorf("get realtime info format error, donot found =\"")
	}

	infos := bytes.Split(info[1], []byte(","))
	if len(infos) < 33 {
		log.Println("sina hq api, res format changed")
		return market_begin_day, fmt.Errorf("sina hq api, res format changed")
	}

	return time.Parse("2006-01-02 15:04:05", string(infos[30])+" "+string(infos[31]))
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
	end := t.AddDate(0, 0, 1)
	for i := len(p.Data) - 1; i > -1; i-- {
		if p.Data[i].Time.Equal(t) {
			return true
		} else if p.Data[i].Time.After(t) && p.Data[i].Time.Before(end) {
			return true
		}
	}
	return false
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
