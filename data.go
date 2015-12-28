package main

import (
	"fmt"
	"html"
	"net/http"

	"./crawl"
)

var stocks crawl.Stocks

func search_handle(w http.ResponseWriter, r *http.Request) {
	sid := r.FormValue("s")
	if len(sid) < 1 {
		http.NotFound(w, r)
		return
	}

	name, _ := crawl.Tick_get_name(sid)
	if len(name) < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "{\"sid\": \"%s\", \"name\": \"%s\"}",
		html.EscapeString(sid), html.EscapeString(name))
}
