package main

import (
	"flag"
	"log"
	"net/http"
	//_ "net/http/pprof"
	"os"

	"./crawl"
	_ "./crawl/robot/jqka"
	_ "./crawl/robot/qq"
	_ "./crawl/robot/sina"
	_ "./crawl/store/mem"
	_ "./crawl/store/mysql"
	"./dev"
	"github.com/golang/glog"
)

type Opt struct {
	debug bool
	serve bool
	play  int
	https bool
	store string

	update_cate bool

	update_factor bool

	min_hub_height int
}

var opt Opt

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.BoolVar(&opt.debug, "debug", false, "debug")
	flag.BoolVar(&opt.serve, "serve", true, "serve mode")
	flag.BoolVar(&opt.update_cate, "update_cate", false, "update cate")
	flag.BoolVar(&opt.update_factor, "update_factor", false, "update factor")
	flag.IntVar(&opt.play, "play", 0, "play mode, ms/tick")
	flag.BoolVar(&opt.https, "https", false, "https")
	flag.StringVar(&opt.store, "store", "mem", "back store with")
	flag.IntVar(&opt.min_hub_height, "min_hub_height", 0, "min hub height")
}

func serve() {
	if opt.debug {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	}

	go h.run()
	http.HandleFunc("/socket.io/", serveWs)
	http.HandleFunc("/search", search_handle)
	http.HandleFunc("/plate", plates_handle)
	http.HandleFunc("/star", star_handle)
	http.HandleFunc("/lucky", lucky_handle)

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
		glog.Infoln("debug on")
	}

	if opt.update_cate {
		glog.Infoln("update_cate on")
		crawl.UpdateCate(opt.store)
	}

	if opt.update_factor {
		glog.Infoln("update_factor on")
		crawl.UpdateFactor(opt.store)
	}

	glog.Infoln("serve mode", opt.serve)
	if !opt.serve {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			glog.Warningln(err)
		}
	}()

	serve()
}
