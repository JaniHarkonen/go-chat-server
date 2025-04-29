package server

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type client struct {
	connection *websocket.Conn
	user       *userInfo
	responses  chan []byte
	server     *Server
}

func (client *client) init() {
	// Use two threads, one for reading and another for writing
	go client.sendMessages()
	client.receiveMessages()
	client.server.outbound <- client
}

func (client *client) receiveMessages() {
	defer client.connection.Close() // Defer because the code is otherwise flagged as unreachable

	for {
		_, bytes, err := client.connection.ReadMessage()

		if err != nil {
			fmt.Println("ERROR: Failed to read a message from client!")
			fmt.Println(err.Error())

			return
		}

		client.server.requests <- &request{
			head:   (header)(bytes[0]),
			client: client,
			bytes:  bytes[1:],
		}
	}
}

func (client *client) sendMessages() {
	for req := range client.server.requests {
		err := client.connection.WriteMessage(websocket.BinaryMessage, req.bytes)

		if err != nil {
			fmt.Println("ERROR: Failed to write a message to client!")
			fmt.Println(err.Error())

			return
		}
	}

	client.connection.Close()
}

func (client *client) send(bytes []byte) {
	client.responses <- bytes
}
