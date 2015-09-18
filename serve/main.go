package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"./crawl"

	"gopkg.in/mgo.v2"
)

type Opt struct {
	debug       bool
	today       bool
	stock       string
	mongo       string
	update_days int
}

var opt Opt

func init() {
	flag.BoolVar(&opt.debug, "debug", true, "debug")
	flag.BoolVar(&opt.today, "today", true, "update today's tick data")
	flag.StringVar(&opt.stock, "stock", "", "stock id")
	flag.StringVar(&opt.mongo, "mongo", "localhost", "mongo uri")
	flag.IntVar(&opt.update_days, "update_days", 5, "update days")
}

func dev_static_handle(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if strings.HasPrefix(upath, "/bower_components") {
		http.ServeFile(w, r, upath[1:])
		return
	}

	if strings.HasSuffix(upath, "/") {
		upath = upath + "index.html"
	}

	rpath := path.Join("app", upath)
	if _, err := os.Stat(rpath); err == nil {
		http.ServeFile(w, r, rpath)
		return
	}

	rpath = path.Join(".tmp", upath)
	if _, err := os.Stat(rpath); err == nil {
		http.ServeFile(w, r, rpath)
		return
	}
	http.NotFound(w, r)
}

func main() {
	flag.Parse()
	h.init()
  if opt.debug {
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
  }
	http.Handle("/socket.io/", h.io)

	if opt.debug {
		http.HandleFunc("/", dev_static_handle)
	} else {
		http.Handle("/", http.FileServer(http.Dir("static")))
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	addr := ":" + port
	log.Println("serve on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

	if len(opt.stock) < 1 {
		panic("need stock id")
	}

	session, err := mgo.Dial(opt.mongo)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db := session.DB("stock")

	s := crawl.NewQuote(opt.stock)
	if opt.today {
		s.UpdateToday(db)
		return
	}
	s.Update(db, opt.update_days)
}
