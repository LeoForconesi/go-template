package user

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrNameRequired  = errors.New("user: name required")
	ErrEmailRequired = errors.New("user: email required")
)

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrNameRequired
	}
	if u.Email == "" {
		return ErrEmailRequired
	}

	return nil
}
