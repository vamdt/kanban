package sina

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/glog"

	. "../"
	. "../../base"
)

type SinaRobot struct {
	RobotBase
}

func init() {
	for i := DefaultRobotConcurrent; i > 0; i-- {
		robot := &SinaRobot{}
		Registry(robot)
	}
}

func (p *SinaRobot) Can(id string, task int32) bool {
	switch task {
	case TaskDay:
		return true
	case TaskMin1:
		return false
	case TaskMin5:
		return false
	case TaskTick:
		return true
	case TaskRealTicks:
		fallthrough
	case TaskRealTick:
		return true
	default:
		return false
	}
	return false
}

func (p *SinaRobot) Day_url(id string, t time.Time) string {
	return fmt.Sprintf("http://biz.finance.sina.com.cn/stock/flash_hq/kline_data.php?rand=9000&symbol=%s&end_date=&begin_date=%s&type=plain",
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

func (p *SinaRobot) stock_in_cate(item *CategoryItem, code string) {
	for i := 1; ; i++ {
		url := fmt.Sprintf("http://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData?page=%d&num=80&sort=symbol&asc=1&node=%s&symbol=&_s_r_a=page",
			i, code)
		c, err := Http_get_gbk(url, nil)
		if err != nil {
			break
		}

		n := item.LeafCount()
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
			item.AddStock(id)
		}
		if item.LeafCount()-n < 80 {
			break
		}
	}
}

func (p *SinaRobot) sub_cate(c *CategoryItem, cont []byte) {
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
		if c.Sub == nil {
			c.Sub = *NewCategory()
		}
		if _, ok := c.Sub[name]; !ok {
			c.Sub[name] = *NewCategoryItem(name)
		}
		sc := c.Sub[name]
		p.stock_in_cate(&sc, code)
		c.Sub[name] = sc
	}
}

func (p *SinaRobot) Cate(tc Category) {
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
			tc[key] = *NewCategoryItem(key)
		}
		item := tc[key]
		p.sub_cate(&item, cont)
		tc[key] = item
	}

}

func parse_realtime_tick(r *RealtimeTick, line []byte) {
	infos := bytes.Split(line, []byte(","))
	if len(infos) < 33 {
		glog.Warningln("sina hq api, res format changed")
		return
	}

	r.Name = string(infos[0])
	nul := []byte("")
	t, _ := time.Parse("2006-01-02", string(infos[30]))
	timestr := infos[31]
	price := infos[3]
	change := nul
	volume := infos[8]
	turnover := infos[9]
	typestr := nul
	r.FromString(t, timestr, price, change, volume, turnover, typestr)
	r.High = ParseCent(string(infos[4]))
	r.Low = ParseCent(string(infos[5]))
	r.Buyone = ParseCent(string(infos[11]))
	r.Sellone = ParseCent(string(infos[21]))
	open := ParseCent(string(infos[1]))

	//"00":"","01":"临停1H","02":"停牌","03":"停牌","04":"临停","05":"停1/2","07":"暂停","-1":"无记录","-2":"未上市","-3":"退市"
	r.Status, _ = strconv.Atoi(string(infos[32]))
	if r.Status == 3 {
		r.Status = 2
	}

	if r.Price > open {
		r.Type = Buy_tick
	} else if r.Price < open {
		r.Type = Sell_tick
	} else {
		r.Type = Eq_tick
	}
	r.Change = open
}

func (p *SinaRobot) GetRealtimeTick(ids string) (res []RealtimeTickRes) {
	if len(ids) < 1 {
		return
	}
	url := fmt.Sprintf("http://hq.sinajs.cn/rn=%d&list=%s",
		time.Now().UnixNano()/int64(time.Millisecond), ids)
	body, err := Http_get_gbk(url, nil)
	if err != nil {
		glog.Warningln(err)
		return
	}

	for _, line := range bytes.Split(body, []byte("\";")) {
		line = bytes.TrimSpace(line)
		info := bytes.Split(line, []byte("=\""))
		if len(info) != 2 {
			continue
		}
		prefix := "var hq_str_"
		if !bytes.HasPrefix(info[0], []byte(prefix)) {
			continue
		}
		id := string(info[0][len(prefix):])
		rt := RealtimeTickRes{Id: id}
		parse_realtime_tick(&rt.RealtimeTick, info[1])
		res = append(res, rt)
	}
	return
}
