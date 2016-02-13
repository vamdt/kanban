package crawl

import "github.com/golang/glog"

type Store interface {
	Close()
	LoadTDatas(table string) ([]Tdata, error)
	SaveTData(table string, data *Tdata) error
	LoadTicks(table string) ([]Tick, error)
	SaveTick(table string, tick *Tick) error
	LoadCategories() (TopCategory, error)
	SaveCategories(TopCategory) error
}

func getStore(s string) Store {
	var store Store
	var err error
	if s == "mongo" {
		store, err = NewMongoStore()
	} else if s == "mysql" {
		store, err = NewMysqlStore()
	} else {
		store, err = NewMemStore()
	}
	if err != nil {
		glog.Fatalln("new [", s, "] store", err)
	}
	return store
}
