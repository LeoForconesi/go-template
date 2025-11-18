package useruc

import (
	"context"

	"github.com/LeonardoForconesi/go-template/internal/domain"
	duser "github.com/LeonardoForconesi/go-template/internal/domain/user"
	"github.com/google/uuid"
)

type GetInput struct {
	ID uuid.UUID
}

type GetOutput struct {
	ID    uuid.UUID
	Name  string
	Email string
	Phone string
}

type Getter struct {
	Repo duser.Repository
}

func (uc *Getter) Execute(ctx context.Context, in GetInput) (GetOutput, error) {
	if in.ID == uuid.Nil {
		return GetOutput{}, domain.ErrInvalidArgument
	}

	u, err := uc.Repo.GetByID(ctx, in.ID)
	if err != nil {
		return GetOutput{}, err
	}
	return GetOutput{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}, nil
}
