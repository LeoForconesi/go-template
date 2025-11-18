package useruc

/*
Delete: idempotente (el repo no debería fallar si no existe).
*/

import (
	"context"

	"github.com/LeonardoForconesi/go-template/internal/domain"
	duser "github.com/LeonardoForconesi/go-template/internal/domain/user"
	"github.com/google/uuid"
)

type DeleteInput struct {
	ID uuid.UUID
}

type Deleter struct {
	Repo duser.Repository
}

func (uc *Deleter) Execute(ctx context.Context, in DeleteInput) error {
	if in.ID == uuid.Nil {
		return domain.ErrInvalidArgument
	}
	// Repo.Delete debe ser idempotente: no falla si ya no existe.
	return uc.Repo.Delete(ctx, in.ID)
}
