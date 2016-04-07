package crawl

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	. "./base"
	"./robot"
	"./store"
	"github.com/golang/glog"
)

const (
	tickPeriod = 5 * time.Second
	minPlay    = 1
)

var market_begin_day time.Time

func init() {
	market_begin_day, _ = time.Parse("2006-01-02", "1990-12-19")
}

type Stock struct {
	Id        string `json:"id"`
	M1s       Tdatas `json:"m1s"`
	M5s       Tdatas `json:"m5s"`
	M30s      Tdatas `json:"m30s"`
	Days      Tdatas `json:"days"`
	Weeks     Tdatas `json:"weeks"`
	Months    Tdatas `json:"months"`
	Ticks     Ticks  `json:"-"`
	last_tick RealtimeTick
	hash      int
	count     int32
	loaded    int32
	broadcast bool
	lst_trade time.Time
	rw        sync.RWMutex
	Name      string
}

func (p *Stock) MarshalTail(tail bool) ([]byte, error) {
	p.rw.RLock()
	defer p.rw.RUnlock()
	s := Stock{
		Id:   p.Id,
		Name: p.Name,
	}
	if !tail || !p.broadcast {
		p.broadcast = true
		// full
		p.M1s.tail(&s.M1s, 0)
		p.M5s.tail(&s.M5s, 0)
		p.M30s.tail(&s.M30s, 0)
		p.Days.tail(&s.Days, 0)
		p.Weeks.tail(&s.Weeks, 0)
		p.Months.tail(&s.Months, 0)
	} else {
		// tail
		p.M1s.tail(&s.M1s, 240)
		p.M5s.tail(&s.M5s, 60)
		p.M30s.tail(&s.M30s, 8)
		p.Days.tail(&s.Days, 8)
		p.Weeks.tail(&s.Weeks, 8)
		p.Months.tail(&s.Months, 8)
	}
	return json.Marshal(s)
}

func NewStock(id string, hub_height int) *Stock {
	p := &Stock{
		Id:    id,
		hash:  StockHash(id),
		count: 1,
	}

	p.M1s.Init(hub_height, id+" f1", nil, &p.M5s)
	p.M5s.Init(hub_height, id+" f5", &p.M1s, &p.M30s)
	p.M30s.Init(hub_height, id+" f30", &p.M5s, &p.Days)
	p.Days.Init(hub_height, id+" day", &p.M30s, &p.Weeks)
	p.Weeks.Init(hub_height, id+" week", &p.Days, &p.Months)
	p.Months.Init(hub_height, id+" month", &p.Weeks, nil)

	return p
}

type Stocks struct {
	stocks  PStockSlice
	rwmutex sync.RWMutex
	store   store.Store
	play    int
	ch      chan *Stock

	min_hub_height int
}

func NewStocks(storestr string, play, min_hub_height int) *Stocks {
	store := store.Get(storestr)
	if min_hub_height < 0 {
		min_hub_height = 0
	}
	return &Stocks{
		min_hub_height: min_hub_height,
		store:          store,
		play:           play,
	}
}

func (p *Stocks) Store() store.Store { return p.store }

func (p *Stocks) Run() {
	if p.play > minPlay {
		for {
			p.play_next_tick()
			time.Sleep(time.Duration(p.play) * time.Millisecond)
		}
	}

	robot.Work()
	for {
		if IsTradeTime(time.Now()) {
			p.Ticks_update_real()
		}
		time.Sleep(tickPeriod)
	}
}

func (p *Stocks) Chan(ch chan *Stock) {
	p.ch = ch
}

func (p *Stocks) res(stock *Stock) {
	if p.ch != nil {
		p.ch <- stock
	}
}

func (p *Stocks) update(s *Stock) {
	if s.Update(p.store, p.play > minPlay) {
		p.res(s)
	}
}

func (p *Stocks) Insert(id string) (int, *Stock, bool) {
	p.rwmutex.RLock()
	i, ok := p.stocks.Search(id)
	if ok {
		s := p.stocks[i]
		p.rwmutex.RUnlock()
		if atomic.AddInt32(&s.count, 1) < 1 {
			atomic.StoreInt32(&s.count, 1)
		}
		return i, s, false
	}

	s := NewStock(id, p.min_hub_height)

	p.rwmutex.RUnlock()
	p.rwmutex.Lock()
	defer p.rwmutex.Unlock()

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
	p.rwmutex.RLock()
	defer p.rwmutex.RUnlock()
	if i, ok := p.stocks.Search(id); ok {
		atomic.AddInt32(&p.stocks[i].count, -1)
	}
}

