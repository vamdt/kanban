package crawl

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

var market_begin_day time.Time

func init() {
	market_begin_day, _ = time.Parse("2006-01-02", "2000-01-01")
}

type Stock struct {
	Id     string  `json:"id"`
	Months []Month `json:"months"`
	Days   Days    `json:"days"`
	Ticks  []Tick  `json:"ticks"`
}

func (p *Stock) day_collection_name() string {
	return fmt.Sprintf("%s.tdata", p.Id)
}

func (p *Stock) day_collection(db *mgo.Database) *mgo.Collection {
	return db.C(p.day_collection_name())
}

func (p *Stock) day_sina_url(t time.Time) string {
	return fmt.Sprintf("http://biz.finance.sina.com.cn/stock/flash_hq/kline_data.php?&rand=9000&symbol=%s&end_date=&begin_date=%s&type=plain",
		p.Id, t.Format("2006-01-02"))
}

func (p *Stock) days_download(t time.Time) (bool, error) {
	body := p.downloadDaysFromSina(t)
	body = bytes.TrimSpace(body)
	lines := bytes.Split(body, []byte("\n"))
	count := len(lines)
	if count < 1 {
		return false, nil
	}

	day := Day{}
	for i := 0; i < count; i++ {
		line := bytes.TrimSpace(lines[i])
		infos := bytes.Split(line, []byte(","))
		if len(infos) != 6 {
			err := errors.New("could not parse line " + string(line))
			return false, err
		}

		day.FromString(infos[0], infos[1], infos[2], infos[3], infos[4], infos[5])
		p.Days.Add(day)
	}
	return true, nil
}

func (p *Stock) downloadDaysFromSina(t time.Time) []byte {
	body, err := Http_get_raw(p.day_sina_url(t), nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	return body
}

func (p *Stock) Days_sync(db *mgo.Database) int {
	c := p.day_collection(db)
	p.Days.Load(c)
	t := p.get_days_latest_time()
	l := len(p.Days.Data)
	p.days_download(t)
	count := len(p.Days.Data)
	if count > l {
		for i, j := l, count; i < j; i++ {
			p.Days.Data[i].Save(c)
		}
	}
	p.Days.Delta = count - l
	return count - l
}

func (p *Stock) get_days_latest_time() time.Time {
	if len(p.Days.Data) < 1 {
		return market_begin_day
	}
	return p.Days.Data[len(p.Days.Data)-1].Time
}

func (p *Stock) get_latest_time_from_db(c *mgo.Collection) time.Time {
	d := Day{}
	err := c.Find(nil).Sort("-_id").Limit(1).One(&d)
	if err != nil {
		log.Println("find fail", err)
		return market_begin_day
	}
	return ObjectId2Time(d.Id)
}
