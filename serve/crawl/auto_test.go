package crawl

import (
	"sort"
	"testing"
)

type test_stock_hash_pair struct {
	id   string
	hash int
}

var stock_hash_tests = []test_stock_hash_pair{
	{"sh600000", 600000},
	{"sz300000", 300000},
	{"sh000001", 1},
	{"sz000001", 1},
	{"s_sz000001", 1},
}

func TestStockHash(t *testing.T) {
	for _, pair := range stock_hash_tests {
		v := StockHash(pair.id)
		if v != pair.hash {
			t.Error(
				"For", pair.id,
				"expected", pair.hash,
				"got", v,
			)
		}
	}
}

var pstockslice = PStockSlice{
	&Stock{Id: "sh600000", hash: 600000},
	&Stock{Id: "sh600001", hash: 600001},
	&Stock{Id: "sh600003", hash: 600003},
	&Stock{Id: "sh600004", hash: 600004},
}

type test_stock_slice_pair struct {
	id  string
	idx int
	has bool
}

var stock_slice_tests = []test_stock_slice_pair{
	{"sh600000", 0, true},
	{"sh600003", 2, true},
	{"sh600002", 2, false},
	{"sh600005", 4, false},
	{"sh600006", 4, false},
	{"sh500000", 0, false},
}

func TestSearchPStockSlice(t *testing.T) {
	if !sort.IsSorted(pstockslice) {
		t.Error("test data should be sorted")
	}
	for _, pair := range stock_slice_tests {
		if v, ok := pstockslice.Search(pair.id); v != pair.idx || ok != pair.has {
			t.Error(
				"For", pair.id,
				"expected", pair.idx, pair.has,
				"got", v, ok,
			)
		}
	}
}

func TestStocksInsert(t *testing.T) {
  s := Stocks{stocks: pstockslice}
  i := s.Insert("sh600003")
  if i != 2 {
    t.Error(
      "For", "sh600003",
      "expected", 2,
      "got", i,
    )
  }
  i = s.Insert("sh600002")
  if i != 2 || s.stocks[2].Id == "sh600002" {
    t.Error(
      "For", "sh600002",
      "expected", 2,
      "got", i,
    )
  }
}
