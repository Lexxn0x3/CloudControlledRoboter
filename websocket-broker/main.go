package main

import (
	"bufio"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[*websocket.Conn]bool)
)

func main() {
	log.Println("Starting server...")

	// Set up a simple HTTP server
	http.HandleFunc("/", handleConnections)
	go func() {
		log.Println("HTTP server listening on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	// Connect to the TCP server
	log.Println("Connecting to TCP server at 127.0.0.1:9001...")
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Fatal("Connection to TCP server failed:", err)
	}
	defer conn.Close()
	log.Println("Connected to TCP server successfully.")

	// Read from the TCP connection
	go readTCPStream(conn)
	select {} // Keep the main goroutine running
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("Received new WebSocket connection request.")

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer ws.Close()

	log.Printf("WebSocket connection established: %s\n", ws.RemoteAddr())
	clients[ws] = true

	// Infinite loop to keep the connection open
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("WebSocket error from %s: %v\n", ws.RemoteAddr(), err)
			delete(clients, ws)
			break
		}
	}
}

func readTCPStream(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		//log.Println("Received data from TCP server.")
		message := scanner.Bytes()
		go broadcastToClients(message)
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading from TCP connection:", err)
	}
}

func broadcastToClients(message []byte) {
	//log.Printf("Broadcasting message to %d clients.\n", len(clients))
	for client := range clients {
		err := client.WriteMessage(websocket.BinaryMessage, message)
		if err != nil {
			log.Printf("Error broadcasting to client %s: %v\n", client.RemoteAddr(), err)
			client.Close()
			delete(clients, client)
		}
	}
}
