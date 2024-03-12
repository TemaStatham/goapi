package postgres

import (
	"github.com/jmoiron/sqlx"
	"log/slog"
)

const (
	productsTable = "product"
)

type ProductRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewProductRepository(db *sqlx.DB, l *slog.Logger) *ProductRepository {
	return &ProductRepository{
		db:  db,
		log: l,
	}
}
