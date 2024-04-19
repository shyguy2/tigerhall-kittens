package main

import (
	"log"

	conf "github.com/tigerhall-kittens/config"
	inits "github.com/tigerhall-kittens/pkg"
	"github.com/tigerhall-kittens/pkg/auth"
	"github.com/tigerhall-kittens/pkg/server"
)

func main() {
	// Read the configuration from server.yml
	config, err := conf.ReadConfig("config/local/server.yml")
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	// Initialize the service
	service, err := inits.InitializeService(config)
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
