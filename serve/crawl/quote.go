package crawl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	buy_tick  = 1
	sell_tick = 2
	eq_tick   = 3
)

type Tick struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Time     time.Time
	Price    int
	Change   int
	Volume   int // 手
	Turnover int // 元
	Type     int
}

type Ticks struct {
	Data    []Tick `json:"data"`
	EndTime time.Time
	Delta   int
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

func TickHasInDB(t time.Time, c *mgo.Collection) bool {
	begin_id := Time2ObjectId(t)
	end_id := Time2ObjectId(t.AddDate(0, 0, 1))
	n, err := c.Find(bson.M{"_id": bson.M{"$gt": begin_id, "$lt": end_id}}).Count()
	if err != nil {
		log.Println("count fail", err)
		return false
	}
	return n > 0
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

func FixTickId(ticks []Tick) {
	if ticks == nil {
		return
	}
	c := len(ticks)
	for i := 0; i < c; i++ {
		ticks[i].Id = Time2ObjectId(ticks[i].Time)
	}
}

func FixTickData(ticks []Tick) {
	if ticks == nil {
		return
	}
	c := len(ticks)
	for i := 0; i < c; i++ {
		ticks[i].Id = Time2ObjectId(ticks[i].Time)
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
	url := fmt.Sprintf("http://vip.stock.finance.sina.com.cn/quotes_service/view/CN_TransListV2.php?num=9000&symbol=%s&rn=%ld",
		id, time.Now().UnixNano()/int64(time.Millisecond))
	body, err := Http_get_gbk(url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	return body
}

func Tick_download_real_from_sina(id string) []byte {
	url := fmt.Sprintf("http://hq.sinajs.cn/rn=%ld&list=%s",
		id, time.Now().UnixNano()/int64(time.Millisecond))
	body, err := Http_get_gbk(url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	return body
}

func Tick_collection_name(id string) string {
	return fmt.Sprintf("%s.tick", id)
}

func Tick_collection(db *mgo.Database, id string) *mgo.Collection {
	return db.C(Tick_collection_name(id))
}

func Tick_sina_url(id string, t time.Time) string {
	return fmt.Sprintf("http://market.finance.sina.com.cn/downxls.php?date=%s&symbol=%s",
		t.Format("2006-01-02"), id)
}

func (p *Ticks) Load(c *mgo.Collection) {
	d := Tick{}
	iter := c.Find(nil).Sort("_id").Iter()
	num := len(p.Data)
	for iter.Next(&d) {
		d.Time = ObjectId2Time(d.Id)
		p.Data = append(p.Data, d)
	}
	if err := iter.Close(); err != nil {
		log.Println(err)
	}
	nnum := len(p.Data)
	p.Delta = nnum - num
	if nnum > 0 {
		p.EndTime = p.Data[nnum-1].Time
	}
}

func (p *Ticks) latest_time() time.Time {
	if len(p.Data) < 1 {
		return market_begin_day
	}
	return p.Data[len(p.Data)-1].Time
}

func (p *Tick) Save(c *mgo.Collection) {
	_, err := c.Upsert(bson.M{"_id": p.Id}, p)
	if err != nil {
		log.Println("insert tick error", err, *p)
	}
}

func (p *Ticks) Add(data Tick) {
	if len(p.Data) < 1 {
		p.Data = append(p.Data, data)
		p.Delta++
	} else if data.Time.After(p.Data[len(p.Data)-1].Time) {
		p.Data = append(p.Data, data)
		p.Delta++
	}
	p.EndTime = p.Data[len(p.Data)-1].Time
}

func (p *Ticks) Insert(data Tick) {
	if len(p.Data) < 1 {
		p.Data = []Tick{data}
		p.Delta = 1
	} else if data.Time.After(p.Data[len(p.Data)-1].Time) {
		p.Data = append(p.Data, data)
		p.Delta++
	} else {
		j := len(p.Data) - 1
		should_insert := true
		for i := j - 1; i > -1; i-- {
			if p.Data[i].Time.After(data.Time) {
				j = i
				continue
			} else if p.Data[i].Time.Equal(data.Time) {
				p.Data[i] = data
				should_insert = false
			}
			break
		}

		if should_insert {
			if j < 1 {
				p.Data = append([]Tick{data}, p.Data...)
			} else {
				p.Data = p.Data[0 : len(p.Data)+1]
				copy(p.Data[j+1:], p.Data[j:])
				p.Data[j] = data
			}
		}
	}
	p.EndTime = p.Data[len(p.Data)-1].Time
}
