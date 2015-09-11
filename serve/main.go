package main

import (
	"flag"

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

func main() {
	flag.Parse()
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
