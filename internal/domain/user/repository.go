package user

import (
	"context"

	"github.com/google/uuid"
)

// Puerto de salida: persistencia (será implementado con GORM).
type Repository interface {
	// Create persiste un nuevo usuario. Debe fallar con ErrDuplicateEmail si email ya existe.
	Create(ctx context.Context, u *User) error

	// GetByID retorna un usuario o ErrNotFound si no existe.
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// GetByEmail retorna un usuario por email (útil para idempotencia/validaciones).
	GetByEmail(ctx context.Context, email string) (*User, error)

	// List pagina resultados (page >=1, size >0). Retorna slice y total opcional.
	List(ctx context.Context, page, size int) ([]User, int64, error)

	// Update actualiza campos mutables (name, phone).
	Update(ctx context.Context, u *User) error

	// Delete elimina por ID (debe ser idempotente: no error si ya no existe).
	Delete(ctx context.Context, id uuid.UUID) error
}

// (Opcional para más adelante) Transacciones si las necesitamos.
type Tx interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
