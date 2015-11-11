package crawl

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

func M_sina_url(id, mins string) string {
	return fmt.Sprintf("http://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?symbol=%s&scale=%s&ma=no&datalen=1000",
		id, mins)
}

func M_collection_name(id, mins string) string {
	return fmt.Sprintf("%s.tdata.k%s", id, mins)
}

func M1_collection(db *mgo.Database, id string) *mgo.Collection {
	return db.C(M_collection_name(id, "1"))
}
