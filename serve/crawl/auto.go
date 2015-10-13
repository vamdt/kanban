package crawl

import (
	"bytes"
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	tickPeriod = 5 * time.Second
)

var market_begin_day time.Time

func init() {
	market_begin_day, _ = time.Parse("2006-01-02", "2000-01-01")
}

type Stock struct {
	Id        string `json:"id"`
	M1s       M1s    `json:"m1s"`
	M5s       M5s    `json:"m5s"`
	M30s      M30s   `json:"m30s"`
	Days      Days   `json:"days"`
	Weeks     Weeks  `json:"weeks"`
	Months    Months `json:"months"`
	Ticks     Ticks  `json:"-"`
	last_tick RealtimeTick
	hash      int
	count     int
	loaded    int
	lst_trade time.Time
}

type Stocks struct {
	stocks  PStockSlice
	rwmutex sync.RWMutex
	db      *mgo.Database
	ch      chan *Stock
}

func (p *Stocks) Run() {
	go p.update()
	for {
		p.Ticks_update_real()
		time.Sleep(tickPeriod)
	}
}

func (p *Stocks) DB(db *mgo.Database) {
	p.db = db
}

func (p *Stocks) Chan(ch chan *Stock) {
	p.ch = ch
}

func (p *Stocks) res(stock *Stock) {
	if p.ch != nil {
		p.ch <- stock
	}
}

func (p *Stocks) update() {
	for {
		p.rwmutex.RLock()
		for i, c := 0, len(p.stocks); i < c; i++ {
			if p.stocks[i].Update(p.db) {
				p.res(p.stocks[i])
			}
		}
		p.rwmutex.RUnlock()
		time.Sleep(time.Second)
	}
}

func (p *Stocks) Insert(id string) (int, *Stock, bool) {
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()
	s := &Stock{Id: id, hash: StockHash(id), count: 1}
	i, ok := p.stocks.Search(id)
	if ok {
		p.stocks[i].count++
		if p.stocks[i].count < 1 {
			p.stocks[i].count = 1
		}
		return i, p.stocks[i], false
	}

	if i < 1 {
		p.stocks = append(PStockSlice{s}, p.stocks...)
		return 0, s, true
	} else if i >= p.stocks.Len() {
		p.stocks = append(p.stocks, s)
		return p.stocks.Len() - 1, s, true
	}
	p.stocks = append(p.stocks, s)
	copy(p.stocks[i+1:], p.stocks[i:])
	p.stocks[i] = s
	return i, s, true
}

func (p *Stocks) Remove(id string) {
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()
	if i, ok := p.stocks.Search(id); ok {
		p.stocks[i].count--
	}
}

func (p *Stocks) Watch(id string) (*Stock, bool) {
	i, s, isnew := p.Insert(id)
	if isnew {
		log.Println("watch new stock", id, i)
	} else {
		log.Println("watch stock", id, i)
	}
	return s, isnew
}

func (p *Stocks) UnWatch(id string) {
	p.Remove(id)
}

func (p *Stocks) Ticks_update_real() {
	p.rwmutex.RLock()
	defer p.rwmutex.RUnlock()
	var wg sync.WaitGroup

	for i, l := 0, len(p.stocks); i < l; {
		var b bytes.Buffer
		var pstocks PStockSlice
		if i+10 < l {
			pstocks = p.stocks[i : i+10]
		} else {
			pstocks = p.stocks[i:l]
		}
		for j := 0; j < 50 && i < l; i, j = i+1, j+1 {
			if p.stocks[i].loaded < 2 {
				continue
			}
			if p.stocks[i].count < 1 {
				continue
			}
			if b.Len() > 0 {
				b.WriteString(",")
			}
			b.WriteString(p.stocks[i].Id)
		}
		if b.Len() < 1 {
			continue
		}

		wg.Add(1)
		go func(ids string, pstocks PStockSlice) {
			defer wg.Done()
			body := Tick_download_real_from_sina(ids)
			if body == nil {
				return
			}
			for _, line := range bytes.Split(body, []byte("\";")) {
				line = bytes.TrimSpace(line)
				info := bytes.Split(line, []byte("=\""))
				if len(info) != 2 {
					continue
				}
				prefix := "var hq_str_"
				if !bytes.HasPrefix(info[0], []byte(prefix)) {
					continue
				}
				id := info[0][len(prefix):]
				if idx, ok := pstocks.Search(string(id)); ok {
					if pstocks[idx].tick_get_real(info[1]) {
						pstocks[idx].Merge()
						p.res(pstocks[idx])
					}
				}
			}
		}(b.String(), pstocks)

	}
	wg.Wait()
}

