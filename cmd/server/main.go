package main

import (
	"log"

	"github.com/babadreamz/Simple-Chat-Service/internal/config"
	"github.com/babadreamz/Simple-Chat-Service/internal/database"
	"github.com/babadreamz/Simple-Chat-Service/internal/handlers"
	"github.com/babadreamz/Simple-Chat-Service/internal/websocket"
	"github.com/gin-gonic/gin"
)

func main() {
	configs := config.Load()
	database.Connect(
		configs.MongoHost,
		configs.MongoPort,
		configs.MongoUser,
		configs.MongoPass,
	)

	hub := websocket.NewTrafficHub()
	go hub.Run()

	router := gin.Default()

	handlers.SetupRoutes(router, hub)
	log.Printf("Chat service started on %s:", configs.AppPort)
	if err := router.Run(":" + configs.AppPort); err != nil {
		log.Fatal("Server failed to start", err)
	}
}
