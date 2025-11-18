package user

import (
	"context"
	"time"
)

// Evento de dominio a publicar en Kafka.
type CreatedEvent struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	OccurredAt time.Time `json:"occurred_at"`
}

// Puerto de salida: publicador de eventos (Kafka).
type EventPublisher interface {
	PublishUserCreated(ctx context.Context, ev CreatedEvent) error
}

// Puerto de salida: notificador a cola (RabbitMQ).
type Notifier interface {
	NotifyUser(ctx context.Context, msg NotifyMessage) error
}

// Mensaje para RabbitMQ.
type NotifyMessage struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	// TTL lo maneja el adapter; aquí no lo fijamos para mantener el dominio agnóstico.
}
