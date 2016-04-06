package store

import (
	"sync"
	"time"

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

type TicksStore interface {
	LoadTicks(table string, start time.Time) ([]Tick, error)
	SaveTicks(string, []Tick) error
	HasTickData(table string, t time.Time) bool
}

type StarStore interface {
	LoadCategories() ([]CategoryItemInfo, error)
	SaveCategories(Category, int) error
	SaveCategoryItemInfoFactor([]CategoryItemInfo)
	Lucky(pid int, symbol string) string
	GetSymbolName(symbol string) string
	Star(int, string)
	UnStar(int, string)
	IsStar(pid int, symbol string) bool
	LoadStar(uid int) ([]CategoryItemInfo, error)
}

type TDataStore interface {
	LoadTDatas(table string, start time.Time) ([]Tdata, error)
	SaveTDatas(string, []Tdata) error
	GetStartTime(symbol string, typ int) time.Time
	LoadMacd(symbol string, typ int, start time.Time) (*Tdata, error)
	SaveMacds(symbol string, typ int, datas []Tdata) error
}

type Store interface {
	TicksStore
	StarStore
	TDataStore
	Open() error
	Close()
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
