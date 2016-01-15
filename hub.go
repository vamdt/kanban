package main

import (
	"log"
	"sync"

	"github.com/golang/glog"

	"./crawl"
)

type hub struct {
	connections map[string][]*connection

	broadcast chan *crawl.Stock

	register   chan *watchRequest
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan *crawl.Stock),
	register:    make(chan *watchRequest),
	unregister:  make(chan *connection),
	connections: make(map[string][]*connection),
}

func (h *hub) do_register(r *watchRequest) {
	name := r.StockId
	log.Println("register", name)
	if _, ok := h.connections[name]; !ok {
		h.connections[name] = []*connection{}
	}
	has := false
	conns := h.connections[name]
	for _, conn := range conns {
		if conn == r.Conn {
			has = true
			log.Println("the conn had regged")
			break
		}
	}
	if !has {
		conns = append(conns, r.Conn)
		h.connections[name] = conns
	}
	if s, isnew := stocks.Watch(name); !isnew {
		h.send(s, r.Conn)
	}
}

func (h *hub) do_unregister(c *connection) {
	glog.Infoln("in unregister c closed=", c.closed)
	holder := make(map[string][]*connection)
	for name, conns := range h.connections {
		if conns == nil {
			continue
		}
		has := false
		for _, conn := range conns {
			if conn == c {
				has = true
			}
		}
		if !has {
			continue
		}
		connections := []*connection{}
		for _, conn := range conns {
			if conn != c {
				connections = append(connections, conn)
			}
		}
		holder[name] = connections
	}
	for name, conns := range holder {
		h.connections[name] = conns
		stocks.UnWatch(name)
	}
	c.Close()
}

func (h *hub) do_broadcast(m *crawl.Stock) {
	conns, ok := h.connections[m.Id]
	if !ok {
		return
	}

	data, err := m.MarshalTail(true)
	if err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup
	for i, c := 0, len(conns); i < c; i++ {
		wg.Add(1)
		go func(c *connection) {
			defer wg.Done()
			c.Send(data)
		}(conns[i])
	}
	wg.Wait()
}

func (h *hub) send(m *crawl.Stock, c *connection) {
	data, err := m.MarshalTail(false)
	if err != nil {
		log.Println(err)
		return
	}
	c.Send(data)
}

func (h *hub) run() {
	if !opt.debug {
		opt.play = 0
	}
	glog.Infof("new stocks with store [%s] min_hub_height [%d]", opt.store, opt.min_hub_height)
	stocks = crawl.NewStocks(opt.store, opt.play, opt.min_hub_height)
	stocks.Chan(h.broadcast)
	go stocks.Run()

	for {
		select {
		case req := <-h.register:
			h.do_register(req)
		case c := <-h.unregister:
			h.do_unregister(c)
		case m := <-h.broadcast:
			h.do_broadcast(m)
		}
	}
}
