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
	Time   time.Time
	Open   int
	Close  int
	High   int
	Low    int
	Volume int
}

type Month struct {
	Id     bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Time   time.Time
	Open   int
	Close  int
	High   int
	Low    int
	Volume int
}

func (p *Day) Save(c *mgo.Collection) {
	_, err := c.Upsert(bson.M{"_id": p.Id}, p)
	if err != nil {
		log.Println("insert Day data error", err, *p)
	}
}

func (p *Day) FromString(timestr, open, high, low, cloze, volume []byte) {
	p.Time, _ = time.Parse("2006-01-02", string(timestr))
	p.Id = Time2ObjectId(p.Time)
	p.Open = ParseCent(string(open))
	p.High = ParseCent(string(high))
	p.Low = ParseCent(string(low))
	p.Close = ParseCent(string(cloze))
	p.Volume, _ = strconv.Atoi(string(volume))
}
