package userService

import (
	"context"
	"encoding/json"

	duser "github.com/LeonardoForconesi/go-template/internal/domain/user"
	"go.uber.org/zap"
)

// QueuePublisher es el publisher genérico de Rabbit.
type QueuePublisher interface {
	Publish(ctx context.Context, queue string, body []byte, ttlMs int) error
}

type UserNotifyService struct {
	log       *zap.Logger
	publisher QueuePublisher
	queue     string
	ttlMs     int
}

func NewUserNotifyService(log *zap.Logger, publisher QueuePublisher, queue string, ttlMs int) *UserNotifyService {
	return &UserNotifyService{log: log, publisher: publisher, queue: queue, ttlMs: ttlMs}
}

// Satisface duser.Notifier
func (s *UserNotifyService) NotifyUser(ctx context.Context, msg duser.NotifyMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return s.publisher.Publish(ctx, s.queue, body, s.ttlMs)
}
