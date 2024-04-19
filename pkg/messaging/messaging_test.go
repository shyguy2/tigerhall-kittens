package messaging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockMessageBroker is a mock implementation of the MessageBroker interface.
type mockMessageBroker struct {
	publishedMessage []byte
	consumeFunc      func([]byte) error
}

func (m *mockMessageBroker) PublishMessage(message []byte) error {
	m.publishedMessage = message
	return nil
}

func (m *mockMessageBroker) ConsumeMessages(processMessage func([]byte) error) {
	m.consumeFunc = processMessage
}

func TestMessageBroker_PublishMessage(t *testing.T) {
	// Create the mock message broker
	mockBroker := &mockMessageBroker{}

	// Publish a test message
	message := []byte("Test Message")
	err := mockBroker.PublishMessage(message)
	assert.NoError(t, err)

	// Assert that the message was published correctly
	assert.Equal(t, message, mockBroker.publishedMessage)
}
