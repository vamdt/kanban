package main

import (
	"encoding/json"
	"log"

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
	if _, ok := h.connections[name]; !ok {
		h.connections[name] = []*connection{}
	}
	conns := h.connections[name]
	for _, conn := range conns {
		if conn == r.Conn {
			return
		}
	}
	conns = append(conns, r.Conn)
	h.connections[name] = conns
	stocks.Watch(r.StockId)
}

func (h *hub) do_unregister(c *connection) {
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
	close(c.send)
}

func (h *hub) run() {
	stocks.DB(db)
	stocks.Chan(h.broadcast)
	go stocks.Run()

	for {
		select {
		case req := <-h.register:
			h.do_register(req)
		case c := <-h.unregister:
			h.do_unregister(c)
		case m := <-h.broadcast:
			if conns, ok := h.connections[m.Id]; ok {
				data, err := json.Marshal(m)
				if err != nil {
					log.Println(err)
					continue
				}
				for _, c := range conns {
					select {
					case c.send <- data:
					default:
						h.do_unregister(c)
					}
				}
			}
		}
	}
}
