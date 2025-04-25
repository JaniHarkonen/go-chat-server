package main

import (
	"fmt"
	"net/http"

	"github.com/JaniHarkonen/go-chat-server/internal/server"
)

func main() {
	const address string = "localhost"
	const port string = "12345"

	server := server.NewServer()
	http.Handle("/chat", server)
	fmt.Println("Server started!")
	go server.Run()

	fmt.Println("Listening to port " + port + "...")

	if err := http.ListenAndServe(address+":"+port, nil); err != nil {
		fmt.Println("ERROR: Unable to start server! Shutting down...")
		fmt.Println(err.Error())
	}
}
