package jqka

import (
	"bytes"
	"fmt"
	"time"

	. "../"
	. "../../base"
	"github.com/golang/glog"
)

type JQKARobot struct {
}

func init() {
	for i := DefaultRobotConcurrent; i > 0; i-- {
		robot := &JQKARobot{}
		Registry(robot)
	}
}

func (p *JQKARobot) Day_url(id string, t time.Time) string {
	return fmt.Sprintf("http://d.10jqka.com.cn/v2/line/hs_%s/01/%s.js",
		id[2:], t.Format("2006"))
}

func (p *JQKARobot) Day_latest_url(id string) string {
	return fmt.Sprintf("http://d.10jqka.com.cn/v2/line/hs_%s/01/last.js",
		id[2:])
}

func (p *JQKARobot) tdata_from_line(td *Tdata, line []byte) bool {
	infos := bytes.Split(line, []byte(","))
	if len(infos) != 8 {
		return false
	}

	//timestr, open,   high,   low,    close,  volume
	//20160217,2829.76,2868.70,2824.36,2867.34,21690992000,225964250000.00,

	//timestr, open, high, cloze, low, volume
	timestr := infos[0]
	open := infos[1]
	high := infos[2]
	close := infos[4]
	low := infos[3]
	volume := infos[5]
	if l := len(volume); l > 2 {
		volume = volume[:l-2]
	}
	td.FromBytes(timestr, open, high, close, low, volume)
	return true
}

func (p *JQKARobot) parse_tdatas(res []Tdata, body []byte) []Tdata {
	data := ParseParamBeginEnd(body, []byte(`"data":"`), []byte(`"`))
	if data == nil {
		return res
	}
	// 20160104,18.28,18.28,17.55,17.80,42240610,754425780.00,0.226;
	lines := bytes.Split(data, []byte(";"))
	for i, count := 0, len(lines); i < count; i++ {
		td := Tdata{}
		if !p.tdata_from_line(&td, lines[i]) {
			continue
		}
		res = append(res, td)
	}
	return res
}

func (p *JQKARobot) Days_download(id string, start time.Time) (res []Tdata, err error) {
	if id == "sh000001" {
		id = "sh1A0001"
	}
	url := p.Day_latest_url(id)
	body := Download(url)
	if !bytes.HasPrefix(body, []byte(`quotebridge_v2_line_hs_`)) {
		return
	}
	body = ParseParamBeginEnd(body, []byte(`(`), nil)
	if body == nil {
		return
	}
	body = bytes.TrimRight(body, ")")

	// "start":"19901219"
	start_str := string(ParseParamBeginEnd(body, []byte(`"start":"`), []byte(`"`)))
	start_date, _ := time.Parse(JQKAmt, start_str)
	if start.Before(start_date) {
		start = start_date
	}

	res = p.parse_tdatas(res, body)

	i, ok := (TdataSlice(res)).Search(start.Truncate(time.Hour * 24))
	if !ok {
		return p.years_download(id, start)
	}

	if i >= len(res) {
		res = []Tdata{}
	} else {
		res = res[i:]
	}
	glog.Infoln("10jqka get item", len(res))
	return
}

func (p *JQKARobot) years_download(id string, start time.Time) (res []Tdata, err error) {
	for t, ys, ye := start, start.Year(), time.Now().Year()+1; ys < ye; ys++ {
		url := p.Day_url(id, t)
		t = t.AddDate(1, 0, 0)
		body := Download(url)
		if !bytes.HasPrefix(body, []byte(`quotebridge_v2_line_hs_`)) {
			continue
		}
		body = ParseParamBeginEnd(body, []byte(`(`), nil)
		if body == nil {
			continue
		}

		res = p.parse_tdatas(res, body)
	}

	if len(res) < 1 {
		return
	}

	i, _ := (TdataSlice(res)).Search(start.Truncate(time.Hour * 24))
	if i >= len(res) {
		res = []Tdata{}
	} else {
		res = res[i:]
	}
	glog.Infoln("10jqka get item", len(res))
	return
}
