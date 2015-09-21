package main

import "time"

const (
	tickPeriod = 2 * time.Second
)

type hub struct {
	connections map[string][]*connection

	broadcast chan []byte

	register   chan *watchRequest
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *watchRequest),
	unregister:  make(chan *connection),
	connections: make(map[string][]*connection),
}

func (h *hub) do_register(r *watchRequest) {
	name := r.StockId + ":" + r.Fq + ":" + r.Level
	if _, ok := h.connections[name]; !ok {
		h.connections[name] = []*connection{}
	}
	conns := h.connections[name]
	conns = append(conns, r.Conn)
	h.connections[name] = conns
}

func (h *hub) do_unregister(c *connection) {
	holder := make(map[string][]*connection)
	for name, conns := range h.connections {
		if conns == nil {
			continue
		}
		connections := []*connection{}
		for _, conn := range conns {
			if conn != c {
				connections = append(connections, conn)
			}
		}
		holder[name] = connections
		h.connections[name] = conns[:len(conns)-1]
	}
	for name, conns := range holder {
		h.connections[name] = conns
	}
	close(c.send)
}

func (h *hub) run() {
	ticker := time.NewTicker(tickPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case req := <-h.register:
			h.do_register(req)
		case c := <-h.unregister:
			h.do_unregister(c)
		case <-ticker.C:
		case m := <-h.broadcast:
			for _, conns := range h.connections {
				if conns == nil {
					continue
				}
				for _, c := range conns {
					select {
					case c.send <- m:
					default:
						h.do_unregister(c)
					}
				}
			}
		}
	}
}
