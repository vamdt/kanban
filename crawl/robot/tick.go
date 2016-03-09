package robot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func raw_cache_filename(id string, t time.Time) string {
	return path.Join(os.Getenv("HOME"), "cache", t.Format("2006/0102"), id)
}

func tick_read_raw_cache(id string, t time.Time) ([]byte, error) {
	return ioutil.ReadFile(raw_cache_filename(id, t))
}

func tick_write_raw_cache(c []byte, id string, t time.Time) {
	if len(c) < 1 {
		return
	}
	f := raw_cache_filename(id, t)
	os.MkdirAll(path.Dir(f), 0755)
	ioutil.WriteFile(f, c, 0644)
}

func Tick_download_from_sina(id string, t time.Time) []byte {
	body, err := tick_read_raw_cache(id, t)
	if err == nil {
		return body
	}

	body, err = Http_get_gbk(Tick_sina_url(id, t), nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	tick_write_raw_cache(body, id, t)
	return body
}

func Tick_download_today_from_sina(id string) []byte {
	url := fmt.Sprintf("http://vip.stock.finance.sina.com.cn/quotes_service/view/CN_TransListV2.php?num=9000&symbol=%s&rn=%d",
		id, time.Now().UnixNano()/int64(time.Millisecond))
	body, err := Http_get_gbk(url, nil)
	if err != nil {
		log.Println(err)
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
	body, err := Http_get_gbk(url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	return body
}

func Tick_sina_url(id string, t time.Time) string {
	return fmt.Sprintf("http://market.finance.sina.com.cn/downxls.php?date=%s&symbol=%s",
		t.Format("2006-01-02"), id)
}
