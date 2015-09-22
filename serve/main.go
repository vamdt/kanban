package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"gopkg.in/mgo.v2"
)

type Opt struct {
	debug bool
	mongo string
}

var opt Opt
var db *mgo.Database

func init() {
	flag.BoolVar(&opt.debug, "debug", true, "debug")
	flag.StringVar(&opt.mongo, "mongo", "127.0.0.1", "mongo uri")
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
	if opt.debug {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	}

	session, err := mgo.Dial(opt.mongo)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db = session.DB("stock")

	go h.run()
	http.HandleFunc("/socket.io/", serveWs)

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
}
