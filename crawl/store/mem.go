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

func (p *Mem) LoadTDatas(table string, start time.Time) (res []Tdata, err error) {
	return
}

func (p *Mem) SaveTDatas(string, []Tdata, []int) (err error) { return }

func (p *Mem) LoadTicks(table string, start time.Time) (res []Tick, err error) {
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

func (p *Mem) LoadStar(uid int) (res []CategoryItemInfo, err error) {
	return
}

func (p *Mem) GetStartTime(symbol string, typ int) time.Time {
	return Market_begin_day
}

func (p *Mem) LoadMacd(symbol string, typ int, start time.Time) (*Tdata, error) {
	return nil, nil
}

func (p *Mem) SaveMacds(symbol string, typ int, datas []Tdata) error {
	return nil
}

func (p *Mem) UpdateFactor(name string, factor int) {}
