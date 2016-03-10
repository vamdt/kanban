package base

import "testing"

func TestString(t *testing.T) {
	tests := map[int]string{
		0: "",
		1: "sz000001",
	}

	for code, exp := range tests {
		c := Code(code)
		if c.String() != exp {
			t.Error(
				"For", code,
				"expected", exp,
				"got", c,
			)
		}
	}
}

func TestEvilSymbol(t *testing.T) {
	tests := []string{
		"",
		"sh000001",
		"sw000001",
	}

	for i, sym := range tests {
		if _, ok := FromSymbol(sym); ok {
			t.Error(
				"For", "case", i, "FromSymbol",
				"expected", "false",
				"got", "true",
			)
		}
	}
}

func TestFromSymbol(t *testing.T) {
	tests := map[string]int{
		"sh600000": 600000,
		"sz000001": 1,
	}

	for sym, code := range tests {
		exp := NewCode(code)
		if c, ok := FromSymbol(sym); !ok {
			t.Error(
				"For", "case", sym,
				"expected", "FromSymbol true",
				"got", false,
			)
		} else if c == nil || *c != *exp {
			t.Error(
				"For", "case", sym,
				"expected", code,
				"got", c,
			)
		}
	}
}
