package main

import (
	"fmt"
	"html"
	"net/http"

	"github.com/golang/glog"

	"./crawl"
)

var stocks *crawl.Stocks = crawl.NewStocks()

func search_handle(w http.ResponseWriter, r *http.Request) {
	sid := r.FormValue("s")
	if len(sid) < 1 {
		http.NotFound(w, r)
		return
	}

	name, err := crawl.Tick_get_name(sid)
	if len(name) < 1 {
		glog.V(100).Infoln(err)
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "{\"sid\": \"%s\", \"name\": \"%s\"}",
		html.EscapeString(sid), html.EscapeString(name))
}
