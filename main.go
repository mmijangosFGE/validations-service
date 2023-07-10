package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mmijangosFGE/validations-service/adapters/api/httpServer"
	mongoDriver "github.com/mmijangosFGE/validations-service/internal/db/drivers/mongo"
	"github.com/mmijangosFGE/validations-service/pkg/constants"
	"github.com/mmijangosFGE/validations-service/pkg/messages"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Load environment variables from .env file
	errEnv := godotenv.Load()
	// Check if the .env file was loaded successfully
	if errEnv != nil {
		log.Fatal(messages.LoadEnvError)
	}
	// Retrieve environment variables
	port := os.Getenv("PORT")
	databaseMongoUrl := os.Getenv("DATABASE_MONGO_URL")
	mongoDB := os.Getenv("MONGO_DB")
	jwtSecret := os.Getenv("JWT_SECRET")
	env := os.Getenv("ENV")

	// Establish a resilient connection to MongoDB
	connection := &mongoDriver.Connection{
		ClientChan: make(chan *mongo.Client, 1),
		Connector:  &mongoDriver.DBConnector{}, // Pass the real MongoDB connector here
		State:      constants.Disconnected,
	}

	// Start a goroutine to monitor the MongoDB connection
	go connection.MonitorConnection(databaseMongoUrl)

	go func() {
		// Receive the connected MongoDB client from the channel
		client := <-connection.ClientChan

		// Add a small delay to let the connection goroutine start running
		time.Sleep(500 * time.Millisecond)

		// Create a new server instance and pass the environment variables to the configuration
		s, errServer := httpServer.NewServer(&httpServer.Config{
			Env:       env,
			JWTSecret: jwtSecret,
			MongoDb:   mongoDB,
			Port:      port,
		})
		// Check if there was an error when creating the server
		if errServer != nil {
			panic(errServer)
		}
		// Start the server
		s.Start(client, s)
	}()

	// Prevent the main function from exiting
	select {}
}
