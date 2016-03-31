package crawl

import (
	"sync"
	"sync/atomic"

	. "./base"
	"./robot"
	"./store"
	"github.com/golang/glog"
)

func (p *Stock) Days_fix(store store.Store) {
	c := Day_collection_name(p.Id)
	p.Days.Data, _ = store.LoadTDatas(c)
	l := len(p.Days.Data)
	if l < 1 {
		return
	}
	t := p.Days.Data[0].Time
	inds, _ := p.days_download(t)
	tdatas := []Tdata{}
	l = len(p.Days.Data)
	for _, ind := range inds {
		if ind < l {
			tdatas = append(tdatas, p.Days.Data[ind])
		}
	}
	if len(tdatas) > 0 {
		store.SaveTDatas(c, tdatas)
	}
}

func (p *Stock) Ticks_fix(store store.Store) {
	daylen := len(p.Days.Data)
	if daylen < 1 {
		return
	}
	c := Tick_collection_name(p.Id)

	for i := daylen - 1; i > -1; i-- {
		if daylen-i > 60 {
			break
		}
		t := p.Days.Data[i].Time
		if store.HasTickData(c, t) {
			glog.V(LogV).Infoln(t, "already in db, skip")
			continue
		}

		p.Ticks.Data = []Tick{}
		if _, err := p.ticks_download(t); err != nil {
			glog.Warningln("fix ticks err", err)
		}
		if len(p.Ticks.Data) < 1 {
			glog.Warningln("got empty ticks")
			continue
		}
		store.SaveTicks(c, p.Ticks.Data)
	}
}

func FixData(storestr string) {
	store := store.Get(storestr)
	data, err := store.LoadCategories()
	if err != nil {
		glog.Infoln("load categories err", err)
	}

	if len(data) < 1 {
		glog.Infoln("load categories empty")
		return
	}

	for i, _ := range data {
		data[i].Factor = 0
	}

	robot.Work()
	stocks := Stocks{store: store}
	var wg sync.WaitGroup
	for i, _ := range data {

		if !data[i].Leaf {
			continue
		}
		_, s, ok := stocks.Insert(data[i].Name)
		if !ok {
			continue
		}

		wg.Add(1)
		go func(s *Stock, i int) {
			defer wg.Done()
			s.Days_fix(store)
			s.loaded = int32(i) + 2
		}(s, i)
	}

	glog.Infoln("wait all fix done")
	wg.Wait()
	glog.Infoln("all fix done")

	stocks.rwmutex.RLock()
	for i, l := 0, len(stocks.stocks); i < l; i++ {
		s := stocks.stocks[i]
		if atomic.LoadInt32(&s.loaded) < 2 {
			continue
		}
		s.Ticks_fix(store)
	}
	stocks.rwmutex.RUnlock()
}
