package crawl

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	lmt   = "2006-01-02 15:04:05"
	smt   = "2006-01-02"
	l_lmt = len(lmt)
	l_smt = len(smt)
)

type Tdata struct {
	Id     bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Time   time.Time     `json:"time"`
	Open   int           `json:"open"`
	Close  int           `json:"close"`
	High   int           `json:"high"`
	Low    int           `json:"low"`
	Volume int           `json:"volume"`
}

type Day struct {
	Tdata
}

type Week struct {
	Tdata
}

type Month struct {
	Tdata
}

type Tdatas struct {
	Data    []Tdata `json:"data"`
	EndTime time.Time
	Delta   int
}

type Days struct {
	Tdatas
}

type Weeks struct {
	Tdatas
}

type Months struct {
	Tdatas
}

func (p *Tdatas) Load(c *mgo.Collection) {
	d := Tdata{}
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

func (p *Tdatas) Add(data Tdata) {
	if len(p.Data) < 1 {
		p.Data = append(p.Data, data)
		p.Delta++
	} else if data.Time.After(p.Data[len(p.Data)-1].Time) {
		p.Data = append(p.Data, data)
		p.Delta++
	}
	p.EndTime = p.Data[len(p.Data)-1].Time
}

func (p *Tdatas) latest_time() time.Time {
	if len(p.Data) < 1 {
		return market_begin_day
	}
	return p.Data[len(p.Data)-1].Time
}

func (p *Tdata) Save(c *mgo.Collection) {
	_, err := c.Upsert(bson.M{"_id": p.Id}, p)
	if err != nil {
		log.Println("insert tdata error", err, *p)
	}
}

func (p *Tdata) FromBytes(timestr, open, high, cloze, low, volume []byte) {
	p.FromString(string(timestr), string(open), string(high), string(cloze),
		string(low), string(volume))
}

func (p *Tdata) FromString(timestr, open, high, cloze, low, volume string) {
	if len(timestr) == l_lmt {
		p.Time, _ = time.Parse(lmt, timestr)
	} else {
		p.Time, _ = time.Parse(smt, timestr)
	}
	p.Id = Time2ObjectId(p.Time)
	p.Open = ParseCent(open)
	p.High = ParseCent(high)
	p.Low = ParseCent(low)
	p.Close = ParseCent(cloze)
	p.Volume, _ = strconv.Atoi(volume)
}

func Day_collection_name(id string) string {
	return fmt.Sprintf("%s.tdata.kday", id)
}

func Day_collection(db *mgo.Database, id string) *mgo.Collection {
	return db.C(Day_collection_name(id))
}

func Day_sina_url(id string, t time.Time) string {
	return fmt.Sprintf("http://biz.finance.sina.com.cn/stock/flash_hq/kline_data.php?&rand=9000&symbol=%s&end_date=&begin_date=%s&type=plain",
		id, t.Format("2006-01-02"))
}

func DownloadDaysFromSina(id string, t time.Time) []byte {
	url := Day_sina_url(id, t)
	return Download(url)
}
