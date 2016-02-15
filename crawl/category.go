package crawl

import (
	"sync"

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
	for i, info := range data {
		if !info.Leaf {
			continue
		}
		_, s, ok := stocks.Insert(info.Name)
		if !ok {
			continue
		}

		wg.Add(1)
		go func(s *Stock, i int) {
			defer wg.Done()
			s.Days_update(store)
			data[i].Factor = s.Days.Factor()
		}(s, i)
	}
	wg.Wait()

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
