package goappsocket

import (
	"encoding/json"

	"github.com/angadthandi/golangmongoapp/golangapp/jsondefinitions"
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

	// Buffered channel of inbound messages.
	ChClientCorrelationIds chan []byte

	// Buffered channel of inbound messages.
	ChSendMsgToClientWithCorrelationId chan []byte
}

func NewHub() *Hub {
	return &Hub{
		broadcast:                          make(chan []byte),
		register:                           make(chan *Client),
		unregister:                         make(chan *Client),
		clients:                            make(map[*Client]bool),
		ChClientCorrelationIds:             make(chan []byte),
		ChSendMsgToClientWithCorrelationId: make(chan []byte),
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
		case b, ok := <-h.ChClientCorrelationIds:
			if ok {
				setClientCorrelationId(h, b)
			}
		case b, ok := <-h.ChSendMsgToClientWithCorrelationId:
			if ok {
				sendMsgToClientWithCorrelationId(h, b)
			}
		}
	}
}

// SendMsgToAllClients sends message to
// all connected clients
func (h *Hub) SendMsgToAllClients(
	jsonMsg json.RawMessage,
) {
	log.Debugf("hub SendMsgToAllClients jsonMsg: %v", jsonMsg)
	h.broadcast <- jsonMsg
}

// setClientCorrelationId sets correlationId for client
// in clientCorrelationIds map
//
// This func does writes to client data
// so should only called via goroutine running hub
// so its unexported
// if exported, can result in DATA RACE CONDITIONS
func setClientCorrelationId(h *Hub, b []byte) {
	log.Debug("setClientCorrelationId")

	var clientUUIDCorrId jsondefinitions.ClientUUIDCorrelationID
	err := json.Unmarshal(b, &clientUUIDCorrId)
	if err != nil {
		log.Debugf("setClientCorrelationId Unable to unmarshal: %v",
			err)
		return
	}
	log.Debugf("setClientCorrelationId clientUUIDCorrId: %v",
		clientUUIDCorrId)

	for c, _ := range h.clients {
		if c.ClientUUID == clientUUIDCorrId.ClientUUID {
			c.clientCorrelationIds[clientUUIDCorrId.ClientCorrelationId] = true
		}
	}
}

// sendMsgToClientWithCorrelationId
// sends message to client
// which has a correlationId stored
// in its clientCorrelationIds map
//
// This func does read & writes to client data
// so should only called via goroutine running hub
// so its unexported
// if exported, can result in DATA RACE CONDITIONS
func sendMsgToClientWithCorrelationId(h *Hub, b []byte) {
	log.Debug("sendMsgToClientWithCorrelationId")

	var msg jsondefinitions.MicroServiceResponseMsgForHub
	err := json.Unmarshal(b, &msg)
	if err != nil {
		log.Debugf("sendMsgToClientWithCorrelationId Unable to unmarshal: %v",
			err)
		return
	}
	log.Debugf("sendMsgToClientWithCorrelationId received msg: %v",
		msg)

	for c, _ := range h.clients {
		if _, ok := c.clientCorrelationIds[msg.CorrelationId]; ok {
			c.send <- msg.ReceivedJsonMsg

			delete(c.clientCorrelationIds, msg.CorrelationId)

			break
		}
	}
}
