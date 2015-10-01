package crawl

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func Time2ObjectId(t time.Time) bson.ObjectId {
	var b [12]byte
	binary.BigEndian.PutUint32(b[:4], uint32(t.Unix()))
	binary.BigEndian.PutUint16(b[4:6], uint16(t.Nanosecond()/int(time.Millisecond)))
	return bson.ObjectId(string(b[:]))
}

func ObjectId2Time(oid bson.ObjectId) time.Time {
	id := string(oid)
	if len(oid) != 12 {
		panic(fmt.Sprintf("Invalid ObjectId: %q", id))
	}
	secs := int64(binary.BigEndian.Uint32([]byte(id[0:4])))
	nsec := int64(binary.BigEndian.Uint16([]byte(id[4:6]))) * int64(time.Millisecond)
	return time.Unix(secs, nsec).UTC()
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

func IsTradeDay(t time.Time) bool {
	switch t.Weekday() {
	case time.Sunday:
		return false
	case time.Saturday:
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
