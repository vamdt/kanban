package crawl

import (
	"sync"
	"sync/atomic"
	"time"

	. "./base"
	"./robot"
	"./robot/sina"
	"./store"
	"github.com/golang/glog"
)

func LoadCategoryItem(p *CategoryItem, store store.Store) {
	data, err := store.LoadCategories()
	if err != nil {
		glog.Warningln("load categories err", err)
	}

	if len(data) < 1 {
		return
	}

	p.Assembly(data)
}

func UpdateCate(storestr string) {
	store := store.Get(storestr)
	cate := NewCategoryItem("")
	LoadCategoryItem(cate, store)
	if cate.Sub == nil {
		cate.Sub = *NewCategory()
	}
	robot := sina.SinaRobot{}
	robot.Cate(cate.Sub)
	store.SaveCategories(cate.Sub, cate.Id)
}

func handle_realtimetick(stock *Stock, wg *sync.WaitGroup) func(interface{}, bool) bool {
	return func(t interface{}, ok bool) bool {
		defer wg.Done()
		if !ok || t == nil {
			glog.Infoln("RealtimeTick fail", stock.Id)
			return true
		}

		rt := t.(*RealtimeTick)
		td := Tdata{}
		td.Time = rt.Time
		td.Open = rt.Change
		td.Close = rt.Price
		td.Volume = rt.Volume / 100
		td.HL = rt.HL
		stock.Days.Add(td)
		stock.Name = rt.Name
		return true
	}
}

func (p *Stocks) Days_update_real() {
	now := time.Now().UTC()
	if !IsTradeDay(now) {
		return
	}

	h, m, _ := now.Clock()
	if h < 1 || (h == 1 && m < 30) {
		return
	}

	var wg sync.WaitGroup

	p.rwmutex.RLock()
	defer p.rwmutex.RUnlock()
	for i, l := 0, len(p.stocks); i < l; i++ {
		if atomic.LoadInt32(&p.stocks[i].loaded) < 2 {
			continue
		}
		wg.Add(1)
		robot.GetRealtimeTick(p.stocks[i].Id, handle_realtimetick(p.stocks[i], &wg))
	}

	glog.V(LogV).Infoln("wait Days_update_real")
	wg.Wait()
	glog.V(LogV).Infoln("Days_update_real done")
}

func UpdateFactor(storestr string) {
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
	latest_time := market_begin_day
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
			s.Days_update(store)
			if t := s.Days.latest_time(); t.After(latest_time) {
				latest_time = t
			}
			s.loaded = int32(i) + 2
		}(s, i)
	}

	glog.Infoln("wait all update done")
	wg.Wait()
	glog.Infoln("all update done")

	stocks.Days_update_real()

	stocks.rwmutex.RLock()
	for i, l := 0, len(stocks.stocks); i < l; i++ {
		if atomic.LoadInt32(&stocks.stocks[i].loaded) < 2 {
			continue
		}
		s := stocks.stocks[i]
		if t := s.Days.latest_time(); t.Before(latest_time) {
			glog.Infoln("latest time", s.Id, t)
			continue
		}
		j := int(s.loaded) - 2
		if j < len(data) {
			data[j].Factor = s.Days.Factor()
			if s.Name != data[j].Tag && len(s.Name) > 0 {
				data[j].Tag = s.Name
			}
		}
	}
	stocks.rwmutex.RUnlock()

	factor := make(map[string]int)
	stats := make([]int, 10)
	for i, info := range data {
		if info.Leaf && info.Factor > 0 {
			factor[info.Name] = i
			if info.Factor > -1 && info.Factor < 10 {
				if IsChinaShareCode(info.Name) {
					stats[info.Factor]++
				}
			}
		}
	}
	mostFactor := maxIndex(stats)

	for i, info := range data {
		if info.Leaf && info.Factor == 0 {
			if f, ok := factor[info.Name]; ok {
				data[i].Factor = data[f].Factor
				data[i].Tag = data[f].Tag
			}
		}
	}

	for i, info := range data {
		if info.Leaf {
			continue
		}
		pid := info.Id
		num := 0
		factor := 0

		for _, item := range data {
			if item.Leaf && item.Pid == pid {
				if item.Factor == 0 {
					continue
				}
				num++
				factor += item.Factor
			}
		}
		if num > 0 {
			data[i].Factor = factor / num
		} else {
			data[i].Factor = mostFactor
		}
	}

	store.SaveCategoryItemInfoFactor(data)
}
