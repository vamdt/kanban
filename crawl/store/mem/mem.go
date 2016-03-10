package mem

import (
	. "../../base"
	"../../store"
)

func init() {
	store.Register("mem", &Mem{})
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
