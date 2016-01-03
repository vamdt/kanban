package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"

	"./crawl"
)

var stocks *crawl.Stocks

func search_handle(w http.ResponseWriter, r *http.Request) {
	sid := r.FormValue("s")
	if len(sid) < 1 {
		glog.V(100).Infoln("s < 1")
		http.NotFound(w, r)
		return
	}

	name := fmt.Sprintf("suggestdata_%d", time.Now().UnixNano()/int64(time.Millisecond))
	url := fmt.Sprintf("http://suggest3.sinajs.cn/suggest/type=11,12,13,14,15&key=%s&name=%s",
		sid, name)
	body, err := crawl.Http_get_gbk(url, nil)
	if err != nil {
		glog.Warningln(err)
		http.NotFound(w, r)
		return
	}
	info := bytes.Split(body, []byte("\""))
	if len(info) < 2 {
		glog.Warningln("fmt err", string(body))
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write(info[1])
}
