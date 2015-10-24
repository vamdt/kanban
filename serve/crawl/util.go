package crawl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
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

func Text2Tdatas(text []byte) Tdatas {
	tds := Tdatas{}
	base := 5
	td := []Tdata{}
	typing := []Typing{}
	tline := []Typing{}
	lines := bytes.Split(text, []byte("\n"))
	lines = lines[1 : len(lines)-1]
	for _, l := range lines {
		if bytes.IndexAny(l, "|-_") > -1 {
			for i, c := 0, len(td); i < c; i++ {
				if td[i].High > 0 {
					td[i].High = td[i].High + base*2
					td[i].Low = td[i].Low + base*2
				}
			}
		}
		for i, c := range l {
			if len(td) <= i {
				td = append(td, Tdata{})
			}

			switch c {
			case 'L':
				tline = append(tline, Typing{I: i, Type: TopTyping})
			case 'l':
				tline = append(tline, Typing{I: i, Type: BottomTyping})
			case '^':
				typing = append(typing, Typing{I: i, Type: TopTyping})
			case '.':
				fallthrough
			case 'v':
				typing = append(typing, Typing{I: i, Type: BottomTyping})
			case '|':
				if td[i].High == 0 {
					td[i].High = base * 3
				}
				td[i].Low = base
			case '-':
				if td[i].High == 0 {
					td[i].High = base * 2
				}
				td[i].Low = base * 2
			case '_':
				if td[i].High == 0 {
					td[i].High = base
				}
				td[i].Low = base
			}
		}
	}
	tds.Data = td
	for i, c := 0, len(typing); i < c; i++ {
		typing[i].High = td[typing[i].I].High
		typing[i].Low = td[typing[i].I].Low
		if typing[i].Type == TopTyping {
			typing[i].Price = td[typing[i].I].High
		} else if typing[i].Type == BottomTyping {
			typing[i].Price = td[typing[i].I].Low
		}
	}
	sort.Sort(TypingSlice(typing))
	tds.Typing.Data = typing

	if llen := len(tline); llen > 0 {
		sort.Sort(TypingSlice(tline))
		for i := llen - 1; i > -1; i-- {
			tline[i].High = td[tline[i].I].High
			tline[i].Low = td[tline[i].I].Low
			if tline[i].Type == TopTyping {
				tline[i].Price = td[tline[i].I].High
				tline[i].Type = DownTyping
			} else if tline[i].Type == BottomTyping {
				tline[i].Price = td[tline[i].I].Low
				tline[i].Type = UpTyping
			}
		}
		if tline[llen-1].Type == DownTyping {
			tline[llen-1].Type = TopTyping
		} else if tline[llen-1].Type == UpTyping {
			tline[llen-1].Type = BottomTyping
		}
		tds.Typing.Line = tline
	}
	return tds
}
