package mock

import "github.com/chrishrb/ezr2mqtt/transport"

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (c *MockClient) Connect() (*transport.Message, error) {
	return nil, nil
}

func (c *MockClient) Send(msg *transport.Message) error {
	return nil
}
