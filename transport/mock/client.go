package mock

import (
	"github.com/chrishrb/ezr2mqtt/transport"
)

type MockClient struct {
	currentMessage *transport.Message
}

func NewMockClient() *MockClient {
	return &MockClient{
		currentMessage: createMockMessage(),
	}
}

func (c *MockClient) Connect() (*transport.Message, error) {
	return nil, nil
}

func (c *MockClient) Send(msg *transport.Message) error {
	// TODO: implement
	return nil
}
