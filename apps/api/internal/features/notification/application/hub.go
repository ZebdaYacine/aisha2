package application

import (
	"sync"

	"github.com/aisha-platform/aisha/apps/api/internal/features/notification/domain"
	"github.com/fasthttp/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]struct{}
}

type Client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func NewHub() *Hub {
	return &Hub{clients: make(map[string]map[*Client]struct{})}
}

func (h *Hub) Add(userID string, conn *websocket.Conn) *Client {
	client := &Client{conn: conn}
	h.mu.Lock()
	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*Client]struct{})
	}
	h.clients[userID][client] = struct{}{}
	h.mu.Unlock()
	return client
}

func (h *Hub) Remove(userID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	clients := h.clients[userID]
	delete(clients, client)
	if len(clients) == 0 {
		delete(h.clients, userID)
	}
}

func (h *Hub) Publish(userID string, notification domain.Notification) {
	if userID == "" {
		return
	}
	h.mu.RLock()
	clients := make([]*Client, 0, len(h.clients[userID]))
	for client := range h.clients[userID] {
		clients = append(clients, client)
	}
	h.mu.RUnlock()
	for _, client := range clients {
		if err := client.WriteJSON(map[string]any{"type": "notification", "notification": notification}); err != nil {
			_ = client.Close()
		}
	}
}

func (c *Client) WriteJSON(value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(value)
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.Close()
}
