package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	MongoURI  string
	MongoHost string
	MongoPort string
	MongoUser string
	MongoPass string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		mongoUser := getEnv("MONGO_DB_ROOT_USERNAME", "admin")
		mongoHost := getEnv("MONGO_HOST", "localhost")
		mongoExternalPort := getEnv("MONGO_EXTERNAL_PORT", "27017")
		mongoPassword := getEnvRequired("MONGO_DB_ROOT_PASSWORD")

		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin",
			mongoUser, mongoPassword, mongoHost, mongoExternalPort)
	}

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = getEnv("APP_PORT", "2020")
	}

	return &Config{
		AppPort:   appPort,
		MongoURI:  uri,
		MongoHost: getEnv("MONGO_HOST", "localhost"),
		MongoPort: getEnv("MONGO_EXTERNAL_PORT", "27017"),
		MongoUser: getEnv("MONGO_DB_ROOT_USERNAME", "admin"),
		MongoPass: os.Getenv("MONGO_DB_ROOT_PASSWORD"),
	}
}

func getEnvRequired(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("CRITICAL ERROR: Environment variable %s is not set.", key)
	}
	return value
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