func StockHash(id string) int {
	for i, c := range []byte(id) {
		if c >= '0' && c <= '9' {
			i, _ = strconv.Atoi(id[i:])
			return i
		}
	}
	return 0
}

func (p *Stock) Merge() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.Ticks2M1s()
		p.M1s2M5s()
		p.M1s2M30s()
		p.M1s.Macd()
		p.M5s.Macd()
		p.M30s.Macd()
		p.M1s.ParseTyping()
		p.M5s.ParseTyping()
		p.M30s.ParseTyping()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.Days2Weeks()
		p.Days2Months()
		p.Days.Macd()
		p.Weeks.Macd()
		p.Months.Macd()
	}()

	wg.Wait()
}

func (p *Stock) Update(db *mgo.Database) bool {
	if p.loaded > 0 {
		return false
	}
	p.loaded = 1
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		p.Days_update(db)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		p.Ticks_update(db)
		p.Ticks_today_update()
		wg.Done()
	}()

	wg.Wait()
	p.Merge()
	p.loaded = 2
	return true
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

func (p *Stock) Days_update(db *mgo.Database) int {
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

func (p *Stock) Ticks_update(db *mgo.Database) int {
	c := Tick_collection(db, p.Id)
	p.Ticks.Load(c)
	begin_time := p.Ticks.latest_time()
	l := len(p.Ticks.Data)

	now := time.Now().UTC()
	end_time := now.Truncate(time.Hour * 24)
	if now.Hour() > 10 {
		end_time = end_time.AddDate(0, 0, 1)
	}

	if begin_time.Equal(market_begin_day) {
		begin_time = end_time.AddDate(0, -2, -1)
	}
	begin_time = begin_time.AddDate(0, 0, 1).Truncate(time.Hour * 24)

	for t := begin_time; t.Before(end_time); t = t.AddDate(0, 0, 1) {
		if !IsTradeDay(t) {
			log.Println(t, "skip non trading day")
			continue
		}

		if TickHasInDB(t, c) {
			log.Println(t, "already in db, skip")
			continue
		}

		log.Println("prepare download ticks", t)
		if ok, err := p.ticks_download(t); ok {
			log.Println("download ticks succ", t)
		} else if err != nil {
			log.Println("download ticks err", err)
		}
	}

	count := len(p.Ticks.Data)
	if count > l {
		for i, j := l, count; i < j; i++ {
			p.Ticks.Data[i].Save(c)
		}
	}
	p.Ticks.Delta = count - l
	log.Println("download ticks", p.Ticks.Delta)
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

func (p *Tdata) parse_mins_from_sina(line []byte) error {
	items := [6]string{"day:", "open:", "high:", "close:", "low:", "volume:"}
	v := [6]string{}
	line = bytes.TrimSpace(line)
	line = bytes.Trim(line, "[{}]")
	infos := bytes.Split(line, []byte(","))
	if len(infos) != 6 {
		return errors.New("could not parse line " + string(line))
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

	p.FromString(v[0], v[1], v[2], v[3], v[4], v[5])
	return nil
}

var UnknowSinaRes error = errors.New("could not find '成交时间' in head line")

func (p *Stock) ticks_download(t time.Time) (bool, error) {
	body := Tick_download_from_sina(p.Id, t)
	if body == nil {
		return false, nil
	}
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

	ticks := make([]Tick, count)
	for i := count; i > 0; i-- {
		line := bytes.TrimSpace(lines[i])
		infos := bytes.Split(line, []byte("\t"))
		if len(infos) != 6 {
			err := errors.New("could not parse line " + string(line))
			return false, err
		}
		ticks[count-i].FromString(t, infos[0], infos[1], infos[2],
			infos[3], infos[4], infos[5])
	}
	FixTickTime(ticks)
	FixTickId(ticks)

	for _, tick := range ticks {
		p.Ticks.Add(tick)
	}
	return true, nil
}

func (p *Stock) Ticks_today_update() int {
	l := len(p.Ticks.Data)

	now := time.Now().UTC()
	if !IsTradeDay(now) {
		return 0
	}

	nhour := now.Hour()
	if nhour < 1 || nhour > 10 {
		return 0
	}

	p.ticks_get_today()

	count := len(p.Ticks.Data)
	p.Ticks.Delta = count - l
	return p.Ticks.Delta
}

func (p *Stock) ticks_get_today() bool {
	last_t, err := Tick_get_today_date(p.Id)
	if err != nil {
		log.Println("get today date fail", err)
		return false
	}
	t := time.Now().UTC().Truncate(time.Hour * 24)
	if t.After(last_t) {
		return false
	}

	body := Tick_download_today_from_sina(p.Id)
	if body == nil {
		return false
	}
	body = bytes.TrimSpace(body)
	lines := bytes.Split(body, []byte("\n"))
	count := len(lines) - 2
	if count < 1 {
		return false
	}

	ticks := make([]Tick, count)
	nul := []byte("")
	for i, j := len(lines)-1, 0; i > 0 && j < count; i-- {
		line := bytes.TrimSpace(lines[i])
		line = bytes.Trim(line, ");")
		infos := bytes.Split(line, []byte("] = new Array("))
		if len(infos) != 2 {
			continue
		}
		line = bytes.Replace(infos[1], []byte(" "), nul, -1)
		line = bytes.Replace(line, []byte("'"), nul, -1)
		infos = bytes.Split(line, []byte(","))
		if len(infos) != 4 {
			continue
		}

		ticks[j].FromString(t, infos[0], infos[2], nul, infos[1], nul, infos[3])
		j++
	}
	FixTickTime(ticks)
	FixTickData(ticks)

	for _, tick := range ticks {
		p.Ticks.Insert(tick)
	}
	return true
}

func (p *Stock) tick_get_real(line []byte) bool {
	infos := bytes.Split(line, []byte(","))
	if len(infos) < 33 {
		log.Println("sina hq api, res format changed")
		return false
	}

	nul := []byte("")
	tick := RealtimeTick{}
	t, _ := time.Parse("2006-01-02", string(infos[30]))
	tick.FromString(t, infos[31], infos[3], nul, infos[8], infos[9], nul)
	tick.buyone = ParseCent(string(infos[11]))
	tick.sellone = ParseCent(string(infos[21]))
	tick.set_status(infos[32])

	if p.last_tick.Volume == 0 {
		p.last_tick = tick
		if tick.Time.Before(p.lst_trade) {
			p.last_tick.Volume = 0
		}
		return false
	}
	if tick.Volume != p.last_tick.Volume {
		if tick.Price >= p.last_tick.sellone {
			tick.Type = buy_tick
		} else if tick.Price <= p.last_tick.buyone {
			tick.Type = sell_tick
		} else {
			tick.Type = eq_tick
		}
		tick.Change = tick.Price - p.last_tick.Price

		volume := (tick.Volume - p.last_tick.Volume) / 100
		p.last_tick = tick
		tick.Volume = volume
		p.Ticks.Insert(tick.Tick)
		p.lst_trade = tick.Time
		return true
	}
	return false
}
