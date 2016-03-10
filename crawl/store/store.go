package store

import (
	"sync"

	. "../base"
	"github.com/golang/glog"
)

var (
	storesMu sync.Mutex
	stores   = make(map[string]Store)
)

// Register makes a store driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, store Store) {
	storesMu.Lock()
	defer storesMu.Unlock()
	if store == nil {
		panic("store: Register driver is nil")
	}
	if _, dup := stores[name]; dup {
		panic("store: Register called twice for driver " + name)
	}
	stores[name] = store
}

func unregisterAllDrivers() {
	storesMu.Lock()
	defer storesMu.Unlock()
	// For tests.
	stores = make(map[string]Store)
}

type Store interface {
	Open() error
	Close()
	LoadTDatas(table string) ([]Tdata, error)
	SaveTDatas(string, []Tdata) error
	LoadTicks(table string) ([]Tick, error)
	SaveTicks(string, []Tick) error
	LoadCategories() ([]CategoryItemInfo, error)
	SaveCategories(Category, int) error
	SaveCategoryItemInfoFactor([]CategoryItemInfo)
	Star(int, string)
	UnStar(int, string)
}

func Get(s string) Store {
	storesMu.Lock()
	store, ok := stores[s]
	storesMu.Unlock()

	if !ok {
		glog.Fatalf("store: unknown store %q (forgotten import?)", s)
		return nil
	}

	err := store.Open()
	if err != nil {
		glog.Fatalln("store: open[", s, "] fail", err)
	}
	return store
}
