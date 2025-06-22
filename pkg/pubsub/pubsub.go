package pubsub

import "context"

type Publisher interface {
	Publish(ctx context.Context, channel string, message []byte) error
}

type Receiver interface {
	ReceiveMessage(ctx context.Context) (message []byte, err error)
}

type Subscriber interface {
	Subscribe(ctx context.Context, channels ...string) Receiver
}

type PubSub interface {
	Publisher
	Subscriber
}
