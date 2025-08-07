package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type Message struct {
	Type string          `json:"type"`
	From string          `json:"from"`
	To   string          `json:"to"`
	Data json.RawMessage `json:"data,omitempty"`
}

var (
	clients   = make(map[string]*Client)
	clientsMu sync.Mutex
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func broadcastUserList() {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	var ids []string
	for id := range clients {
		ids = append(ids, id)
	}

	data, _ := json.Marshal(map[string]interface{}{
		"type": "user-list",
		"data": ids,
	})
	fmt.Println(string(data))

	for _, client := range clients {
		client.Conn.WriteMessage(websocket.TextMessage, data)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	if r.Method == "OPTIONS" {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	var client *Client

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected: ", err)
			if client != nil {
				clientsMu.Lock()
				delete(clients, client.ID)
				clientsMu.Unlock()
				broadcastUserList()
			}
			break
		}

		var msg Message
		fmt.Println("Received message:", string(message))
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Println("Json decode error:", err)
			continue
		}

		switch msg.Type {
		case "register":
			clientsMu.Lock()

			if _, exists := clients[msg.From]; exists {
				errMsg := map[string]string{
					"type": "error",
					"data": "UAE",
				}

				errJSON, _ := json.Marshal(errMsg)
				conn.WriteMessage(websocket.TextMessage, errJSON)

				clientsMu.Unlock()
				broadcastUserList()
				return
			}

			client = &Client{
				ID:   msg.From,
				Conn: conn,
			}

			clients[msg.From] = client
			clientsMu.Unlock()
			broadcastUserList()

		case "offer", "answer", "candidate":
			clientsMu.Lock()
			targetClient, ok := clients[msg.To]
			clientsMu.Unlock()

			if !ok {
				fmt.Println("User %s not found", msg.To)
				continue
			}

			targetClient.Conn.WriteMessage(websocket.TextMessage, message)
		}

	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Server listening on port 5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

}
