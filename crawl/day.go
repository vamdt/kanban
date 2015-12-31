package crawl

import (
	"fmt"
	"time"
)

func Day_collection_name(id string) string {
	return fmt.Sprintf("%s.tdata.kday", id)
}

func Day_sina_url(id string, t time.Time) string {
	return fmt.Sprintf("http://biz.finance.sina.com.cn/stock/flash_hq/kline_data.php?&rand=9000&symbol=%s&end_date=&begin_date=%s&type=plain",
		id, t.Format("2006-01-02"))
}

func DownloadDaysFromSina(id string, t time.Time) []byte {
	url := Day_sina_url(id, t)
	return Download(url)
}
