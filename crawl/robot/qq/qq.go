package qq

import (
	"bytes"
	"fmt"
	"time"

	. "../"
	. "../../base"
)

type QQRobot struct {
	RobotBase
}

func init() {
	for i := DefaultRobotConcurrent; i > 0; i-- {
		robot := &QQRobot{}
		Registry(robot)
	}
}

func (p *QQRobot) Can(id string, task int32) bool {
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

func (p *QQRobot) Day_url(id string, t time.Time) string {
	return fmt.Sprintf("http://data.gtimg.cn/flashdata/hushen/daily/%s/%s.js",
		t.Format("06"), id)
}

func (p *QQRobot) Day_latest_url(id string) string {
	return fmt.Sprintf("http://data.gtimg.cn/flashdata/hushen/latest/daily/%s.js",
		id)
}

func (p *QQRobot) tdata_from_line(td *Tdata, line []byte) bool {
	line = bytes.TrimSpace(line)
	infos := bytes.Split(line, []byte(" "))
	if len(infos) != 6 {
		return false
	}

	//timestr, open, high, cloze, low, volume
	td.FromBytes(infos[0], infos[1], infos[3], infos[2], infos[4], infos[5])
	return true
}

func (p *QQRobot) Days_download(id string, start time.Time) (res []Tdata, err error) {
	url := p.Day_latest_url(id)
	body := Download(url)
	if !bytes.HasPrefix(body, []byte(`latest_daily_data`)) {
		return
	}
	lines := bytes.Split(body, []byte("\\n\\"))
	if len(lines) < 3 {
		return
	}

	// start:901219
	start_str := string(ParseParamByte(lines[1], []byte("start"), []byte(" "), []byte(":")))
	start_date, _ := time.Parse(QQmt, start_str)
	if start.Before(start_date) {
		start = start_date
	}

	for i, count := 2, len(lines)-1; i < count; i++ {
		td := Tdata{}
		if !p.tdata_from_line(&td, lines[i]) {
			continue
		}
		res = append(res, td)
	}

	i, ok := (TdataSlice(res)).Search(start.Truncate(time.Hour * 24))
	if !ok {
		return p.years_download(id, start)
	}

	if i >= len(res) {
		res = []Tdata{}
	} else {
		res = res[i:]
	}
	return
}

func (p *QQRobot) years_download(id string, start time.Time) (res []Tdata, err error) {
	for t, ys, ye := start, start.Year(), time.Now().Year()+1; ys < ye; ys++ {
		url := p.Day_url(id, t)
		t = t.AddDate(1, 0, 0)
		body := Download(url)
		if !bytes.HasPrefix(body, []byte(`daily_data_`)) {
			continue
		}
		lines := bytes.Split(body, []byte("\\n\\"))

		for i, count := 1, len(lines)-1; i < count; i++ {
			td := Tdata{}
			if !p.tdata_from_line(&td, lines[i]) {
				continue
			}

			res = append(res, td)
		}
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
	return
}
