package crawl

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
)

type CategoryItemInfo struct {
	Id     int
	Pid    int
	Factor int
	Leaf   bool
	Name   string
}

type CategoryItem struct {
	Id   int
	Name string
	Info []CategoryItemInfo
	Sub  Category
}

type Category map[string]CategoryItem

func NewCategory() *Category {
	c := make(Category)
	return &c
}

func NewCategoryItem(name string) *CategoryItem {
	return &CategoryItem{Name: name}
}

func (p *CategoryItem) Assembly(data []CategoryItemInfo) {
	for i := len(data) - 1; i > -1; i-- {
		if p.Id != data[i].Pid {
			continue
		}

		name := data[i].Name

		if data[i].Leaf {
			p.Info = append(p.Info, data[i])
		} else {
			if p.Sub == nil {
				p.Sub = *NewCategory()
			}

			if _, ok := p.Sub[name]; !ok {
				item := NewCategoryItem(name)
				item.Id = data[i].Id
				p.Sub[name] = *item
			}
			item := p.Sub[name]
			item.Assembly(data)
			p.Sub[name] = item
		}
	}
}

func (p *CategoryItem) AddStock(id string) {
	info := CategoryItemInfo{Name: id}
	p.Info = append(p.Info, info)
}

func (p *CategoryItem) LeafCount() int {
	return len(p.Info)
}

func LoadCategories(store Store) Category {
	data, err := store.LoadCategories()
	if err != nil {
		glog.Infoln("load categories err", err)
	}

	if len(data) < 1 {
		return nil
	}

	item := NewCategoryItem("")
	item.Assembly(data)
	return item.Sub
}

func UpdateCate(storestr string) {
	store := getStore(storestr)
	tc := LoadCategories(store)
	if tc == nil {
		tc = *NewCategory()
	}
	robot := SinaRobot{}
	robot.Cate(tc)
	store.SaveCategories(tc)
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
					pstocks[idx].day_get_real(info[1])
				}
			}
		}(b.String(), pstocks)

	}
	glog.V(LogV).Infoln("wait Days_update_real")
	wg.Wait()
	glog.V(LogV).Infoln("Days_update_real done")
}

func UpdateFactor(storestr string) {
	store := getStore(storestr)
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

	stocks := Stocks{store: store}
	var wg sync.WaitGroup
	c := 0
	for i, _ := range data {

		if c == 50 {
			wg.Wait()
			c = 0
		}
		if !data[i].Leaf {
			continue
		}
		_, s, ok := stocks.Insert(data[i].Name)
		if !ok {
			continue
		}

		wg.Add(1)
		c++
		go func(s *Stock, i int) {
			defer wg.Done()
			s.Days_update(store)
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
		j := int(s.loaded) - 2
		if j < len(data) {
			data[j].Factor = s.Days.Factor()
		}
	}
	stocks.rwmutex.RUnlock()

	factor := make(map[string]int)
	for _, info := range data {
		if info.Leaf && info.Factor > 0 {
			factor[info.Name] = info.Factor
		}
	}

	for i, info := range data {
		if info.Leaf && info.Factor == 0 {
			if f, ok := factor[info.Name]; ok {
				data[i].Factor = f
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
					glog.Warningln("found an 0 factor", item)
					continue
				}
				num++
				factor += item.Factor
			}
		}
		if num > 0 {
			data[i].Factor = factor / num
		}
	}

	store.SaveCategoryItemInfoFactor(data)
}
