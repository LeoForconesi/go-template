package user_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	duser "github.com/LeonardoForconesi/go-template/internal/domain/user"
)

func TestUserValidation(t *testing.T) {
	t.Run("should create valid user", func(t *testing.T) {
		u := duser.User{
			ID:    uuid.New(),
			Name:  "Leo",
			Email: "leo@email.com",
			Phone: "123",
		}
		require.NotEmpty(t, u.ID)
		require.Equal(t, "Leo", u.Name)
	})

	t.Run("should fail if name missing", func(t *testing.T) {
		u := duser.User{Email: "a@b.com"}
		err := u.Validate()
		require.ErrorIs(t, err, duser.ErrNameRequired)
	})
}
