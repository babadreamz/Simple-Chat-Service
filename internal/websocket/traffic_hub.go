package websocket

import "encoding/json"

type TrafficHub struct {
	//Clients    map[*Client]bool
	Rooms      map[string]map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	UnRegister chan *Client
}

func NewTrafficHub() *TrafficHub {
	return &TrafficHub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		Rooms:      make(map[string]map[*Client]bool),
	}
}
func (hub *TrafficHub) Run() {
	for {
		select {
		case client := <-hub.Register:
			if hub.Rooms[client.ConversationId] == nil {
				hub.Rooms[client.ConversationId] = make(map[*Client]bool)
			}
			hub.Rooms[client.ConversationId][client] = true

		case client := <-hub.UnRegister:
			if _, ok := hub.Rooms[client.ConversationId][client]; ok {
				delete(hub.Rooms[client.ConversationId], client)
				close(client.Send)
			}
			if len(hub.Rooms[client.ConversationId]) == 0 {
				delete(hub.Rooms, client.ConversationId)
			}

		case message := <-hub.Broadcast:
			type MessageHeader struct {
				ConversationId string `json:"conversation_id"`
			}
			var messageHeader MessageHeader
			err := json.Unmarshal(message, &messageHeader)
			if err != nil {
				return
			}
			if clientsInRoom, ok := hub.Rooms[messageHeader.ConversationId]; ok {
				for client := range clientsInRoom {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(clientsInRoom, client)
					}
				}
			}
		}
	}
}
