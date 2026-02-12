package handler

import (
	"context"
	"encoding/json"
	"sync"

	"reelcut/internal/domain"
	"reelcut/internal/notifier"
)

// Connection holds a single WebSocket connection and its send channel.
type Connection struct {
	UserID   string
	Send     chan []byte
	connKey  string // optional unique key for this connection
}

// Hub holds registered connections and broadcasts messages per user.
type Hub struct {
	mu          sync.RWMutex
	connections map[string]map[*Connection]struct{} // userID -> set of connections
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		connections: make(map[string]map[*Connection]struct{}),
	}
}

// Register adds a connection for the given userID.
func (h *Hub) Register(userID string, c *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.connections[userID] == nil {
		h.connections[userID] = make(map[*Connection]struct{})
	}
	h.connections[userID][c] = struct{}{}
}

// Unregister removes a connection.
func (h *Hub) Unregister(userID string, c *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if m := h.connections[userID]; m != nil {
		delete(m, c)
		close(c.Send)
		if len(m) == 0 {
			delete(h.connections, userID)
		}
	}
}

// BroadcastToUser sends a message to all connections for the given user.
func (h *Hub) BroadcastToUser(userID string, message []byte) {
	h.mu.RLock()
	conns := make([]*Connection, 0, len(h.connections[userID]))
	for c := range h.connections[userID] {
		conns = append(conns, c)
	}
	h.mu.RUnlock()
	for _, c := range conns {
		select {
		case c.Send <- message:
		default:
			// skip if channel full
		}
	}
}

// JobNotifierImpl sends job updates to the user over the hub.
type JobNotifierImpl struct {
	Hub *Hub
}

// NewJobNotifier returns a JobNotifier that broadcasts job_updated to the hub.
func NewJobNotifier(hub *Hub) notifier.JobNotifier {
	return &JobNotifierImpl{Hub: hub}
}

// NotifyJob implements notifier.JobNotifier.
func (n *JobNotifierImpl) NotifyJob(ctx context.Context, job *domain.ProcessingJob) {
	if n == nil || n.Hub == nil || job == nil {
		return
	}
	payload := map[string]interface{}{
		"type": "job_updated",
		"job":  job,
	}
	b, _ := json.Marshal(payload)
	n.Hub.BroadcastToUser(job.UserID.String(), b)
}
