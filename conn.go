package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 10240,
}

type connection struct {
	ws *websocket.Conn

	send   chan []byte
	closed int32
	last   int64
}

type watchRequest struct {
	StockId string      `json:"s"`
	Fq      string      `json:"fq"`
	Level   string      `json:"k"`
	Event   string      `json:"event"`
	Num     int         `json:"num,omitempty"`
	Conn    *connection `json:"-"`
}

func (c *connection) Send(data []byte) bool {
	n := time.Now().Unix()
	if atomic.LoadInt64(&c.last)+2 > n {
		return false
	}
	atomic.StoreInt64(&c.last, n)
	select {
	case c.send <- data:
	default:
		c.Close()
		return false
	}
	return true
}

func (c *connection) Close() {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return
	}
	go func() { h.unregister <- c }()
	c.ws.Close()
	close(c.send)
}

func (c *connection) watch(req *watchRequest) {
	log.Println("watch", req)
	req.Conn = c
	h.register <- req
	log.Println("watch send", req)
}

func (c *connection) readPump() {
	defer c.Close()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		wreq := &watchRequest{}
		err := c.ws.ReadJSON(wreq)
		if err != nil {
			log.Println("read request from ws fail", err)
			break
		}
		c.watch(wreq)
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	go c.writePump()
	c.readPump()
}
