package database

import (
	"context"
	"errors"
	"time"

	"github.com/babadreamz/Simple-Chat-Service/internal/dtos"
	"github.com/babadreamz/Simple-Chat-Service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveMessage(message *models.Message) error {
	collection := GetCollection("messages")
	_, err := collection.InsertOne(context.TODO(), message)
	return err
}
func CreateConversation(conversationId, responderId, reporterId string) (*dtos.ConversationDTO, error) {
	collection := GetCollection("conversations")
	participants := []models.Participant{
		{
			UserId:         responderId,
			ConversationId: conversationId,
			Role:           models.RoleResponder,
			JoinedAt:       time.Now(),
		},
		{
			UserId:         reporterId,
			ConversationId: conversationId,
			Role:           models.RoleReporter,
			JoinedAt:       time.Now(),
		},
	}
	newConversation := models.Conversation{
		ID:              conversationId,
		Participants:    participants,
		LastMessageTime: time.Now(),
		CreatedAt:       time.Now(),
		Status:          models.StatusActive,
	}
	_, err := collection.InsertOne(context.TODO(), newConversation)
	if err != nil {
		return nil, err
	}
	ParticipantsDtos := []dtos.ParticipantDTO{
		{
			UserId: responderId, Role: string(models.RoleResponder)},
		{UserId: reporterId, Role: string(models.RoleReporter)},
	}
	ConversationDto := &dtos.ConversationDTO{
		ID:           conversationId,
		Participants: ParticipantsDtos,
		CreatedAt:    time.Now(),
		Status:       string(newConversation.Status),
	}
	return ConversationDto, nil
}
func UpdateConversationStatus(conversationId string, newStatus models.ConversationStatus) error {
	collection := GetCollection("conversations")
	filter := bson.M{"_id": conversationId}
	update := bson.M{"$set": bson.M{"newStatus": newStatus}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}
func IsConversationActive(conversationId string) (bool, error) {
	collection := GetCollection("conversations")
	var convo models.Conversation
	err := collection.FindOne(context.TODO(), bson.M{"_id": conversationId}).Decode(&convo)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return convo.Status == models.StatusActive, nil
}
