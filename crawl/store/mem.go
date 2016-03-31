package store

import (
	"time"

	. "../base"
)

func init() {
	Register("mem", &Mem{})
}

type Mem struct {
}

func (p *Mem) Open() error { return nil }

func (p *Mem) Close() {
}

func (p *Mem) LoadTDatas(table string) (res []Tdata, err error) {
	return
}

func (p *Mem) SaveTDatas(string, []Tdata) (err error) { return }

func (p *Mem) LoadTicks(table string) (res []Tick, err error) {
	return
}

func (p *Mem) SaveTicks(string, []Tick) (err error) { return }

func (p *Mem) LoadCategories() (res []CategoryItemInfo, err error) { return }

func (p *Mem) SaveCategories(Category, int) (err error) { return }

func (p *Mem) SaveCategoryItemInfoFactor([]CategoryItemInfo) {}

func (p *Mem) Star(pid int, symbol string) {}

func (p *Mem) UnStar(pid int, symbol string) {}

func (p *Mem) IsStar(pid int, symbol string) bool {
	return false
}

func (p *Mem) Lucky(pid int, symbol string) string {
	return symbol
}

func (p *Mem) GetSymbolName(symbol string) string {
	return symbol
}

func (p *Mem) HasTickData(table string, t time.Time) bool {
	return true
}
