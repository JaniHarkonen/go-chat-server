package server

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/JaniHarkonen/go-chat-server/internal/chat"
	"github.com/gorilla/websocket"
)

type Server struct {
	clients     map[chat.UserID]*client
	inbound     chan *client
	outbound    chan *client
	requests    chan *request
	chatManager *chat.Manager
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
		clients:     make(map[chat.UserID]*client),
		inbound:     make(chan *client),
		outbound:    make(chan *client),
		requests:    make(chan *request),
		chatManager: nil,
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
	chatManager := chat.NewManager(10, 10)
	server.chatManager = chatManager
	nextUserID := firstUserID

	// Send a given message to all clients
	broadcastToAll := func(buffer *bytes.Buffer) {
		bytes := buffer.Bytes()

		for _, client := range server.clients {
			client.send(bytes)
		}
	}

	executeCommand := commandExecutor(server)
	requestHandlers := make(map[header]func(in *bytes.Buffer, req *request))

	// Handle first connection and exchanging of user info
	requestHandlers[iHeadClientInfo] = func(in *bytes.Buffer, req *request) {
		client := req.client
		name := readString(in)
		client.user.SetName(&name)
		chatManager.RegisterUser(client.user) // Register the user's username

		// Write a table of active users paired with their IDs
		res := createResponse(oHeadCompleteUpdate)
		writeUInt32((uint32)(len(chatManager.ActiveUsers())), res) // Active user count

		for active := range chatManager.ActiveUsers() {
			writeUserInfo(active, res)
		}

		// Write a snapshot of latest chat messages
		visible := chatManager.VisibleMessages()
		writeUInt32((uint32)(len(chatManager.Snapshot())), res)

		for _, msg := range visible {
			writeUserId(msg.User().ID(), res)
			writeString(*msg.Message(), res)
		}

		client.send(res.Bytes())
	}

	// Handle user name change
	requestHandlers[iHeadNameChange] = func(in *bytes.Buffer, req *request) {
		client := req.client
		name := readString(in)

		// Unregister then re-register the user after its name has been changed to reflect the name change
		chatManager.UnregisterUser(client.user)
		client.user.SetName(&name)
		chatManager.RegisterUser(client.user)

		// Only broadcast name change if the user is active
		if chatManager.IsUserActive(client.user) {
			res := createResponse(oHeadNameChange)
			writeUserId(client.user.ID(), res)
			writeString(name, res)

			broadcastToAll(res)
		}
	}

	// Handle chat input (message/command)
	requestHandlers[iHeadChatInput] = func(in *bytes.Buffer, req *request) {
		client := req.client
		msg := readString(in)

		// Handle possible command
		if (msg)[0] == '/' {
			executeCommand(&msg)
		} else {
			// Not a command -> handle chat message
			if !chatManager.IsUserMuted(client.user) {
				activated, deactivated := chatManager.Post(client.user, &msg)

				res := createResponse(oHeadDeltaUpdate)
				writeUserInfo(activated, res)

				if deactivated != nil {
					writeUserId(deactivated.ID(), res)
				} else {
					writeUserId(0, res)
				}

				writeUserId(client.user.ID(), res)
				writeString(msg, res)

				broadcastToAll(res)
			}
		}
	}

	for {
		select {
		// Client joined
		case client := <-server.inbound:
			server.clients[nextUserID] = client
			client.user = chat.NewUser(nextUserID, nil)
			nextUserID++
		// Client left
		case client := <-server.outbound:
			fmt.Println("disconnected client")
			server.chatManager.UnregisterUser(client.user)
			delete(server.clients, client.user.ID())
			close(client.responses)
		// Message received
		case req := <-server.requests:
			handler, ok := requestHandlers[req.head]

			if ok {
				handler(bytes.NewBuffer(req.bytes), req)
			} else {
				fmt.Println("ERROR: Invalid message header '" + strconv.Itoa((int)(req.head)) + "' received!")
			}
		}
	}
}

func (server *Server) ResolveClient(user *chat.User) *client {
	if user == nil {
		return nil
	}

	client, ok := server.clients[user.ID()]

	if ok {
		return client
	}

	return nil
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
