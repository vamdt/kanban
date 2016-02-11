package crawl

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

type SinaRobot struct {
}

func init() {
	for i := 6; i > -1; i-- {
		robot := &SinaRobot{}
		Registry(robot)
	}
}

func (p *SinaRobot) Day_url(id string, t time.Time) string {
	return fmt.Sprintf("http://biz.finance.sina.com.cn/stock/flash_hq/kline_data.php?&rand=9000&symbol=%s&end_date=&begin_date=%s&type=plain",
		id, t.Format("2006-01-02"))
}

func (p *SinaRobot) Days_download(id string, start time.Time) (res []Tdata, err error) {
	url := p.Day_url(id, start)
	body := Download(url)
	body = bytes.TrimSpace(body)
	lines := bytes.Split(body, []byte("\n"))

	for i, count := 0, len(lines); i < count; i++ {
		td := Tdata{}
		line := bytes.TrimSpace(lines[i])
		infos := bytes.Split(line, []byte(","))
		if len(infos) != 6 {
			err = errors.New("could not parse line " + string(line))
			return
		}

		//timestr, open, high, cloze, low, volume
		td.FromBytes(infos[0], infos[1], infos[2], infos[3], infos[4], infos[5])
		res = append(res, td)
	}
	return
}
