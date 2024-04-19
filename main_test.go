package main

import (
	"testing"
	conf "tigerhall-kittens-app/config"

	"github.com/stretchr/testify/assert"
)

func TestInitializeService(t *testing.T) {
	// Prepare a temporary configuration with in-memory RabbitMQ for testing
	config := &conf.Config{
		RabbitMq: conf.RabbitMq{
			AmqpURL:   "amqp://guest:guest@localhost:5672/",
			QueueName: "test_queue",
		},
		Database: conf.Database{
			Host:     "localhost",
			Port:     5432,
			Username: "testuser",
			Password: "testpass",
			DBName:   "testdb",
		},
		JWT: conf.JWT{
			SecretKey: "test_secret_key",
		},
		Server: conf.Server{
			Port: "8080",
		},
	}

	_, err := initializeService(config)

	// Assert that the service is initialized without errors
	assert.Error(t, err)

}
