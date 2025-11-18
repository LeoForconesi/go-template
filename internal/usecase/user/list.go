package useruc

/*
Get/List: validaciones simples (UUID, paginado).
*/

import (
	"context"

	"github.com/LeonardoForconesi/go-template/internal/domain"
	duser "github.com/LeonardoForconesi/go-template/internal/domain/user"
)

type ListInput struct {
	Page int
	Size int
}

type ListItem struct {
	ID    string
	Name  string
	Email string
	Phone string
}

type ListOutput struct {
	Items []ListItem
	Total int64
	Page  int
	Size  int
}

type Lister struct {
	Repo duser.Repository
}

func (uc *Lister) Execute(ctx context.Context, in ListInput) (ListOutput, error) {
	page := in.Page
	size := in.Size

	if page < 1 || size < 1 {
		return ListOutput{}, domain.ErrInvalidArgument
	}

	users, total, err := uc.Repo.List(ctx, page, size)
	if err != nil {
		return ListOutput{}, err
	}

	out := ListOutput{
		Items: make([]ListItem, 0, len(users)),
		Total: total,
		Page:  page,
		Size:  size,
	}
	for _, u := range users {
		out.Items = append(out.Items, ListItem{
			ID:    u.ID.String(),
			Name:  u.Name,
			Email: u.Email,
			Phone: u.Phone,
		})
	}
	return out, nil
}
