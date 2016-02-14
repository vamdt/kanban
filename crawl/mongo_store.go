// +build ignore

package crawl

import (
	"encoding/binary"
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongo string

func init() {
	flag.StringVar(&mongo, "mongo", "mongodb://127.0.0.1/stock", "mongo uri")
	RegisterStore("mongo", &MongoStore{})
}

func (p *MongoStore) Open() (err error) {
	if p.session != nil {
		p.Close()
	}

	session, err := mgo.Dial(mongo)
	if err != nil {
		return
	}
	p.session = session
	return
}

type MongoStore struct {
	session *mgo.Session
}

func (p *MongoStore) Close() {
	if p.session != nil {
		p.session.Close()
		p.session = nil
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

func Time2ObjectId(t time.Time) bson.ObjectId {
	var b [12]byte
	binary.BigEndian.PutUint32(b[:4], uint32(t.Unix()))
	binary.BigEndian.PutUint16(b[4:6], uint16(t.Nanosecond()/int(time.Millisecond)))
	return bson.ObjectId(string(b[:]))
}

func ObjectId2Time(oid bson.ObjectId) time.Time {
	id := string(oid)
	if len(oid) != 12 {
		panic(fmt.Sprintf("Invalid ObjectId: %q", id))
	}
	secs := int64(binary.BigEndian.Uint32([]byte(id[0:4])))
	nsec := int64(binary.BigEndian.Uint16([]byte(id[4:6]))) * int64(time.Millisecond)
	return time.Unix(secs, nsec).UTC()
}

func (p *MongoStore) LoadCategories() (res TopCategory, err error) { return }
func (p *MongoStore) SaveCategories(c TopCategory) (err error)     { return }
