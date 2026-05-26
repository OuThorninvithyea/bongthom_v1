package websocket

import (
	// Community Pacakges
	"fmt"
	"log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	// Intenral Pacakges
	"admin-api/pkg/share"
)

type WebSocketRoute struct {
	manager *WebSocketManager
	db      *sqlx.DB
}

func NewRoute(app *fiber.App, db *sqlx.DB) *WebSocketRoute {
	manager := NewWebSocketManager()
	route := &WebSocketRoute{manager: manager, db: db}
	v1 := app.Group("/api/v1")
	v1.Get("/ws", websocket.New(route.HandleConnection))

	return route
}

func (route *WebSocketRoute) HandleConnection(conn *websocket.Conn) {
	defer conn.Close()

	userContext := conn.Locals("UserContext")
	if userContext == nil {
		conn.WriteMessage(websocket.CloseMessage, []byte("Unauthorized"))
		return
	}

	userID := userContext.(share.UserContext).UserID
	clientID := fmt.Sprintf("user-%v", userID)

	client := &Client{Conn: conn, Id: clientID}
	route.manager.AddClient(client)

	defer route.manager.RemoveClient(clientID)

	for {
		_, message, err := conn.ReadMessage()
		fmt.Println(string(message), "this is messs>>>>>>>>>>>")
		route.manager.Broadcast(message)

		if err != nil {
			log.Printf("Error reading message from client %s: %v", clientID, err)
			break
		}
		log.Printf("Message from client %s: %s", clientID, message)
	}
}
