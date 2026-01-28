package wss

import (
	"bytes"
	"encoding/json"
)

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

func (h *Hub) ReadPump(msg WssMessage) error {
	message, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	h.broadcast <- message

	return nil
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
		case message := <-h.broadcast:
			var msg WssMessage

			err := json.Unmarshal(message, &msg)
			if err != nil {
				continue
			}

			send(h.clients, msg)
		}
	}
}

func send(clients map[*Client]bool, msg WssMessage) {
	for client := range clients {

		if msg.Recipient.Key != "" && client.Ud[msg.Recipient.Key] != msg.Recipient.Value {
			continue
		}

		select {
		case client.Send <- []byte(msg.Content):
		default:
			close(client.Send)
			delete(clients, client)
		}
	}
}

func (h *Hub) Clients() map[*Client]bool {
	return h.clients
}
