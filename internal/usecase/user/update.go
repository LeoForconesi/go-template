package useruc

/*
Update: trae actual, muta campos permitidos, valida y persiste.
*/

import (
	"context"
	"time"

	"github.com/LeonardoForconesi/go-template/internal/domain"
	duser "github.com/LeonardoForconesi/go-template/internal/domain/user"
	"github.com/google/uuid"
)

type UpdateInput struct {
	ID    uuid.UUID
	Name  *string
	Phone *string
}

type UpdateOutput struct {
	ID    uuid.UUID
	Name  string
	Email string
	Phone string
}

type Updater struct {
	Repo duser.Repository
}

func (uc *Updater) Execute(ctx context.Context, in UpdateInput) (UpdateOutput, error) {
	if in.ID == uuid.Nil {
		return UpdateOutput{}, domain.ErrInvalidArgument
	}
	if in.Name == nil && in.Phone == nil {
		return UpdateOutput{}, domain.ErrInvalidArgument
	}

	// Traemos el actual
	u, err := uc.Repo.GetByID(ctx, in.ID)
	if err != nil {
		return UpdateOutput{}, err
	}

	// Mutaciones permitidas
	if in.Name != nil {
		u.Name = *in.Name
	}
	if in.Phone != nil {
		u.Phone = *in.Phone
	}
	u.UpdatedAt = time.Now().UTC()

	// Validación de dominio (p.ej. name requerido)
	if err := u.Validate(); err != nil {
		return UpdateOutput{}, domain.ErrInvalidArgument
	}

	if err := uc.Repo.Update(ctx, u); err != nil {
		return UpdateOutput{}, err
	}

	return UpdateOutput{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}, nil
}
