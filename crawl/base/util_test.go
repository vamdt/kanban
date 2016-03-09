package base

import "testing"

type test_cent_pair struct {
	str string
	exp int
}

var cent_tests = []test_cent_pair{
	{"10.00", 1000},
	{"10.01", 1001},
	{"-10.00", -1000},
	{"-10.01", -1001},
	{"0.01", 1},
	{"-0.01", -1},
	{"1", 100},
	{"", 0},
}

func TestParseCent(t *testing.T) {
	for _, pair := range cent_tests {
		v := ParseCent(pair.str)
		if v != pair.exp {
			t.Error(
				"For", pair.str,
				"expected", pair.exp,
				"got", v,
			)
		}
	}
}
