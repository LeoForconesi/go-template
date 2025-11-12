package userService

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	duser "github.com/LeonardoForconesi/go-template/internal/domain/user"
)

// EventProducer define lo que necesitamos del productor Kafka genérico.
type EventProducer interface {
	Produce(ctx context.Context, topic string, key, value []byte) error
}

type UserEventService struct {
	log      *zap.Logger
	producer EventProducer
	topic    string
}

func NewUserEventService(log *zap.Logger, producer EventProducer, topic string) *UserEventService {
	return &UserEventService{log: log, producer: producer, topic: topic}
}

// Satisface duser.EventPublisher
func (s *UserEventService) PublishUserCreated(ctx context.Context, ev duser.CreatedEvent) error {
	payload, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	// Usamos user-id como key (particionado estable)
	return s.producer.Produce(ctx, s.topic, []byte(ev.ID), payload)
}