func (p *Stocks) Watch(id string) (*Stock, bool) {
	i, s, isnew := p.Insert(id)
	if isnew {
		go p.update(s)
		glog.V(LogV).Infof("watch new stock id=%s index=%d", id, i)
	} else {
		glog.V(LogV).Infof("watch stock id=%s index=%d count=%d", id, i, s.count)
	}
	return s, isnew
}

func (p *Stocks) UnWatch(id string) {
	p.Remove(id)
}

func (p *Stocks) Find_need_update_tick_ids() (pstocks PStockSlice) {
	p.rwmutex.RLock()
	defer p.rwmutex.RUnlock()
	for i, l := 0, len(p.stocks); i < l; i++ {
		if atomic.LoadInt32(&p.stocks[i].loaded) < 2 {
			continue
		}
		pstocks = append(pstocks, p.stocks[i])
	}
	return
}

func (p *Stocks) play_next_tick() {
	p.rwmutex.RLock()
	defer p.rwmutex.RUnlock()
	for i, l := 0, len(p.stocks); i < l; i++ {
		if atomic.LoadInt32(&p.stocks[i].loaded) < 2 {
			continue
		}
		if atomic.LoadInt32(&p.stocks[i].count) < 1 {
			continue
		}

		p.stocks[i].rw.Lock()
		if p.stocks[i].Ticks.play == nil || len(p.stocks[i].Ticks.play) < 1 {
			p.stocks[i].Ticks.play = p.stocks[i].Ticks.Data
			p.stocks[i].Ticks.Data = []Tick{}
			if len(p.stocks[i].Ticks.play) > 240 {
				p.stocks[i].Ticks.Data = p.stocks[i].Ticks.play[:240]
			}
		}
		lplay := len(p.stocks[i].Ticks.play)
		ldata := len(p.stocks[i].Ticks.Data)
		if ldata < lplay {
			p.stocks[i].Ticks.Data = p.stocks[i].Ticks.play[:ldata+1]
			p.stocks[i].Merge(false, p.store)
			p.res(p.stocks[i])
		}
		p.stocks[i].rw.Unlock()
	}
}

