package models

import (
	"time"
)

type ConversationStatus string

const (
	StatusActive   ConversationStatus = "active"
	StatusClosed   ConversationStatus = "closed"
	StatusArchived ConversationStatus = "archived"
)

type Role string

const (
	RoleResponder Role = "RESPONDER"
	RoleReporter  Role = "REPORTER"
)

type Conversation struct {
	ID              string             `bson:"_id" json:"id"`
	Participants    []Participant      `bson:"participants" json:"participants"`
	LastMessageTime time.Time          `bson:"last_message_time" json:"last_message_time"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	Status          ConversationStatus `bson:"status" json:"status"`
}
type Participant struct {
	UserId         string    `bson:"user_id" json:"user_id"`
	ConversationId string    `bson:"conversation_id" json:"conversation_id"`
	Role           Role      `bson:"role" json:"role"`
	JoinedAt       time.Time `bson:"joined_at" json:"joined_at"`
}
type Message struct {
	ID             string    `bson:"_id,omitempty" json:"id"`
	ConversationID string    `bson:"conversation_id" json:"conversation_id"`
	SenderID       string    `bson:"sender_id" json:"sender_id"`
	Content        string    `bson:"content" json:"content"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
}
type IncomingMessage struct {
	SenderID       string `json:"sender_id"`
	Content        string `json:"content"`
	ConversationID string `json:"conversation_id"`
}
