package dtos

import "time"

type ParticipantDTO struct {
	UserId string `json:"userId"`
	Role   string `json:"role"`
}
type ConversationDTO struct {
	ID           string           `json:"conversationId"`
	Participants []ParticipantDTO `json:"participants"`
	CreatedAt    time.Time        `json:"createdAt"`
	Status       string           `json:"status"`
}
type CreateConversationRequest struct {
	ConversationID string `json:"conversation_id"`
	ResponderID    string `json:"responder_id"`
	ReporterID     string `json:"reporter_id"`
}
