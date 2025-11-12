package repository

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres abre una conexión GORM sobre Postgres (acepta DSN completo).
func NewPostgres(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// MapNotFound traduce el error de GORM a un booleano “not found”.
func MapNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
