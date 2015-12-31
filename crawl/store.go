package crawl

import "time"

type Store interface {
	Close()
	LoadTDatas(table string) ([]Tdata, error)
	SaveTData(table string, data *Tdata) error
	LoadTicks(table string) ([]Tick, error)
	SaveTick(table string, tick *Tick) error
	TickHasTimeData(table string, t time.Time) bool
}
