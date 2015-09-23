package crawl

import (
	"bytes"
	"errors"
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

var market_begin_day time.Time

func init() {
	market_begin_day, _ = time.Parse("2006-01-02", "2000-01-01")
}

type Stock struct {
	Id     string `json:"id"`
	M1s    M1s    `json:"m1s"`
	M5s    M5s    `json:"m5s"`
	M30s   M30s   `json:"m30s"`
	Days   Days   `json:"days"`
	Weeks  Weeks  `json:"weeks"`
	Months Months `json:"months"`
	Ticks  []Tick `json:"-"`
}

func (p *Stock) days_download(t time.Time) (bool, error) {
	body := DownloadDaysFromSina(p.Id, t)
	body = bytes.TrimSpace(body)
	lines := bytes.Split(body, []byte("\n"))
	count := len(lines)
	if count < 1 {
		return false, nil
	}

	day := Tdata{}
	for i := 0; i < count; i++ {
		line := bytes.TrimSpace(lines[i])
		infos := bytes.Split(line, []byte(","))
		if len(infos) != 6 {
			err := errors.New("could not parse line " + string(line))
			return false, err
		}

		day.FromBytes(infos[0], infos[1], infos[2], infos[3], infos[4], infos[5])
		p.Days.Add(day)
	}
	return true, nil
}

func (p *Stock) Days_sync(db *mgo.Database) int {
	c := Day_collection(db, p.Id)
	p.Days.Load(c)
	t := p.Days.latest_time()
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

func (p *Stock) M5s_sync(db *mgo.Database) int {
	c := M5_collection(db, p.Id)
	p.M5s.Load(c)
	l := len(p.M5s.Data)
	p.m5s_download()
	count := len(p.M5s.Data)
	if count > l {
		for i, j := l, count; i < j; i++ {
			p.M5s.Data[i].Save(c)
		}
	}
	p.M5s.Delta = count - l
	return count - l
}

func (p *Stock) M30s_sync(db *mgo.Database) int {
	c := M30_collection(db, p.Id)
	p.M30s.Load(c)
	l := len(p.M30s.Data)
	p.m30s_download()
	count := len(p.M30s.Data)
	if count > l {
		for i, j := l, count; i < j; i++ {
			p.M30s.Data[i].Save(c)
		}
	}
	p.M30s.Delta = count - l
	return count - l
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

func (p *Stock) m30s_download() (bool, error) {
	body := DownloadM30sFromSina(p.Id)
	body = bytes.TrimSpace(body)
	lines := bytes.Split(body, []byte("},{"))
	count := len(lines)
	if count < 1 {
		return false, nil
	}

	data := Tdata{}
	items := [6]string{"day:", "open:", "high:", "close:", "low:", "volume:"}
	v := [6]string{}

	for i := 0; i < count; i++ {
		line := bytes.TrimSpace(lines[i])
		line = bytes.Trim(line, "[{}]")
		infos := bytes.Split(line, []byte(","))
		if len(infos) != 6 {
			err := errors.New("could not parse line " + string(line))
			return false, err
		}

		for i, item := range items {
			v[i] = ""
			for _, info := range infos {
				if bytes.HasPrefix(info, []byte(item)) {
					info = bytes.TrimPrefix(info, []byte(item))
					info = bytes.Trim(info, "\"")
					v[i] = string(info)
				}
			}
		}

		data.FromString(v[0], v[1], v[2], v[3], v[4], v[5])
		p.M30s.Add(data)
	}

	return true, nil
}

func (p *Stock) m5s_download() (bool, error) {
	body := DownloadM5sFromSina(p.Id)
	body = bytes.TrimSpace(body)
	lines := bytes.Split(body, []byte("},{"))
	count := len(lines)
	if count < 1 {
		return false, nil
	}

	data := Tdata{}
	items := [6]string{"day:", "open:", "high:", "close:", "low:", "volume:"}
	v := [6]string{}

	for i := 0; i < count; i++ {
		line := bytes.TrimSpace(lines[i])
		line = bytes.Trim(line, "[{}]")
		infos := bytes.Split(line, []byte(","))
		if len(infos) != 6 {
			err := errors.New("could not parse line " + string(line))
			return false, err
		}

		for i, item := range items {
			v[i] = ""
			for _, info := range infos {
				if bytes.HasPrefix(info, []byte(item)) {
					info = bytes.TrimPrefix(info, []byte(item))
					info = bytes.Trim(info, "\"")
					v[i] = string(info)
				}
			}
		}

		data.FromString(v[0], v[1], v[2], v[3], v[4], v[5])
		p.M5s.Add(data)
	}
	return true, nil
}
