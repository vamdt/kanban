package crawl

func NewMemStore() (ms *MemStore, err error) {
	ms = &MemStore{}
	return
}

type MemStore struct {
}

func (p *MemStore) Close() {
}

func (p *MemStore) LoadTDatas(table string) (res []Tdata, err error) {
	return
}

func (p *MemStore) SaveTData(table string, data *Tdata) (err error) {
	return
}

func (p *MemStore) LoadTicks(table string) (res []Tick, err error) {
	return
}

func (p *MemStore) SaveTick(table string, tick *Tick) (err error) {
	return
}
