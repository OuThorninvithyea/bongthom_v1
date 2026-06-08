package websocket

import (

	// Community pacakges
	"encoding/base64"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

var global_clients = make(map[string]*Client)

type Client struct {
	Conn *websocket.Conn
	ID   string
}

type WebSocketManager struct {
	// Clients map[string]*Client
	mu sync.RWMutex
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		// Clients: make(map[string]*Client),
	}
}

func (wm *WebSocketManager) PrintlnClient() {
	fmt.Println("admin client for websocket", global_clients)
}

func (wm *WebSocketManager) AddClient(client *Client) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	global_clients[client.ID] = client
	log.Printf("Client added: %s", client.ID)

}

func (wm *WebSocketManager) RemoveClient(clientId string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	delete(global_clients, clientId)
	log.Printf("Client removed: %s", clientId)
}

func (wm *WebSocketManager) Broadcast(data interface{}) {

	wm.mu.RLock()
	defer wm.mu.RUnlock()
	wm.PrintlnClient()

	jsonData, ok := data.(string)
	if ok {
		// Check if it's Base64 encoded
		decodedBytes, err := base64.StdEncoding.DecodeString(jsonData)
		if err == nil {
			data = string(decodedBytes)
		}
	}

	fmt.Println("data")

	for _, client := range global_clients {

		if err := client.Conn.WriteJSON(data); err != nil {
			log.Printf("Broadcast error for client %s: %v", client.ID, err)
			client.Conn.Close()
			delete(global_clients, client.ID)
		}
	}
}

func (manager *WebSocketManager) Emit(clientID string, data interface{}) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	client, ok := global_clients[clientID]
	if !ok {
		log.Printf("Client %s not found", clientID)
		return
	}

	err := client.Conn.WriteJSON(data)
	if err != nil {
		log.Printf("Error sending message to client %s: %v", clientID, err)
		client.Conn.Close()
		delete(global_clients, clientID)
	}
}

func (wm *WebSocketManager) NotifyUser(userID string, data interface{}) {
	wm.PrintlnClient()
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	wm.PrintlnClient()
	log.Printf("Clients available: %+v", global_clients)

	clientID := fmt.Sprintf("user-%s", userID)
	client, ok := global_clients[clientID]
	if !ok {
		log.Printf("Client %s not connected", clientID)
		return
	}

	err := client.Conn.WriteJSON(data)
	if err != nil {
		log.Printf("Error sending notification to client %s: %v", clientID, err)
		client.Conn.Close()
		delete(global_clients, clientID)
	}
}
