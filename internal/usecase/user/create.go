package useruc

/*
Create: valida dominio, chequea duplicado por email,
crea y publica users.created de forma best-effort (si querés que falle si no se puede publicar, lo cambiamos).
*/

import (
	"context"
	"time"

	"github.com/LeonardoForconesi/go-template/internal/domain"

	duser "github.com/LeonardoForconesi/go-template/internal/domain/user"
	"github.com/google/uuid"
)

type CreateInput struct {
	Name  string
	Email string
	Phone string
}

type CreateOutput struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Creator struct {
	Repo      duser.Repository
	Publisher duser.EventPublisher
	// (Opcional) Tx duser.Tx
}

func (uc *Creator) Execute(ctx context.Context, in CreateInput) (CreateOutput, error) {
	u := &duser.User{
		ID:        uuid.New(),
		Name:      in.Name,
		Email:     in.Email,
		Phone:     in.Phone,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := u.Validate(); err != nil {
		return CreateOutput{}, domain.ErrInvalidArgument
	}

	// Duplicado por email
	if existing, err := uc.Repo.GetByEmail(ctx, u.Email); err == nil && existing != nil {
		return CreateOutput{}, domain.ErrAlreadyExists
	}

	// Persistencia
	if err := uc.Repo.Create(ctx, u); err != nil {
		// El repo puede mapear unique violation a ErrAlreadyExists
		return CreateOutput{}, err
	}

	// Publicar evento de dominio (best-effort; si querés, podés decidir fallar si el publish falla)
	if uc.Publisher != nil {
		_ = uc.Publisher.PublishUserCreated(ctx, duser.CreatedEvent{
			ID:         u.ID.String(),
			Name:       u.Name,
			Email:      u.Email,
			OccurredAt: time.Now().UTC(),
		})
	}

	return CreateOutput{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}
