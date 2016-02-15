package crawl

import (
	"bytes"
	"fmt"
	"time"
)

type QQRobot struct {
}

func init() {
	for i := 4; i > 0; i-- {
		robot := &QQRobot{}
		Registry(robot)
	}
}

func (p *QQRobot) Day_url(id string, t time.Time) string {
	return fmt.Sprintf("http://data.gtimg.cn/flashdata/hushen/daily/%s/%s.js",
		t.Format("06"), id)
}

func (p *QQRobot) Day_latest_url(id string) string {
	return fmt.Sprintf("http://data.gtimg.cn/flashdata/hushen/latest/daily/%s.js",
		id)
}

func (p *QQRobot) Days_latest_download(id string, start time.Time) (res []Tdata, err error) {
	url := p.Day_latest_url(id)
	body := Download(url)
	if !bytes.HasPrefix(body, []byte(`latest_daily_data`)) {
		return
	}
	lines := bytes.Split(body, []byte("\\n\\"))
	if len(lines) < 3 {
		return
	}

	for i, count := 2, len(lines)-1; i < count; i++ {
		line := bytes.TrimSpace(lines[i])
		infos := bytes.Split(line, []byte(" "))
		if len(infos) != 6 {
			continue
		}

		td := Tdata{}
		//timestr, open, high, cloze, low, volume
		td.FromBytes(infos[0], infos[1], infos[3], infos[2], infos[4], infos[5])
		res = append(res, td)
	}
	return
}

func (p *QQRobot) Days_download(id string, start time.Time) (res []Tdata, err error) {
	for t, ys, ye := start, start.Year(), time.Now().Year()+1; ys < ye; ys++ {
		url := p.Day_url(id, t)
		t = t.AddDate(1, 0, 0)
		body := Download(url)
		if !bytes.HasPrefix(body, []byte(`daily_data_`)) {
			continue
		}
		lines := bytes.Split(body, []byte("\\n\\"))

		for i, count := 1, len(lines)-1; i < count; i++ {
			line := bytes.TrimSpace(lines[i])
			infos := bytes.Split(line, []byte(" "))
			if len(infos) != 6 {
				continue
			}

			td := Tdata{}
			//timestr, open, high, cloze, low, volume
			td.FromBytes(infos[0], infos[1], infos[3], infos[2], infos[4], infos[5])
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
