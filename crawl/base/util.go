package base

import (
	"strconv"
	"strings"
)

func ParseCent(s string) int {
	ms := strings.SplitN(s, ".", 3)
	if len(ms) < 1 {
		return 0
	}

	m, _ := strconv.Atoi(ms[0])

	var cent string
	if len(ms) > 1 {
		cent = ms[1]
	}
	cent = cent + "00"
	cent = cent[:2]
	c, _ := strconv.Atoi(cent)
	if strings.HasPrefix(s, "-") {
		return 100*m - c
	}
	return 100*m + c
}
