package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Tick struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	time     time.Time
	Price    int
	Change   int
	Volume   int // 手
	Turnover int // 元
	Type     int
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

func ParseCent(s string) int {
	ms := strings.SplitN(s, ".", 3)

	m, _ := strconv.Atoi(ms[0])
	if m < 0 {
		m = -m
	}

	var cent string
	if len(ms) > 1 {
		cent = ms[1]
	}
	cent = cent + "00"
	cent = cent[:2]
	c, _ := strconv.Atoi(cent)
	if s[:1] == "-" {
		return 100*m - c
	}
	return 100*m + c
}

func (p *Tick) FromString(date time.Time, timestr, price, change, volume, turnover, typestr []byte) {
	p.time, _ = time.Parse("15:04:05", string(timestr))
	p.time = date.Add(time.Second * time.Duration(TheSeconds(p.time)))

	p.Price = ParseCent(string(price))
	p.Change = ParseCent(string(change))

	p.Volume, _ = strconv.Atoi(string(volume))
	p.Turnover, _ = strconv.Atoi(string(turnover))

	switch string(typestr) {
	case "买盘":
		p.Type = 1
	case "卖盘":
		p.Type = 2
	case "中性盘":
		p.Type = 3
	}
}

type Quote struct {
	Stock string
	Ticks []Tick
}

func NewQuote(stock string) *Quote {
	return &Quote{Stock: stock}
}

func (p *Quote) collectionName() string {
	return fmt.Sprintf("%s.tick", p.Stock)
}

func (p *Quote) Collection(db *mgo.Database) *mgo.Collection {
	return db.C(p.collectionName())
}

func (p *Quote) collectionTodayName() string {
	return fmt.Sprintf("%s.tick.0day", p.Stock)
}

func (p *Quote) CollectionToday(db *mgo.Database) *mgo.Collection {
	return db.C(p.collectionTodayName())
}

func (p *Quote) sinaQuoteUrl(t time.Time) string {
	return fmt.Sprintf("http://market.finance.sina.com.cn/downxls.php?date=%s&symbol=%s",
		t.Format("2006-01-02"), p.Stock)
}

func (p *Quote) sinaTodayQuoteUrl(t time.Time) string {
	return fmt.Sprintf("http://vip.stock.finance.sina.com.cn/quotes_service/view/CN_TransListV2.php?num=100000&symbol=%s&rn=%ld",
		p.Stock, t.UnixNano()/int64(time.Millisecond))
}

func (p *Quote) rawCacheFilename(t time.Time) string {
	return path.Join(os.Getenv("HOME"), "cache", t.Format("2006/0102"), p.Stock)
}

func (p *Quote) readRawCache(t time.Time) ([]byte, error) {
	return ioutil.ReadFile(p.rawCacheFilename(t))
}

func (p *Quote) writeRawCache(c []byte, t time.Time) {
	if len(c) < 1 {
		return
	}
	f := p.rawCacheFilename(t)
	os.MkdirAll(path.Dir(f), 0755)
	ioutil.WriteFile(p.rawCacheFilename(t), c, 0644)
}

func TheSeconds(t time.Time) int {
	return t.Hour()*60*60 + t.Minute()*60 + t.Second()
}

func (p *Quote) UpdateToday(db *mgo.Database) {
}

func (p *Quote) Update(db *mgo.Database, days int) {
	now := time.Now().UTC()
	t := now.Truncate(time.Hour * 24)
	if now.Hour() > 10 {
		t = t.AddDate(0, 0, 1)
	}

	if days < 1 {
		days = 5
	}

	i := 0
	for i < days {
		t = t.AddDate(0, 0, -1)
		if !IsTradeDay(t) {
			log.Println(t, "skip non trading day")
			continue
		}

		i++
		if p.HasInDB(t, db) {
			log.Println(t, "already in db, skip")
			continue
		}

		log.Println("prepare download ticks", t)
		if i > 1 {
			time.Sleep(time.Second)
		}

		if ok, err := p.downloadTicks(t); ok {
			p.Save(db, t)
		} else if err != nil {
			log.Println(err)
		}
	}
}

func (p *Quote) HasInDB(t time.Time, db *mgo.Database) bool {
	c := p.Collection(db)
	begin_id := Time2ObjectId(t)
	end_id := Time2ObjectId(t.AddDate(0, 0, 1))
	n, err := c.Find(bson.M{"_id": bson.M{"$gt": begin_id, "$lt": end_id}}).Count()
	if err != nil {
		log.Println("count fail", err)
		return false
	}
	return n > 0
}

var UnknowSinaRes error = errors.New("could not find '成交时间' in head line")

func (p *Quote) downloadTicks(t time.Time) (bool, error) {
	body := p.downloadFromSina(t)
	body = bytes.TrimSpace(body)
	lines := bytes.Split(body, []byte("\n"))
	count := len(lines) - 1
	if count < 1 {
		return false, nil
	}
	if bytes.Contains(lines[0], []byte("script")) {
		return false, nil
	}
	if !bytes.Contains(lines[0], []byte("成交时间")) {
		return false, UnknowSinaRes
	}

	p.Ticks = make([]Tick, count)
	for i := count; i > 0; i-- {
		line := bytes.TrimSpace(lines[i])
		infos := bytes.Split(line, []byte("\t"))
		if len(infos) != 6 {
			err := errors.New("could not parse line " + string(line))
			return false, err
		}
		p.Ticks[count-i].FromString(t, infos[0], infos[1], infos[2],
			infos[3], infos[4], infos[5])
	}
	p.fixTickTime()
	p.fixTickId()
	return true, nil
}

func (p *Quote) fixTickTime() {
	c := len(p.Ticks)
	if c < 2 {
		return
	}
	for i := 1; i < c; i++ {
		t := p.Ticks[i-1].time
		for j := 1; i < c && t == p.Ticks[i].time; i++ {
			p.Ticks[i].time = t.Add(time.Duration(j) * time.Millisecond)
			j++
		}
	}
}

func (p *Quote) fixTickId() {
	c := len(p.Ticks)
	for i := 0; i < c; i++ {
		p.Ticks[i].Id = Time2ObjectId(p.Ticks[i].time)
	}
}

func (p *Quote) Dump() {
	c := len(p.Ticks)
	for i := 0; i < c; i++ {
		fmt.Printf("%+v\n", p.Ticks[i])
	}
}

func (p *Quote) downloadFromSina(t time.Time) []byte {
	body, err := p.readRawCache(t)
	if err == nil {
		return body
	}

	log.Println(p.sinaQuoteUrl(t))
	resp, err := http_get(p.sinaQuoteUrl(t), nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(transform.NewReader(resp.Body,
		simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		log.Println(err)
		return nil
	}
	p.writeRawCache(body, t)
	return body
}

func (p *Quote) Save(db *mgo.Database, t time.Time) {
	c := p.Collection(db)
	dlen := len(p.Ticks)
	for i := 0; i < dlen; i++ {
		_, err := c.Upsert(bson.M{"_id": p.Ticks[i].Id}, &p.Ticks[i])
		if err != nil {
			log.Fatal("insert Quote error", err, p.Ticks[i])
		}
	}
}
