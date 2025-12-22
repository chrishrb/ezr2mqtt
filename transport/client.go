package transport

type Client interface {
	Connect() (*Message, error)
	Send(message *Message) error
}
