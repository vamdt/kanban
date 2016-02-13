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

func (p *SinaRobot) Cate(tc TopCategory) {
	p.top_cate(tc)
	return
}

func (p *SinaRobot) stock_in_cate(item *CategoryItem, code string) {
	for i := 1; ; i++ {
		url := fmt.Sprintf("http://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData?page=%d&num=80&sort=symbol&asc=1&node=%s&symbol=&_s_r_a=page",
			i, code)
		c, err := Http_get_gbk(url, nil)
		if err != nil {
			break
		}

		n := len(item.Id)
		for len(c) > 0 {
			end := []byte(`symbol`)
			if i := bytes.Index(c, end); i > -1 {
				c = c[i+len(end):]
			} else {
				break
			}

			end = []byte(`,`)
			id := ""
			if i := bytes.Index(c, end); i > -1 {
				id = string(bytes.Trim(c[:i], `:" `))
				c = c[i+len(end):]
			} else {
				break
			}

			if len(id) < 1 {
				break
			}
			item.Id = append(item.Id, id)
		}
		if len(item.Id)-n < 80 {
			break
		}
	}
}

func (p *SinaRobot) real_cate(c Category, cont []byte) {
	cont = bytes.Trim(cont, `[],"`)
	cont = bytes.Replace(cont, []byte(`","","`), []byte(","), -1)
	cont = bytes.Replace(cont, []byte(`"`), []byte(""), -1)
	lines := bytes.Split(cont, []byte(`],[`))
	for _, l := range lines {
		kv := bytes.Split(l, []byte(`,`))
		if len(kv) < 2 {
			break
		}
		name, code := string(kv[0]), string(kv[1])
		citem := NewCategoryItem(name)
		p.stock_in_cate(citem, code)
		c[name] = *citem
	}
}

func (p *SinaRobot) top_cate(tc TopCategory) {
	url := "http://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodes"
	c, err := Http_get_gbk(url, nil)
	if err != nil {
		return
	}

	end := []byte(`[["沪深股市",[`)
	if i := bytes.Index(c, end); i > -1 {
		c = c[i+len(end):]
	} else {
		return
	}

	end = []byte(`"hangye","cn"]`)
	if i := bytes.Index(c, end); i > -1 {
		c = c[:i+len(end)]
	} else {
		return
	}

	empty := []byte("")
	for len(c) > 0 {
		name, cont := empty, empty

		end = []byte(`[[`)
		if i := bytes.Index(c, end); i > -1 {
			name = c[:i]
			c = c[i+len(end):]
		} else {
			break
		}

		if i := bytes.LastIndexByte(name, '['); i > -1 {
			name = name[i+1:]
		}
		name = bytes.Trim(name, `[],"`)
		if len(name) < 1 {
			break
		}

		key := string(name)
		if key == "证监会行业" {
			break
		}

		end = []byte(`]]`)
		if i := bytes.Index(c, end); i > -1 {
			cont = c[:i]
			c = c[i+len(end):]
		} else {
			break
		}

		if _, ok := tc[key]; !ok {
			tc[key] = *NewCategory()
		}
		p.real_cate(tc[key], cont)
	}

}
