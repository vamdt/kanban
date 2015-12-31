package crawl

import (
	"flag"
	"time"

	"github.com/golang/glog"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongo string

func init() {
	flag.StringVar(&mongo, "mongo", "mongodb://127.0.0.1/stock", "mongo uri")
}

func NewMongoStore() (ms *MongoStore, err error) {
	session, err := mgo.Dial(mongo)
	if err != nil {
		return
	}
	ms = &MongoStore{session: session}
	return
}

type MongoStore struct {
	session *mgo.Session
}

func (p *MongoStore) Close() {
	if p.session != nil {
		p.session.Close()
	}
}

func (p *MongoStore) LoadTDatas(table string) (res []Tdata, err error) {
	c := p.session.DB("").C(table)
	d := Tdata{}
	iter := c.Find(nil).Sort("_id").Iter()
	for iter.Next(&d) {
		d.Time = ObjectId2Time(d.Id)
		res = append(res, d)
	}
	if err := iter.Close(); err != nil {
		glog.Warningln(err)
	}
	return
}

func (p *MongoStore) SaveTData(table string, data *Tdata) (err error) {
	c := p.session.DB("").C(table)
	_, err = c.Upsert(bson.M{"_id": data.Id}, data)
	if err != nil {
		glog.Warningln("insert tdata error", err, *data)
	}
	return
}

func (p *MongoStore) LoadTicks(table string) (res []Tick, err error) {
	c := p.session.DB("").C(table)
	d := Tick{}
	iter := c.Find(nil).Sort("_id").Iter()
	for iter.Next(&d) {
		d.Time = ObjectId2Time(d.Id)
		res = append(res, d)
	}
	if err := iter.Close(); err != nil {
		glog.Warningln(err)
	}
	return
}

func (p *MongoStore) SaveTick(table string, tick *Tick) (err error) {
	c := p.session.DB("").C(table)
	_, err = c.Upsert(bson.M{"_id": tick.Id}, tick)
	if err != nil {
		glog.Warningln("insert tick error", err, *tick)
	}
	return
}

func (p *MongoStore) TickHasTimeData(table string, t time.Time) bool {
	c := p.session.DB("").C(table)
	begin_id := Time2ObjectId(t)
	end_id := Time2ObjectId(t.AddDate(0, 0, 1))
	n, err := c.Find(bson.M{"_id": bson.M{"$gt": begin_id, "$lt": end_id}}).Count()
	if err != nil {
		glog.Warningln("count fail", err)
		return false
	}
	return n > 0
}
