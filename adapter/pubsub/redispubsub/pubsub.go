package redispubsub

import (
	"context"

	redis2 "github.com/redis/go-redis/v9"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	"github.com/syntaxfa/quick-connect/pkg/pubsub"
)

type PubSub struct {
	re *redis.Adapter
}

func New(re *redis.Adapter) *PubSub {
	return &PubSub{
		re: re,
	}
}

func (p *PubSub) Publish(ctx context.Context, channel string, message []byte) error {
	return p.re.Client().Publish(ctx, channel, message).Err()
}

func (p *PubSub) Subscribe(ctx context.Context, channels ...string) pubsub.Receiver {
	pubSub := p.re.Client().Subscribe(ctx, channels...)

	return newReceiver(pubSub)
}

type receiver struct {
	pubSub *redis2.PubSub
}

func newReceiver(pubSub *redis2.PubSub) *receiver {
	return &receiver{
		pubSub: pubSub,
	}
}

func (r *receiver) ReceiveMessage(ctx context.Context) ([]byte, error) {
	redisMsg, rErr := r.pubSub.ReceiveMessage(ctx)
	if rErr != nil {
		return nil, rErr
	}

	return []byte(redisMsg.Payload), nil
}
