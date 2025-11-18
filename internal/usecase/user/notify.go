package useruc

/*
Notify: verifica existencia de usuario y publica en RabbitMQ (TTL lo maneja el adapter).
*/

import (
	"context"

	"github.com/LeonardoForconesi/go-template/internal/domain"
	duser "github.com/LeonardoForconesi/go-template/internal/domain/user"
	"github.com/google/uuid"
)

type NotifyInput struct {
	UserID  uuid.UUID
	Message string
}

type Notifier struct {
	Repo     duser.Repository
	Notifier duser.Notifier // RabbitMQ adapter
}

func (uc *Notifier) Execute(ctx context.Context, in NotifyInput) error {
	if in.UserID == uuid.Nil || in.Message == "" {
		return domain.ErrInvalidArgument
	}

	// Aseguramos que el usuario exista (puede ser útil para validaciones).
	if _, err := uc.Repo.GetByID(ctx, in.UserID); err != nil {
		return err
	}
	return uc.Notifier.NotifyUser(ctx, duser.NotifyMessage{
		UserID:  in.UserID.String(),
		Message: in.Message,
	})
}
