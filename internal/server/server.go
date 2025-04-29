package server

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	clients  map[*client]bool
	inbound  chan *client
	outbound chan *client
	requests chan *request
}

const (
	readWriteBufferSize = 1024
	messageBufferSize   = 256
)

type requestHandler interface {
	handle(bytes []byte)
}

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
		requests: make(chan *request),
	}
}

func (server *Server) createClient(connection *websocket.Conn) *client {
	return &client{
		connection: connection,
		user:       nil,
		responses:  make(chan []byte, messageBufferSize),
		server:     server,
	}
}

func (server *Server) Run() {
	chatManager := newChatManager(10)
	nextUserID := firstUserID

	// Send a given message to all clients
	broadcastToAll := func(buffer *bytes.Buffer) {
		bytes := buffer.Bytes()

		for client := range server.clients {
			client.send(bytes)
		}
	}

	requestHandlers := make(map[header]func(in *bytes.Buffer, req *request))

	// Handle first connection and exchanging of user info
	requestHandlers[iHeadClientInfo] = func(in *bytes.Buffer, req *request) {
		name := readString(in)
		req.client.user.name = name

		// Send a table of active users paired with their IDs
		res := createResponse(oHeadActiveUsers)
		for active := range chatManager.activeUsers {
			writeUserId(active.id, res)
			writeString(*name, res)
		}

		req.client.send(res.Bytes())
	}

	// Handle user name change
	requestHandlers[iHeadNameChange] = func(in *bytes.Buffer, req *request) {
		name := readString(in)
		req.client.user.name = name

		// Only broadcast name change if the user is active
		if chatManager.contains(req.client.user) {
			res := createResponse(oHeadActiveUsers)
			writeUserId(req.client.user.id, res)
			writeString(*name, res)

			broadcastToAll(res)
		}
	}

	// Handle chat input (message/command)
	requestHandlers[iHeadChatInput] = func(in *bytes.Buffer, req *request) {
	}

	for {
		select {
		// Client joined
		case client := <-server.inbound:
			server.clients[client] = true
			client.user = newUserInfo(nextUserID, nil)
			nextUserID++
		// Client left
		case client := <-server.outbound:
			delete(server.clients, client)
			close(client.responses)
		// Message received
		case req := <-server.requests:
			requestHandlers[req.head](bytes.NewBuffer(req.bytes), req)
		}
	}
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	connection, err := upgrader.Upgrade(writer, req, nil)

	if err != nil {
		fmt.Println("ERROR: Unable to upgrade HTTP request to the WebSocket protocol!")
		fmt.Println(err.Error())
	}

	client := server.createClient(connection)
	fmt.Println("Client connected from '" + client.connection.RemoteAddr().String() + "'!")
	server.inbound <- client
	client.init()
}