func (p *Stocks) Ticks_update_real() {
	var wg sync.WaitGroup

	stocks := p.Find_need_update_tick_ids()
	l := len(stocks)
	if l < 1 {
		return
	}

	for i := 0; i < l; {
		var b bytes.Buffer
		var pstocks PStockSlice
		step := 50
		if i+step < l {
			pstocks = stocks[i : i+step]
		} else {
			pstocks = stocks[i:l]
		}
		for j := 0; j < step && i < l; i, j = i+1, j+1 {
			if b.Len() > 0 {
				b.WriteString(",")
			}
			b.WriteString(stocks[i].Id)
		}
		if b.Len() < 1 {
			continue
		}

		wg.Add(1)
		go func(ids string, pstocks PStockSlice) {
			defer wg.Done()
			body := robot.Tick_download_real_from_sina(ids)
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
						pstocks[idx].Merge(false, p.store)
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

func (p *Stock) Merge(day bool, store store.Store) {
	m1_fresh_index := p.Ticks2M1s()
	m5_fresh_index := p.M5s.MergeFrom(&p.M1s, false, Minute5end)
	m30_fresh_index := p.M30s.MergeFrom(&p.M1s, false, Minute30end)

	if day {
		p.Ticks.clean()
		td, _ := store.LoadMacd(p.Id, L1, p.M1s.start)
		p.M1s.Macd(m1_fresh_index, td)
		store.SaveMacds(p.Id, L1, p.M1s.Data)
		td, _ = store.LoadMacd(p.Id, L5, p.M1s.start)
		p.M5s.Macd(m5_fresh_index, td)
		store.SaveMacds(p.Id, L5, p.M5s.Data)
		td, _ = store.LoadMacd(p.Id, L30, p.M1s.start)
		p.M30s.Macd(m30_fresh_index, td)
		store.SaveMacds(p.Id, L30, p.M30s.Data)
	} else {
		p.M1s.Macd(m1_fresh_index, nil)
		p.M5s.Macd(m5_fresh_index, nil)
		p.M30s.Macd(m30_fresh_index, nil)
	}

	p.M1s.ParseChan()
	p.M5s.ParseChan()
	p.M30s.ParseChan()

	if day {
		p.Weeks.MergeFrom(&p.Days, true, Weekend)
		p.Months.MergeFrom(&p.Days, true, Monthend)
		td, _ := store.LoadMacd(p.Id, LDay, p.Days.start)
		p.Days.Macd(0, td)
		store.SaveMacds(p.Id, LDay, p.Days.Data)
		td, _ = store.LoadMacd(p.Id, LWeek, p.Days.start)
		p.Weeks.Macd(0, td)
		store.SaveMacds(p.Id, LWeek, p.Weeks.Data)
		td, _ = store.LoadMacd(p.Id, LMonth, p.Days.start)
		p.Months.Macd(0, td)
		store.SaveMacds(p.Id, LMonth, p.Months.Data)
		p.Days.ParseChan()
		p.Weeks.ParseChan()
		p.Months.ParseChan()
	}
}

func (p *Tdatas) ParseChan() {
	p.ParseTyping()

	if p.base == nil {
		p.Typing.LinkTyping()
		p.ParseSegment()
		p.Segment.LinkTyping()
	}
	p.ParseHub()
	p.LinkHub()
}

func (p *Stock) Update(store store.Store, play bool) bool {
	if !atomic.CompareAndSwapInt32(&p.loaded, 0, 1) {
		return false
	}

	p.Days_update(store)

	p.Ticks_update(store)
	p.Ticks_today_update()

	if play {
		glog.Warningln("WITH PLAY MODE")
	} else {
		p.Merge(true, store)
	}
	atomic.StoreInt32(&p.loaded, 2)
	return true
}

func (p *Stock) days_download(t time.Time) ([]int, error) {
	inds := []int{}
	tds, err := robot.Days_download(p.Id, t)
	if err != nil {
		return inds, err
	}
	for i, count := 0, len(tds); i < count; i++ {
		ind, isnew := p.Days.Add(tds[i])
		if isnew || ind > 0 {
			inds = append(inds, ind)
		}
	}
	return inds, nil
}

func (p *Stock) Days_update(store store.Store) int {
	c := Day_collection_name(p.Id)
	p.Days.start = store.GetStartTime(p.Id, LDay)
	p.Days.Data, _ = store.LoadTDatas(c, p.Days.start)
	t := p.Days.latest_time()
	now := time.Now().AddDate(0, 0, -1).UTC().Truncate(time.Hour * 24)
	for !IsTradeDay(now) {
		now = now.AddDate(0, 0, -1)
	}
	if t.Equal(now) || t.After(now) {
		return 0
	}

	inds, _ := p.days_download(t)
	if len(inds) > 0 {
		store.SaveTDatas(c, p.Days.Data, inds)
		factor := p.Days.Factor()
		store.UpdateFactor(p.Id, factor)
	}
	return len(inds)
}

func (p *Stock) Ticks_update(store store.Store) int {
	c := Tick_collection_name(p.Id)
	p.M1s.start = store.GetStartTime(p.Id, L1)
	p.Ticks.Data, _ = store.LoadTicks(c, p.M1s.start)
	begin_time := p.M1s.start
	l := len(p.Ticks.Data)
	if l > 0 {
		begin_time = p.Ticks.Data[0].Time
	}

	now := time.Now().UTC()
	end_time := now.Truncate(time.Hour * 24)
	if now.Hour() > 10 {
		end_time = end_time.AddDate(0, 0, 1)
	}

	if begin_time.Equal(market_begin_day) {
		begin_time = end_time.AddDate(0, -2, -1)
	}
	begin_time = begin_time.AddDate(0, 0, 1).Truncate(time.Hour * 24)
	daylen := len(p.Days.Data)
	if daylen < 1 {
		return 0
	}
	i, _ := ((TdataSlice)(p.Days.Data)).Search(begin_time)
	glog.V(LogV).Infof("from %d/%d %s begin_time=%s end_time=%s", i, daylen, p.M1s.start, begin_time, end_time)
	var t time.Time
	for ; i <= daylen; i++ {
		if i < daylen {
			t = p.Days.Data[i].Time
		} else if i == daylen {
			t = p.Days.Data[i-1].Time.AddDate(0, 0, 1)
		}
		if !end_time.After(t) {
			glog.V(LogV).Infoln(t, "reach end_time", end_time)
			break
		}

		if p.Ticks.hasTimeData(t) {
			continue
		}

		glog.V(LogV).Infoln("prepare download ticks", t)
		if ticks, err := p.ticks_download(t); ticks != nil {
			for j, _ := range ticks {
				p.Ticks.Add(ticks[j])
			}
			store.SaveTicks(c, ticks)
			glog.V(LogV).Infoln("download ticks succ", t)
		} else if err != nil {
			glog.V(LogD).Infoln("download ticks err", err)
		}
	}

	count := len(p.Ticks.Data)
	glog.V(LogV).Infof("download ticks %d/%d", count-l, count)
	return count - l
}

/*
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
*/

var UnknowSinaRes error = errors.New("could not find '成交时间' in head line")

func (p *Stock) ticks_download(t time.Time) ([]Tick, error) {
	body := robot.Tick_download_from_sina(p.Id, t)
	if body == nil {
		return nil, UnknowSinaRes
	}
	body = bytes.TrimSpace(body)
	lines := bytes.Split(body, []byte("\n"))
	count := len(lines) - 1
	if count < 1 {
		return nil, UnknowSinaRes
	}
	if bytes.Contains(lines[0], []byte("script")) {
		return nil, UnknowSinaRes
	}
	if !bytes.Contains(lines[0], []byte("成交时间")) {
		return nil, UnknowSinaRes
	}

	ticks := make([]Tick, count)
	for i := count; i > 0; i-- {
		line := bytes.TrimSpace(lines[i])
		infos := bytes.Split(line, []byte("\t"))
		if len(infos) != 6 {
			err := errors.New("could not parse line " + string(line))
			return nil, err
		}
		ticks[count-i].FromString(t, infos[0], infos[1], infos[2],
			infos[3], infos[4], infos[5])
	}
	FixTickTime(ticks)

	return ticks, nil
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
	return count - l
}

func (p *Stock) ticks_get_today() bool {
	last_t, name, err := Tick_get_today_date(p.Id)
	if err != nil {
		glog.Warningln("get today date fail", err)
		return false
	}
	p.Name = name
	t := time.Now().UTC().Truncate(time.Hour * 24)
	if t.After(last_t) {
		return false
	}

	body := robot.Tick_download_today_from_sina(p.Id)
	if body == nil {
		return false
	}
	body = bytes.TrimSpace(body)
	lines := bytes.Split(body, []byte("\n"))

	ticks := []Tick{}
	tick := Tick{}
	nul := []byte("")
	for i := len(lines) - 1; i > 0; i-- {
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

		tick.FromString(t, infos[0], infos[2], nul, infos[1], nul, infos[3])
		if tick.Volume == 0 && tick.Price == 0 {
			continue
		}
		ticks = append(ticks, tick)
	}
	FixTickTime(ticks)
	FixTickData(ticks)

	for _, tick := range ticks {
		p.Ticks.Add(tick)
	}
	return true
}

func (p *Stock) tick_get_real(line []byte) bool {
	infos := bytes.Split(line, []byte(","))
	if len(infos) < 33 {
		glog.Warningln("sina hq api, res format changed")
		return false
	}

	p.Name = string(infos[0])
	nul := []byte("")
	tick := RealtimeTick{}
	t, _ := time.Parse("2006-01-02", string(infos[30]))
	tick.FromString(t, infos[31], infos[3], nul, infos[8], infos[9], nul)
	tick.Buyone = ParseCent(string(infos[11]))
	tick.Sellone = ParseCent(string(infos[21]))
	tick.SetStatus(infos[32])

	if p.last_tick.Volume == 0 {
		p.last_tick = tick
		if tick.Time.Before(p.lst_trade) {
			p.last_tick.Volume = 0
		}
		return false
	}
	if tick.Volume != p.last_tick.Volume {
		if tick.Price >= p.last_tick.Sellone {
			tick.Type = Buy_tick
		} else if tick.Price <= p.last_tick.Buyone {
			tick.Type = Sell_tick
		} else {
			tick.Type = Eq_tick
		}
		tick.Change = tick.Price - p.last_tick.Price

		volume := (tick.Volume - p.last_tick.Volume) / 100
		p.last_tick = tick
		tick.Volume = volume
		p.Ticks.Add(tick.Tick)
		p.lst_trade = tick.Time
		return true
	}
	return false
}
