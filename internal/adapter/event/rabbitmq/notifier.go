package rabbitmq

import (
	"context"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Publisher struct {
	log *zap.Logger
	ch  *amqp.Channel
}

func NewPublisher(log *zap.Logger, url string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{log: log, ch: ch}, nil
}

// Publish publica en una cola directa (sin exchange) con TTL opcional (ms).
func (p *Publisher) Publish(ctx context.Context, queue string, body []byte, ttlMs int) error {
	// Asegura que la cola exista (durable)
	if _, err := p.ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return err
	}
	pub := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}
	if ttlMs > 0 {
		pub.Expiration = strconv.Itoa(ttlMs)
	}
	if err := p.ch.PublishWithContext(ctx, "", queue, false, false, pub); err != nil {
		p.log.Warn("rabbit publish failed", zap.Error(err), zap.String("queue", queue))
		return err
	}
	p.log.Debug("rabbit message published", zap.String("queue", queue))
	return nil
}
