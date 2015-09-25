package crawl

import "testing"

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
