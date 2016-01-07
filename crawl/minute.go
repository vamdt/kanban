package crawl

import "fmt"

func M_sina_url(id, mins string) string {
	return fmt.Sprintf("http://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?symbol=%s&scale=%s&ma=no&datalen=1000",
		id, mins)
}

func M_collection_name(id, mins string) string {
	return fmt.Sprintf("%s.tdata.k%s", id, mins)
}
