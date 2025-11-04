package websocket

import (
	"sync"
)

type Client struct {
	hub  *Hub
	conn *Conn
	send chan []byte
}

type Hub struct {
	Clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case Client := <-h.register:
			h.mutex.Lock()
			h.Clients[Client] = true
			h.mutex.Unlock()

		case Client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.Clients[Client]; ok {
				delete(h.Clients, Client)
				close(Client.send)
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			h.mutex.RLock()
			for Client := range h.Clients {
				select {
				case Client.send <- message:
				default:
					close(Client.send)
					delete(h.Clients, Client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.Clients)
}
