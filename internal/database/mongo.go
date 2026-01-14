package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect(host, port, user, password string) {
	credentials := options.Credential{
		AuthSource: "admin",
		Username:   user,
		Password:   password,
	}
	uri := fmt.Sprintf("mongodb://%s:%s", host, port)
	clientOptions := options.Client().ApplyURI(uri).SetAuth(credentials)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDb", err)
	}
	fmt.Println("Connected to MongoDB")
	Client = client
}
func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("chat_db").Collection(collectionName)
}
