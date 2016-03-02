package crawl

import (
	"bytes"
	"strconv"
	"strings"
	"time"
)

func maxInt(a ...int) int {
	v := a[0]
	for i := len(a) - 1; i > 0; i-- {
		if a[i] > v {
			v = a[i]
		}
	}
	return v
}

func minInt(a ...int) int {
	v := a[0]
	for i := len(a) - 1; i > 0; i-- {
		if a[i] < v {
			v = a[i]
		}
	}
	return v
}

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

func LatestTradeDay() time.Time {
	t := time.Now().UTC()
	h, m, _ := t.Clock()
	if h < 1 || (h == 1 && m < 30) {
		t = t.AddDate(0, 0, -1)
	}
	for !IsTradeDay(t) {
		t = t.AddDate(0, 0, -1)
	}
	return t
}

func IsTradeDay(t time.Time) bool {
	switch t.Weekday() {
	case time.Sunday:
		return false
	case time.Saturday:
		return false
	}
	return true
}

func IsTradeTime(t time.Time) bool {
	t = t.UTC()
	if !IsTradeDay(t) {
		return false
	}
	h, m, _ := t.Clock()
	if h < 1 || h > 7 { // [00:00 - 09:00)  [16:00 - 00:00)
		return false
	} else if h == 1 && m < 25 { // 09:25
		return false
	} else if h == 7 && m > 5 { // 15:05
		return false
	} else if h == 3 && m > 35 { // 11:35
		return false
	} else if h == 4 && m < 59 { // 12:59
		return false
	}
	return true
}

func Monthend(t time.Time) time.Time {
	_, _, d := t.Date()
	t = t.AddDate(0, 1, 1-d)
	return t.Truncate(time.Hour * 24)
}

func Weekend(t time.Time) time.Time {
	wd := t.Weekday()
	if wd < time.Saturday {
		t = t.AddDate(0, 0, int(time.Saturday-wd))
	}
	return t.Truncate(time.Hour * 24)
}

func Minuteend(t time.Time) time.Time {
	return t.Truncate(time.Minute).Add(time.Minute)
}

func Minute5end(t time.Time) time.Time {
	return t.Truncate(5 * time.Minute).Add(5 * time.Minute)
}

func Minute30end(t time.Time) time.Time {
	return t.Truncate(30 * time.Minute).Add(30 * time.Minute)
}

func ParseParamBeginEnd(s, begin, end []byte) []byte {
	i := bytes.Index(s, begin)
	if i < 0 {
		return nil
	}
	s = s[i+len(begin):]

	if end == nil {
		return s
	}

	i = bytes.Index(s, end)
	if i < 0 {
		return nil
	}
	return s[:i]
}

func ParseParamByte(s, name, sep, eq []byte) []byte {
	lines := bytes.Split(s, sep)
	for i, _ := range lines {
		if !bytes.HasPrefix(lines[i], name) {
			continue
		}
		v := bytes.Split(lines[i], eq)
		if len(v) > 2 {
			return v[2]
		}
		break
	}
	return nil
}

func ParseParamInt(s, name, sep, eq []byte, defv int) int {
	b := ParseParamByte(s, name, sep, eq)
	if len(b) > 0 {
		i, _ := strconv.Atoi(string(b))
		return i
	}
	return defv
}
