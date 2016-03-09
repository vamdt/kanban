package crawl

func init() {
	RegisterStore("mem", &MemStore{})
}

type MemStore struct {
}

func (p *MemStore) Open() error { return nil }

func (p *MemStore) Close() {
}

func (p *MemStore) LoadTDatas(table string) (res []Tdata, err error) {
	return
}

func (p *MemStore) SaveTDatas(string, []Tdata) (err error) { return }

func (p *MemStore) LoadTicks(table string) (res []Tick, err error) {
	return
}

func (p *MemStore) SaveTicks(string, []Tick) (err error) { return }

func (p *MemStore) LoadCategories() (res []CategoryItemInfo, err error) { return }

func (p *MemStore) SaveCategories(Category, int) (err error) { return }

func (p *MemStore) SaveCategoryItemInfoFactor([]CategoryItemInfo) {}
