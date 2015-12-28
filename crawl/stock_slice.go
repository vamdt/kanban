package crawl

import (
	"sort"
	"strings"
)

type PStockSlice []*Stock

func (p PStockSlice) Len() int      { return len(p) }
func (p PStockSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PStockSlice) Less(i, j int) bool {
	if p[i].hash == p[j].hash {
		return strings.Compare(p[i].Id, p[j].Id) == -1
	}
	return p[i].hash < p[j].hash
}

func SearchPStockSlice(a PStockSlice, id string) int {
	hash := StockHash(id)
	return sort.Search(len(a), func(i int) bool {
		if a[i].hash == hash {
			return strings.Compare(a[i].Id, id) > -1
		}
		return a[i].hash > hash
	})
}

func (p PStockSlice) Search(id string) (int, bool) {
	i := SearchPStockSlice(p, id)
	if i >= p.Len() || i < 0 {
		return i, false
	}
	if strings.Compare(p[i].Id, id) == 0 {
		return i, true
	}
	return i, false
}
