package handlers

import (
	"fmt"
	"net/http"

	"github.com/babadreamz/Simple-Chat-Service/internal/database"
	"github.com/babadreamz/Simple-Chat-Service/internal/dtos"
	"github.com/babadreamz/Simple-Chat-Service/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateConversation(context *gin.Context) {
	var convRequest dtos.CreateConversationRequest
	if err := context.ShouldBindJSON(&convRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Invalid request body": err.Error()})
		return
	}
	fmt.Printf("Creating Conversation with ID: %s\n", convRequest.ConversationID)
	convoDto, err := database.CreateConversation(convRequest.ConversationID, convRequest.ResponderID, convRequest.ReporterID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Failed to create conversation": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, convoDto)
}
func CloseConversation(context *gin.Context) {
	id := context.Query("conversation_id")
	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{"Bad Request": "conversation_id is required"})
		return
	}
	err := database.UpdateConversationStatus(id, models.StatusClosed)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Failed to close conversation": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Conversation closed successfully"})
}
func ArchiveConversation(context *gin.Context) {
	id := context.Query("conversation_id")
	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{"Bad Request": "conversation_id is required"})
		return
	}
	err := database.UpdateConversationStatus(id, models.StatusArchived)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Failed to archive conversation": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Conversation archived successfully"})
}
