package wss

import (
	"bifrost/common/jwtx"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	period = 9

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * period) / 10

	// send buffer size
	bufSize = 128

	bufferSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  bufferSize,
	WriteBufferSize: bufferSize,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	Ud jwtx.UserData
}

type WssMessage struct {
	SenderKey WssSendTag `json:"sender"`    // 发送人
	Recipient WssSendTag `json:"recipient"` // 接收人
	Content   string     `json:"content"`   // 发送数据
}

type WssSendTag struct {
	Key   string
	Value interface{}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()

		_ = c.Conn.Close()

		c.Hub.unregister <- c
	}()

	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// The hub closed the channel.
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})

				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.Send)
			}

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}

			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(ud map[string]interface{}, hub *Hub, w http.ResponseWriter, r *http.Request) (err error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	client := &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, bufSize),
		Ud:   ud,
	}
	client.Hub.register <- client

	// 只有读取的权限
	go client.writePump()
	//go client.readPump()

	return err
}
