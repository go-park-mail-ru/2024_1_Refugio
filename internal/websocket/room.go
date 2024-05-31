package websocket

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	emailApi "mail/internal/models/delivery_models"
	"mail/internal/models/response"
	"mail/internal/pkg/middleware"
	"net/http"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

type room struct {
	// clients holds all current clients in this room.
	clients map[string]*client

	// join is a channel for clients wishing to join the room.
	join chan *client

	// leave is a channel for clients wishing to leave the room.
	leave chan *client

	// forward is a channel that holds incoming messages that should be forwarded to the other clients.
	forward chan []byte
}

// newRoom create a new chat room

func NewRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[string]*client),
	}
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client.login] = client
		case client := <-r.leave:
			delete(r.clients, client.login)
			//close(client.receive)
		case msg := <-r.forward:
			var newEmail emailApi.Email
			if err := newEmail.UnmarshalJSON(msg); err != nil {
				fmt.Println("Bad JSON in request in Run")
			}
			for login, client := range r.clients {
				if login == newEmail.RecipientEmail {
					client.receive <- msg
				}
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	login, ok := vars["login"]
	if !ok {
		response.HandleError(w, http.StatusBadRequest, "Bad login in request")
		return
	}
	newW, _ := w.(*middleware.LoggingResponseWriter)
	NewW := newW.ResponseWriter
	wrap, _ := NewW.(*middleware.LoggingResponseWriter)
	wr := wrap.ResponseWriter
	socket, err := upgrader.Upgrade(wr, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	_, ok = r.clients[login]
	if ok {
		r.leave <- r.clients[login]
	}
	/*
		for cl := range r.clients {
			if cl.login == login {
				cl = &client{
					socket:  socket,
					receive: make(chan []byte, messageBufferSize),
					room:    r,
					login:   login,
				}

			}
		}
	*/
	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
		login:   login,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
