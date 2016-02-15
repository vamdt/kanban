package crawl

import "github.com/golang/glog"

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
	tc := LoadCategories(store)
	if tc == nil {
		glog.Infoln("load categories empty")
		return
	}

	for _, v := range tc {
		glog.Infoln("top ", v.Name, v.Id)
		if v.Sub != nil {
			for name, item := range v.Sub {
				glog.Infoln(">> ", name)
				if item.Info != nil {
					for _, info := range item.Info {
						glog.Infoln("\t\t>>> ", info)
					}
				}
			}
		}
	}
}
