package handlers

import (
	"github.com/babadreamz/Simple-Chat-Service/internal/websocket"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, hub *websocket.TrafficHub) {
	router.GET("/ws", func(context *gin.Context) {
		websocket.ServeWs(hub, context.Writer, context.Request)
	})
	conversation := router.Group("/conversation")
	{
		conversation.POST("/start", CreateConversation)
		conversation.PATCH("/close", CloseConversation)
		conversation.PATCH("/archive", ArchiveConversation)
	}
}
