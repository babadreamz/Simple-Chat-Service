package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/babadreamz/Simple-Chat-Service/internal/dtos"
	"github.com/babadreamz/Simple-Chat-Service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveMessage(incoming models.IncomingMessage) (*models.Message, error) {
	isActive, err := IsConversationActive(incoming.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to check conversation status: %v", err)
	}
	if !isActive {
		return nil, fmt.Errorf("conversation is not active or does not exist")
	}
	collection := GetCollection("messages")
	newMessage := models.Message{
		SenderID:       incoming.SenderID,
		Content:        incoming.Content,
		ConversationID: incoming.ConversationID,
		CreatedAt:      time.Now(),
	}
	resul, err := collection.InsertOne(context.TODO(), newMessage)
	if err != nil {
		return nil, err
	}
	if id, ok := resul.InsertedID.(primitive.ObjectID); ok {
		newMessage.ID = id.Hex()
	}
	return &newMessage, nil
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
