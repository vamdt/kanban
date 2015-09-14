package main

import (
	"log"

	"github.com/googollee/go-socket.io"
)

var sio *socketio.Server

func init() {
	sio, _ = socketio.NewServer(nil)
	sio.On("connection", func(so socketio.Socket) {
		log.Println("new connect")
		so.Join("chat")
		so.On("chat message", func(msg string) {
			log.Println("chat message", msg)
		})
		sio.BroadcastTo("chat", "new", "f")
	})
	sio.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
}
