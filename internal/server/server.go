package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	clients  map[*client]bool
	inbound  chan *client
	outbound chan *client
	messages chan []byte
}

const (
	readWriteBufferSize = 1024
	messageBufferSize   = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  readWriteBufferSize,
	WriteBufferSize: readWriteBufferSize,
	CheckOrigin:     func(req *http.Request) bool { return true }, // Default origin check -> no checks
}

func NewServer() *Server {
	return &Server{
		clients:  make(map[*client]bool),
		inbound:  make(chan *client),
		outbound: make(chan *client),
		messages: make(chan []byte),
	}
}

func (server *Server) CreateClient(connection *websocket.Conn) *client {
	return &client{
		connection: connection,
		messagesTo: make(chan []byte, messageBufferSize),
		server:     server,
	}
}

func (server *Server) Run() {
	for {
		select {
		// Client joined
		case client := <-server.inbound:
			server.clients[client] = true
		// Client left
		case client := <-server.outbound:
			delete(server.clients, client)
			close(client.messagesTo)
		// Message received
		case message := <-server.messages:
			for client := range server.clients {
				client.messagesTo <- message
			}
		}
	}
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	connection, err := upgrader.Upgrade(writer, req, nil)

	if err != nil {
		fmt.Println("ERROR: Unable to upgrade HTTP request to the WebSocket protocol!")
		fmt.Println(err.Error())
	}

	client := server.CreateClient(connection)
	fmt.Println("Client connected from '" + client.connection.RemoteAddr().String() + "'!")
	server.inbound <- client
	client.init()
}
