package pkg

import (
	"log"

	conf "github.com/tigerhall-kittens/config"
	"github.com/tigerhall-kittens/pkg/auth"
	"github.com/tigerhall-kittens/pkg/messaging"
	"github.com/tigerhall-kittens/pkg/repository"
	"github.com/tigerhall-kittens/pkg/server"
	"github.com/tigerhall-kittens/pkg/service"
)

func InitializeService(config *conf.Config) (service.TigerService, error) {
	// Initialize the database connection
	dbConnectionString := conf.BuildDBConnectionString(config.Database)

	store, err := repository.NewPostgresRepository(dbConnectionString)
	if err != nil {
		return nil, err
	}

	// Initialize the RabbitMQ message broker
	messageBroker, err := messaging.NewMessageBroker(config.RabbitMq.AmqpURL, config.RabbitMq.QueueName)
	if err != nil {
		return nil, err
	}

	// Start the message consumer in a separate Goroutine
	go messageBroker.ConsumeMessages(messaging.ProcessMessage)

	// Initialize the service
	service := service.NewTigerService(store, messageBroker)

	return service, nil
}

func main() {
	// Read the configuration from server.yml
	config, err := conf.ReadConfig("config/local/server.yml")
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	// Initialize the service
	service, err := InitializeService(config)
	if err != nil {
		log.Fatalf("Failed to initialize the service: %v", err)
	}

	// Initialize the server
	srv := server.NewServer()

	// Set up the routes and handlers
	srv.SetupRoutes(service, auth.NewAuth(config.JWT.SecretKey))

	// Start the server
	err = srv.Start(config.Server.Port)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
