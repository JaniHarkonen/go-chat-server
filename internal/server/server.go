package server

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

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
	chatManager := newChatManager(10, 10)
	nextUserID := firstUserID

	nextUserID++

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
		fmt.Println("client info received")
		client := req.client
		name := readString(in)
		client.user.name = name
		fmt.Println("client username is:", "'"+(*name)+"'")

		// Write a table of active users paired with their IDs
		res := createResponse(oHeadCompleteUpdate)
		writeUInt32((uint32)(len(chatManager.activeUsers)), res) // Active user count

		for active := range chatManager.activeUsers {
			writeUserInfo(active, res)
		}

		// Write a snapshot of latest chat messages
		visible := chatManager.visibleMessages()

		if len(visible) <= 0 {
			writeEmpty(res)
		} else {
			for _, msg := range visible {
				writeUserId(msg.user.id, res)
				writeString(*msg.message, res)
			}
		}

		fmt.Println("sending response to client info")
		client.send(res.Bytes())
	}

	// Handle user name change
	requestHandlers[iHeadNameChange] = func(in *bytes.Buffer, req *request) {
		client := req.client
		name := readString(in)
		client.user.name = name

		// Only broadcast name change if the user is active
		if chatManager.contains(client.user) {
			res := createResponse(oHeadNameChange)
			writeUserId(client.user.id, res)
			writeString(*name, res)

			broadcastToAll(res)
		}
	}

	// Handle chat input (message/command)
	requestHandlers[iHeadChatInput] = func(in *bytes.Buffer, req *request) {
		client := req.client
		msg := readString(in)

		activated, deactivated := chatManager.post(client.user, msg)

		res := createResponse(oHeadDeltaUpdate)
		writeUserInfo(activated, res)
		writeUserInfo(deactivated, res)
		writeUserId(client.user.id, res)
		writeString(*msg, res)

		broadcastToAll(res)
	}

	for {
		select {
		// Client joined
		case client := <-server.inbound:
			server.clients[client] = true
			client.user = newUserInfo(nextUserID, nil)
			nextUserID++
			fmt.Println("new user registered")
		// Client left
		case client := <-server.outbound:
			delete(server.clients, client)
			close(client.responses)
			fmt.Println("handling disconnection")
		// Message received
		case req := <-server.requests:
			fmt.Println("message received", req.head)
			handler, ok := requestHandlers[req.head]

			if ok {
				handler(bytes.NewBuffer(req.bytes), req)
			} else {
				fmt.Println("ERROR: Invalid message header '" + strconv.Itoa((int)(req.head)) + "' received!")
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

	client := server.createClient(connection)
	fmt.Println("Client connected from '" + client.connection.RemoteAddr().String() + "'!")
	server.inbound <- client
	client.init()
}
