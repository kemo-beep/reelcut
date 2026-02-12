package handler

import (
	"log"
	"net/http"
	"time"

	"reelcut/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler handles WebSocket connections and registers them with the hub.
type WebSocketHandler struct {
	Hub *Hub
}

// NewWebSocketHandler creates a WebSocket handler that uses the given hub.
func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{Hub: hub}
}

// Handle upgrades the request to WebSocket and runs read/write pumps.
func (h *WebSocketHandler) Handle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	connection := &Connection{
		UserID: userID,
		Send:   make(chan []byte, 256),
	}
	h.Hub.Register(userID, connection)
	defer h.Hub.Unregister(userID, connection)

	go h.writePump(conn, connection)
	h.readPump(conn, connection)
}

func (h *WebSocketHandler) readPump(conn *websocket.Conn, c *Connection) {
	defer conn.Close()
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket read: %v", err)
			}
			break
		}
	}
}

func (h *WebSocketHandler) writePump(conn *websocket.Conn, c *Connection) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
