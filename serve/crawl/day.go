package crawl

import (
	"log"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Day struct {
	Id     bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Time   time.Time     `json:"time"`
	Open   int           `json:"open"`
	Close  int           `json:"close"`
	High   int           `json:"high"`
	Low    int           `json:"low"`
	Volume int           `json:"volume"`
}

type Days struct {
	Data    []Day `json:"data"`
	EndTime time.Time
	Delta   int
}

type Month struct {
	Id     bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Time   time.Time     `json:"time"`
	Open   int           `json:"open"`
	Close  int           `json:"close"`
	High   int           `json:"high"`
	Low    int           `json:"low"`
	Volume int           `json:"volume"`
}

func (p *Days) Load(c *mgo.Collection) {
	d := Day{}
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

func (p *Days) Add(day Day) {
	if len(p.Data) < 1 {
		p.Data = append(p.Data, day)
	} else if day.Time.After(p.Data[len(p.Data)-1].Time) {
		p.Data = append(p.Data, day)
	}
	p.Delta++
	p.EndTime = p.Data[len(p.Data)-1].Time
}

func (p *Day) Save(c *mgo.Collection) {
	_, err := c.Upsert(bson.M{"_id": p.Id}, p)
	if err != nil {
		log.Println("insert Day data error", err, *p)
	}
}

func (p *Day) FromString(timestr, open, high, cloze, low, volume []byte) {
	p.Time, _ = time.Parse("2006-01-02", string(timestr))
	p.Id = Time2ObjectId(p.Time)
	p.Open = ParseCent(string(open))
	p.High = ParseCent(string(high))
	p.Low = ParseCent(string(low))
	p.Close = ParseCent(string(cloze))
	p.Volume, _ = strconv.Atoi(string(volume))
}
