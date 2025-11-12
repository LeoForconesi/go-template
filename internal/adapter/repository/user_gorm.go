package repository

import (
	"context"
	"errors"
	"time"

	"github.com/LeonardoForconesi/go-template/internal/domain"
	duser "github.com/LeonardoForconesi/go-template/internal/domain/user"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

// gormUser es el mapeo a tabla.
type gormUser struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"type:text;not null"`
	Email     string    `gorm:"type:text;uniqueIndex;not null"`
	Phone     string    `gorm:"type:text;not null;default:''"`
	CreatedAt time.Time `gorm:"type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamptz;autoUpdateTime"`
}

// toDomain convierte a entidad de dominio.
func (gu gormUser) toDomain() duser.User {
	return duser.User{
		ID:        gu.ID,
		Name:      gu.Name,
		Email:     gu.Email,
		Phone:     gu.Phone,
		CreatedAt: gu.CreatedAt.UTC(),
		UpdatedAt: gu.UpdatedAt.UTC(),
	}
}

func fromDomain(u *duser.User) gormUser {
	return gormUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt.UTC(),
		UpdatedAt: u.UpdatedAt.UTC(),
	}
}

type UserGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) *UserGormRepository {
	return &UserGormRepository{db: db}
}

func (r *UserGormRepository) Create(ctx context.Context, u *duser.User) error {
	gu := fromDomain(u)
	if err := r.db.WithContext(ctx).Create(&gu).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UserGormRepository) GetByID(ctx context.Context, id uuid.UUID) (*duser.User, error) {
	var gu gormUser
	err := r.db.WithContext(ctx).First(&gu, "id = ?", id).Error
	if err != nil {
		if MapNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	u := gu.toDomain()
	return &u, nil
}

func (r *UserGormRepository) GetByEmail(ctx context.Context, email string) (*duser.User, error) {
	var gu gormUser
	err := r.db.WithContext(ctx).First(&gu, "email = ?", email).Error
	if err != nil {
		if MapNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	u := gu.toDomain()
	return &u, nil
}

func (r *UserGormRepository) List(ctx context.Context, page, size int) ([]duser.User, int64, error) {
	var (
		gu    []gormUser
		total int64
	)

	q := r.db.WithContext(ctx).Model(&gormUser{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := q.Order("created_at DESC").Limit(size).Offset(offset).Find(&gu).Error; err != nil {
		return nil, 0, err
	}

	out := make([]duser.User, 0, len(gu))
	for _, row := range gu {
		out = append(out, row.toDomain())
	}

	return out, total, nil
}

func (r *UserGormRepository) Update(ctx context.Context, u *duser.User) error {
	tx := r.db.WithContext(ctx).Model(&gormUser{}).Where("id = ?", u.ID).Updates(map[string]any{
		"name":  u.Name,
		"phone": u.Phone,
	})
	if tx.Error != nil {
		if isUniqueViolation(tx.Error) {
			return domain.ErrAlreadyExists
		}
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *UserGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx := r.db.WithContext(ctx).Delete(&gormUser{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	// idempotente: no error si RowsAffected == 0
	return nil
}

// ---------------------- helpers ----------------------

// Detecta violación de unique en Postgres (código 23505).
func isUniqueViolation(err error) bool {
	var pqErr interface{ SQLState() string }
	if errors.As(err, &pqErr) {
		return pqErr.SQLState() == "23505"
	}
	// gorm/pgx a veces envuelve en err.Error(); fallback simple:
	return false
}

func (gormUser) TableName() string {
	return "users"
}
