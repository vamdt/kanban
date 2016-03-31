package sse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	. "../"
	. "../../base"
	"github.com/golang/glog"
)

const tout time.Duration = time.Second * 30

type SSE struct {
	RobotBase
}

func init() {
	for i := DefaultRobotConcurrent; i > 0; i-- {
		robot := &SSE{}
		Registry(robot)
	}
}

func (p *SSE) Can(id string, task int32) bool {
	if !strings.HasPrefix(id, "sh") {
		return false
	}
	switch task {
	case TaskDay:
		return true
	case TaskMin1:
		return false
	case TaskMin5:
		return false
	case TaskTick:
		return false
	case TaskRealTick:
		return true
	default:
		return false
	}
	return false
}

const (
	cb string = "jQuery11120940553460502997_"
)

func (p *SSE) day_url(id string, t time.Time) string {
	now := time.Now()
	days := int(now.Sub(t) / time.Hour / 24)
	r := now.UnixNano() / int64(time.Millisecond)
	cols := "date%2Copen%2Chigh%2Clow%2Cclose%2Cvolume"
	return fmt.Sprintf("http://yunhq.sse.com.cn:32041/v1/sh1/dayk/%s?callback=%s%d&select=%s&begin=-%d&end=-1&_=%d",
		id[2:], cb, r, cols, days, r)
}

// [20160316,43.52,44.73,42.10,42.82,60043446]
type DayT []json.Number

type DaysRes struct {
	Code  string `json:"code,omitempty"`
	Total int    `json:"total,omitempty"`
	Begin int    `json:"begin,omitempty"`
	End   int    `json:"end,omitempty"`
	Kline []DayT `json:"kline,omitempty"`
}

func (p *SSE) Days_download(id string, start time.Time) (res []Tdata, err error) {
	url := p.day_url(id, start)
	body := Download(url, tout)
	if !bytes.HasPrefix(body, []byte(cb)) {
		glog.Warningln("sse:", url, "prefix not correct", string(body))
		return
	}
	body = ParseParamBeginEnd(body, []byte(`(`), []byte(`)`))
	if body == nil {
		return
	}

	// {"code":"600570","total":2872,"begin":2870,"end":2872,"kline":[[20160316,43.52,44.73,42.10,42.82,60043446],[20160317,42.20,47.10,42.05,47.10,51304803]]}
	sse_res := DaysRes{}
	err = json.Unmarshal(body, &sse_res)
	if err != nil {
		glog.Warningln("sse:", url, err)
		return
	}

	if len(sse_res.Kline) < 1 {
		glog.Warningln("sse:", url, "no res", string(body), sse_res)
		return
	}

	for _, d := range sse_res.Kline {
		// date,open,high,low,close,volume
		// [20160316,43.52,44.73,42.10,42.82,60043446]
		if len(d) < 6 {
			glog.Warningln("sse: days format error, len < 6", d, url)
			continue
		}
		timestr := d[0].String()
		open := d[1].String()
		high := d[2].String()
		low := d[3].String()
		close := d[4].String()
		volume := d[5].String()
		if l := len(volume); l > 2 {
			volume = volume[:l-2]
		}
		td := Tdata{}
		td.FromString(timestr, open, high, close, low, volume)
		res = append(res, td)
	}

	i, _ := (TdataSlice(res)).Search(start.Truncate(time.Hour * 24))
	if i >= len(res) {
		res = []Tdata{}
	} else {
		res = res[i:]
	}
	glog.Infoln("sse get item", len(res))
	return
}

type SnapRes struct {
	Code string            `json:"code,omitempty"`
	Date int               `json:"date,omitempty"`
	Time int               `json:"time,omitempty"`
	Snap []json.RawMessage `json:"snap,omitempty"`
}

func (p *SSE) GetRealtimeTick(id string) (res []RealtimeTickRes) {
	if !strings.HasPrefix(id, "sh") {
		return
	}
	now := time.Now()
	r := now.UnixNano() / int64(time.Millisecond)
	cols := "name%2Clast%2Cchg_rate%2Cchange%2Camount%2Cvolume%2Copen%2Cprev_close%2Cask%2Cbid%2Chigh%2Clow%2Ctradephase"
	url := fmt.Sprintf("http://yunhq.sse.com.cn:32041/v1/sh1/snap/%s?callback=%s%d&select=%s&_=%d",
		id[2:], cb, r, cols, r)
	body, _ := Http_get_gbk(url, nil, tout)
	if !bytes.HasPrefix(body, []byte(cb)) {
		glog.Warningln("sse:", url, "prefix not correct", string(body))
		return
	}
	body = ParseParamBeginEnd(body, []byte(`(`), []byte(`)`))
	if body == nil {
		return
	}

	sres := SnapRes{}
	err := json.Unmarshal(body, &sres)
	if err != nil {
		glog.Warningln("sse:", url, err)
		return
	}

	if len(sres.Snap) < 1 {
		glog.Warningln("sse:", url, "no res", string(body), sres)
		return
	}
	timestr := fmt.Sprintf("%d%06d", sres.Date, sres.Time)
	rt := RealtimeTickRes{Id: sres.Code}
	// ["",57.80,-0.69,-0.40,473199307,8163864,58.00,58.20,[57.70,5800,57.80,9600,57.81,3700,57.83,1600,57.84,200],[57.65,1000,57.64,25331,57.61,1700,57.60,12600,57.59,9200],58.59,57.18,"T111"]}
	rt.Time, _ = time.Parse("20060102150405", timestr)
	json.Unmarshal(sres.Snap[0], &rt.Name)
	var n json.Number
	json.Unmarshal(sres.Snap[1], &n)
	rt.Price = ParseCent(string(n))
	json.Unmarshal(sres.Snap[4], &rt.Turnover)
	json.Unmarshal(sres.Snap[5], &rt.Volume)

	json.Unmarshal(sres.Snap[6], &n)
	open := ParseCent(string(n))
	json.Unmarshal(sres.Snap[7], &n)
	yes_close := ParseCent(string(n))
	rt.Change = open - yes_close
	if rt.Change > 0 {
		rt.Type = Buy_tick
	} else if rt.Change < 0 {
		rt.Type = Sell_tick
	} else {
		rt.Type = Eq_tick
	}
	rt.Change = open
	var buy []json.Number
	json.Unmarshal(sres.Snap[8], &buy)
	rt.Buyone = ParseCent(string(buy[0]))
	var sell []json.Number
	json.Unmarshal(sres.Snap[9], &sell)
	rt.Sellone = ParseCent(string(sell[0]))

	json.Unmarshal(sres.Snap[10], &n)
	rt.High = ParseCent(string(n))
	json.Unmarshal(sres.Snap[11], &n)
	rt.Low = ParseCent(string(n))
	res = append(res, rt)
	return
}
