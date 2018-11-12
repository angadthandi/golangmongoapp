package goappsocket

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	log.Debug("Hub Run!")
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Debugf("Hub register client: %v", client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Debugf("Hub unregister client: %v", client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					log.Debugf("Hub send message client: %v", message)
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) SendMsgToAllClients(
	jsonMsg json.RawMessage,
) {
	log.Debugf("hub SendMsgToAllClients h: %v", h)
	log.Debugf("hub SendMsgToAllClients jsonMsg: %v", jsonMsg)
	h.broadcast <- jsonMsg
}
