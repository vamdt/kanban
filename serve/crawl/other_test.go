package crawl

import (
	"strconv"
	"testing"
)

func TestAtoi(t *testing.T) {
	i, _ := strconv.Atoi("-05")
	if i != -5 {
		t.Error(
			"For", "-05",
			"expected", -5,
			"got", i,
		)
	}
}

func TestByteString(t *testing.T) {
	var b []byte
	s := string(b)
	if s != "" {
		t.Error(
			"For", "",
			"expected", "",
			"got", s,
		)
	}
}
