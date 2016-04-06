package robot

import (
	"fmt"
	"time"

	"github.com/golang/glog"
)

const tout time.Duration = time.Second * 20

func Tick_download_from_sina(id string, t time.Time) []byte {
	body, err := Http_get(Tick_sina_url(id, t), nil, tout)
	if err != nil {
		glog.Warningln(err)
		return nil
	}

	return body
}

func Tick_download_today_from_sina(id string) []byte {
	url := fmt.Sprintf("http://vip.stock.finance.sina.com.cn/quotes_service/view/CN_TransListV2.php?num=9000&symbol=%s&rn=%d",
		id, time.Now().UnixNano()/int64(time.Millisecond))
	body, err := Http_get(url, nil, tout)
	if err != nil {
		glog.Warningln(err)
		return nil
	}

	return body
}

func Tick_download_real_from_sina(id string) []byte {
	if len(id) < 1 {
		return nil
	}
	url := fmt.Sprintf("http://hq.sinajs.cn/rn=%d&list=%s",
		time.Now().UnixNano()/int64(time.Millisecond), id)
	body, err := Http_get(url, nil, tout)
	if err != nil {
		glog.Warningln(err)
		return nil
	}

	return body
}

func Tick_sina_url(id string, t time.Time) string {
	return fmt.Sprintf("http://market.finance.sina.com.cn/downxls.php?date=%s&symbol=%s",
		t.Format("2006-01-02"), id)
}
