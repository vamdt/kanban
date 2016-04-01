package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"

	"./crawl"
	"./crawl/base"
	"./crawl/robot"
)

var stocks *crawl.Stocks

func jsonp(w http.ResponseWriter, r *http.Request, data interface{}) {
	cb := r.FormValue("cb")
	if len(cb) > 0 {
		cb = strings.Fields(cb)[0]
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if len(cb) > 0 {
		fmt.Fprintf(w, "/**/ typeof %s === 'function' && %s(", cb, cb)
	}
	buf, _ := json.Marshal(data)
	w.Write(buf)
	if len(cb) > 0 {
		fmt.Fprintf(w, ");")
	}
}

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
	body, err := robot.Http_get(url, nil, time.Second*10)
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

func plates_handle(w http.ResponseWriter, r *http.Request) {
	pid := 0
	if pidstr := r.FormValue("pid"); len(pidstr) > 0 {
		pid, _ = strconv.Atoi(pidstr)
	}
	data, _ := stocks.Store().LoadCategories()
	sel := []base.CategoryItemInfo{}
	for _, d := range data {
		if d.Pid == pid {
			sel = append(sel, d)
		}
	}
	jsonp(w, r, sel)
}

func star_handle(w http.ResponseWriter, r *http.Request) {
	sid := r.FormValue("s")
	if len(sid) < 1 {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		stocks.Store().Star(-1, sid)
		w.Write(nil)
		return
	}

	if r.Method == http.MethodDelete {
		stocks.Store().UnStar(-1, sid)
		w.Write(nil)
		return
	}

	if r.Method == http.MethodGet {
		is := stocks.Store().IsStar(-1, sid)
		jsonp(w, r, map[string]bool{"star": is})
		return
	}
}

func expect_json_res(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		return true
	}
	if strings.Contains(accept, "text/json") {
		return true
	}
	return false
}

func lucky_handle(w http.ResponseWriter, r *http.Request) {
	pool, _ := strconv.Atoi(r.FormValue("pool"))
	sid := r.FormValue("s")
	sid = stocks.Store().Lucky(pool, sid)
	if expect_json_res(r) {
		jsonp(w, r, map[string]string{"lucky": sid})
		return
	}
	http.Redirect(w, r, "/#/"+sid+"/1", http.StatusFound)
}
