package outbox

type MessageBroker interface {
	Send(message Message) error
}
