package server

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type client struct {
	connection *websocket.Conn
	messagesTo chan []byte
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
		_, message, err := client.connection.ReadMessage()

		if err != nil {
			fmt.Println("ERROR: Failed to read a message from client!")
			fmt.Println(err.Error())

			return
		}

		client.server.messages <- message
	}
}

func (client *client) sendMessages() {
	for message := range client.messagesTo {
		err := client.connection.WriteMessage(websocket.BinaryMessage, message)

		if err != nil {
			fmt.Println("ERROR: Failed to write a message to client!")
			fmt.Println(err.Error())

			return
		}
	}

	client.connection.Close()
}
