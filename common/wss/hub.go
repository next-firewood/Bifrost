package wss

import (
	"bifrost/common/errorx"
	"bifrost/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Hub struct {
	// Registered clients.
	clients      map[*Client]bool
	reverseIndex map[string]map[interface{}][]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub(ck string) *Hub {
	return &Hub{
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[*Client]bool),
		reverseIndex: make(map[string]map[interface{}][]*Client),
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
			for fieldName, fieldValue := range client.Ud {
				if _, exists := h.reverseIndex[fieldName]; !exists {
					h.reverseIndex[fieldName] = make(map[interface{}][]*Client)
				}
				h.reverseIndex[fieldName][fieldValue] = append(h.reverseIndex[fieldName][fieldValue], client)
			}
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

func (h *Hub) ClientsMap(fieldName string, value string) (clients []*Client, err error) {
	fieldIndex, exists := h.reverseIndex[fieldName]
	if !exists {
		return nil, errorx.BusinessErr(fmt.Sprintf("字段 %s 不存在于反向索引中", fieldName))
	}

	var convertedValue interface{}
	var found bool

	for key := range fieldIndex {
		switch key.(type) {
		case int:
			convertedValue, err = strconv.Atoi(value)
			if err == nil {
				found = true
			}
		case uint64:
			convertedValue, err = strconv.ParseUint(value, 10, 64)
			if err == nil {
				found = true
			}
		case uint32:
			convertedValue, err = strconv.ParseUint(value, 10, 32)
			if err == nil {
				found = true
			}
		case float64:
			convertedValue, err = strconv.ParseFloat(value, 64)
			if err == nil {
				found = true
			}
		case float32:
			convertedValue, err = strconv.ParseFloat(value, 32)
			if err == nil {
				found = true
			}
		case string:
			convertedValue = value
			found = true
		}
		if found {
			break
		}
	}

	if !found {
		return nil, errorx.BusinessErr(fmt.Sprintf("无法将值 %v 转换为字段 %s 的类型", value, fieldName))
	}

	clients, exists = fieldIndex[convertedValue]
	if !exists {
		return nil, errorx.BusinessErr(fmt.Sprintf("未找到字段 %s 中值为 %v 的客户端", fieldName, convertedValue))
	}

	return clients, err
}

func (h *Hub) FilterSendFn(c *gin.Context) (fn func(message string) error, err error) {
	var f model.FilterHeader

	if err = c.ShouldBindHeader(&f); err != nil {
		return nil, err
	}

	if f.FilterValue == "" || f.FilterKey == "" {
		return nil, errorx.NewGinBindParamError()
	}

	clients, err := h.ClientsMap(f.FilterKey, f.FilterValue)

	fn = func(message string) (err error) {
		for _, client := range clients {
			select {
			case client.Send <- []byte(message):
			default:
				close(client.Send)
				delete(h.clients, client)
				return errorx.BusinessErr("发送消息失败")
			}
		}

		return err
	}

	return fn, err
}
