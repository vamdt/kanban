package crawl

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

type M1 struct {
	Tdata
}

type M5 struct {
	Tdata
}

type M30 struct {
	Tdata
}

type M1s struct {
	Tdatas
}

type M5s struct {
	Tdatas
}

type M30s struct {
	Tdatas
}

func M_sina_url(id, mins string) string {
	return fmt.Sprintf("http://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?symbol=%s&scale=%s&ma=no&datalen=1000",
		id, mins)
}

func M5_sina_url(id string) string {
	return M_sina_url(id, "5")
}

func M30_sina_url(id string) string {
	return M_sina_url(id, "30")
}

func DownloadM30sFromSina(id string) []byte {
	url := M30_sina_url(id)
	return Download(url)
}

func DownloadM5sFromSina(id string) []byte {
	url := M5_sina_url(id)
	return Download(url)
}

func M_collection_name(id, mins string) string {
	return fmt.Sprintf("%s.tdata.k%s", id, mins)
}

func M1_collection(db *mgo.Database, id string) *mgo.Collection {
	return db.C(M_collection_name(id, "1"))
}

func M5_collection(db *mgo.Database, id string) *mgo.Collection {
	return db.C(M_collection_name(id, "5"))
}

func M30_collection(db *mgo.Database, id string) *mgo.Collection {
	return db.C(M_collection_name(id, "30"))
}
