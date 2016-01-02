package crawl

import (
	"flag"

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
		res = append(res, d)
	}
	if err := iter.Close(); err != nil {
		glog.Warningln(err)
	}
	return
}

func data2BsonM(data interface{}) (m bson.M, err error) {
	m = make(bson.M)
	buf, err := bson.Marshal(data)
	if err != nil {
		return
	}
	err = bson.Unmarshal(buf, m)
	return
}

func (p *MongoStore) SaveTData(table string, data *Tdata) (err error) {
	c := p.session.DB("").C(table)
	id := Time2ObjectId(data.Time)
	m, err := data2BsonM(*data)
	if err != nil {
		glog.Warningln("convert tdata error", err, *data)
		return
	}
	m["_id"] = id
	_, err = c.Upsert(bson.M{"_id": id}, m)
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
		res = append(res, d)
	}
	if err := iter.Close(); err != nil {
		glog.Warningln(err)
	}
	return
}

func (p *MongoStore) SaveTick(table string, tick *Tick) (err error) {
	c := p.session.DB("").C(table)
	id := Time2ObjectId(tick.Time)
	m, err := data2BsonM(*tick)
	if err != nil {
		glog.Warningln("convert tick error", err, *tick)
		return
	}
	m["_id"] = id
	_, err = c.Upsert(bson.M{"_id": id}, m)
	if err != nil {
		glog.Warningln("insert tick error", err, *tick)
	}
	return
}
