package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"./crawl"

	"gopkg.in/mgo.v2"
)

type Opt struct {
	today       bool
	stock       string
	mongo       string
	update_days int
}

var opt Opt

func init() {
	flag.BoolVar(&opt.today, "today", true, "update today's tick data")
	flag.StringVar(&opt.stock, "stock", "", "stock id")
	flag.StringVar(&opt.mongo, "mongo", "localhost", "mongo uri")
	flag.IntVar(&opt.update_days, "update_days", 5, "update days")
}

func yoWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		log.Println(messageType)
		log.Println(p)
		if err = conn.WriteMessage(messageType, p); err != nil {
			return
		}
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/yo", yoWs)
	http.HandleFunc("/stock", serveWs)
	http.Handle("/", http.FileServer(http.Dir("static")))

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
