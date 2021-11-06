package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type user struct {
	socket *websocket.Conn
	send   chan []byte
	chat   *chat
}

func (u *user) get() {
	defer func(socket *websocket.Conn) {
		err := socket.Close()
		if err != nil {
			panic(err)
		}
	}(u.socket)
	for {
		_, message, err := u.socket.ReadMessage()
		if err != nil {
			log.Fatal("Reading error: ", err)
			return
		}
		u.chat.forward <- message
	}
}

func (u *user) post() {
	defer func(socket *websocket.Conn) {
		err := socket.Close()
		if err != nil {
			panic(err)
		}
	}(u.socket)
	for message := range u.send {
		err := u.socket.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Fatal("Writing error: ", err)
			return
		}
	}
}
