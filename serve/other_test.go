package main

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
