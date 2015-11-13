package crawl

import (
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
	emas   int
	emal   int
	DIFF   int
	DEA    int
	MACD   int
}

type Tdatas struct {
	Data    []Tdata `json:"data"`
	Typing  typing_parser
	Segment segment_parser
	Hub     []Typing
	EndTime time.Time

	min_hub_height int
}

func (p *Tdatas) Load(c *mgo.Collection) {
	var data []Tdata
	d := Tdata{}
	iter := c.Find(nil).Sort("_id").Iter()
	for iter.Next(&d) {
		d.Time = ObjectId2Time(d.Id)
		data = append(data, d)
	}
	if err := iter.Close(); err != nil {
		log.Println(err)
	}
	p.Data = data
	nnum := len(p.Data)
	if nnum > 0 {
		p.EndTime = p.Data[nnum-1].Time
	}
}

func (p *Tdatas) Add(data Tdata) {
	if len(p.Data) < 1 {
		p.Data = append(p.Data, data)
	} else if data.Time.After(p.Data[len(p.Data)-1].Time) {
		p.Data = append(p.Data, data)
	} else if data.Time.Equal(p.Data[len(p.Data)-1].Time) {
		p.Data[len(p.Data)-1] = data
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
				p.Data = append([]Tdata{data}, p.Data...)
			} else {
				p.Data = append(p.Data, data)
				copy(p.Data[j+1:], p.Data[j:])
				p.Data[j] = data
			}
		}
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
