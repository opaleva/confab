package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type chat struct {
	forward chan []byte
	begin   chan *user
	end     chan *user
	users   map[*user]bool
}

func newChat() *chat {
	return &chat{
		forward: make(chan []byte),
		begin:   make(chan *user),
		end:     make(chan *user),
		users:   make(map[*user]bool),
	}
}

func (c *chat) start() {
	for {
		select {
		case user := <-c.begin:
			c.users[user] = true
		case user := <-c.end:
			delete(c.users, user)
			close(user.send)
		case message := <-c.forward:
			for user := range c.users {
				user.send <- message
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

func (c *chat) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	socket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}
	user := &user{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		chat:   c,
	}
	c.begin <- user
	defer func() { c.end <- user }()
	go user.post()
	user.get()
}
