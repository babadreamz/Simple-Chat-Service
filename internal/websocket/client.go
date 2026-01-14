package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/babadreamz/Simple-Chat-Service/internal/database"
	"github.com/babadreamz/Simple-Chat-Service/internal/models"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Hub            *TrafficHub
	Conn           *websocket.Conn
	Send           chan []byte
	ConversationId string
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.UnRegister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, rawData, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var incomingMsg models.IncomingMessage
		if err := json.Unmarshal(rawData, &incomingMsg); err != nil {
			log.Printf("error parsing JSON: %v", err)
			continue
		}
		savedMessage, err := database.SaveMessage(incomingMsg)
		if err != nil {
			log.Printf("error saving message to DB: %v", err)
			continue
		}
		broadcastBytes, _ := json.Marshal(savedMessage)
		c.Hub.Broadcast <- broadcastBytes
	}
}
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			responseWriter, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = responseWriter.Write(message)
			if err != nil {
				return
			}
			if err := responseWriter.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
func ServeWs(hub *TrafficHub, responseWriter http.ResponseWriter, request *http.Request) {
	conversationID := request.URL.Query().Get("conversation_id")
	if conversationID == "" {
		http.Error(responseWriter, "conversation_id is required", http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		Hub:            hub,
		Conn:           conn,
		Send:           make(chan []byte, 256),
		ConversationId: conversationID,
	}
	client.Hub.Register <- client
	go client.WritePump()
	go client.ReadPump()
}
