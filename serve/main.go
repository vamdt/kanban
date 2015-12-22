package main

import (
	"flag"
	"log"
	"net/http"
	//_ "net/http/pprof"
	"os"

	"./dev"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

type Opt struct {
	debug bool
	https bool
	mongo string
}

var opt Opt
var db *mgo.Database

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.BoolVar(&opt.debug, "debug", false, "debug")
	flag.BoolVar(&opt.https, "https", false, "https")
	flag.StringVar(&opt.mongo, "mongo", "127.0.0.1", "mongo uri")
}

func serve() {
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
	http.HandleFunc("/search", search_handle)

	port := os.Getenv("PORT")

	if opt.debug {
    if len(port) == 0 {
      port = ":3000"
    }
		dev.Start(opt.https, port)
		defer dev.Exit()
		http.Handle("/", dev.Dev)
	} else {
		http.Handle("/", http.FileServer(http.Dir("static")))
	}

	glog.Infoln("serve on", port)
	if opt.https {
		http.ListenAndServeTLS(port, "conf/cert.pem", "conf/key.pem", nil)
	} else {
		http.ListenAndServe(port, nil)
	}
}

func main() {
	flag.Parse()
	if opt.debug {
		flag.Lookup("logtostderr").Value.Set("true")
	}

	defer func() {
		if err := recover(); err != nil {
			glog.Warningln(err)
		}
	}()

	serve()
}
