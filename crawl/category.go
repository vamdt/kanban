package crawl

import "github.com/golang/glog"

type CategoryItem struct {
	Name string
	Id   []string
}

type Category map[string]CategoryItem
type TopCategory map[string]Category

func NewTopCategory() *TopCategory {
	c := make(TopCategory)
	return &c
}

func NewCategory() *Category {
	c := make(Category)
	return &c
}

func NewCategoryItem(name string) *CategoryItem {
	return &CategoryItem{Name: name}
}

func UpdateCate(storestr string) {
	store := getStore(storestr)
	tc, err := store.LoadCategories()
	if err != nil {
		glog.Infoln("load categories err", err)
	}
	if tc == nil {
		tc = *NewTopCategory()
	}
	robot := SinaRobot{}
	robot.Cate(tc)
	store.SaveCategories(tc)
}
