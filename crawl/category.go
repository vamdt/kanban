package crawl

import "github.com/golang/glog"

type CategoryItem struct {
	Id   int
	Name string
	Sid  *[]string
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

func (p *CategoryItem) AddStock(id string) {
	if p.Sid == nil {
		p.Sid = &[]string{}
	}
	sid := *p.Sid
	sid = append(sid, id)
	p.Sid = &sid
}

func (p *CategoryItem) LeafCount() int {
	if p.Sid == nil {
		return 0
	}
	return len(*p.Sid)
}

func UpdateCate(storestr string) {
	store := getStore(storestr)
	tc, err := store.LoadCategories()
	if err != nil {
		glog.Infoln("load categories err", err)
	}
	if tc == nil {
		tc = *NewCategory()
	}
	robot := SinaRobot{}
	robot.Cate(tc)
	store.SaveCategories(tc)
}

func UpdateFactor(storestr string) {
	store := getStore(storestr)
	tc, err := store.LoadCategories()
	if err != nil {
		glog.Infoln("load categories err", err)
		return
	}
	if tc == nil {
		glog.Infoln("load categories empty")
		return
	}

	for _, v := range tc {
		glog.Infoln("top ", v.Name, v.Id)
		if v.Sub != nil {
			for name, item := range v.Sub {
				glog.Infoln(">> ", name)
				if item.Sid != nil {
					for _, id := range *item.Sid {
						glog.Infoln("\t\t>>> ", id)
					}
				}
			}
		}
	}
}
