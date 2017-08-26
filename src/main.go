package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// global variables for connected clients and broadcast channel
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

// upgrader config
var upgrader = websocket.Upgrader{}

// object definition
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	// main fileserver
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// configure websockets route
	http.HandleFunc("/ws", handleConnections)

	// listen to incoming chat messages
	go handleMessages()

	// start server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// close connection when function returns
	defer ws.Close()

	// register client into map
	clients[ws] = true

	// listening loop
	for {
		var msg Message
		// read in new message as JSON and map it to a message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// send new message to broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// grab next message from broadcast channel
		msg := <-broadcast
		// send it to all currently connected clients
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
