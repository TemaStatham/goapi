package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"goapi/internal/model"
	"goapi/internal/repository"
	"log/slog"
)

const (
	ErrCategoryID        = 0
	categoryTable        = "categoryies"
	productCategoryTable = "product_category"
)

type CategoryRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewCategoryRepository(db *sqlx.DB, l *slog.Logger) *CategoryRepository {
	return &CategoryRepository{
		db:  db,
		log: l,
	}
}

func (c *CategoryRepository) AddCategory(ctx context.Context, name string) (int64, error) {
	const op = "postgres.AddCategory"

	log := c.log.With(
		slog.String("op", op),
		slog.String("name", name),
	)

	log.Info("save product in db")

	var id int64

	query := fmt.Sprintf(
		"INSERT INTO %s (name) VALUES ($1) RETURNING id",
		categoryTable,
	)

	row := c.db.QueryRow(query, name)
	if err := row.Scan(&id); err != nil {
		log.Error("error insert category in db")
		return id, repository.ErrCategoryExist
	}

	log.Info("product saved in db successfully")

	return id, nil
}

func (c *CategoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	const op = "postgres.DeleteCategory"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("removing a category from the database")

	tx, err := c.db.Beginx()
	if err != nil {
		log.Error(ErrStartTransaction.Error())
		return fmt.Errorf("%s %w", op, ErrStartTransaction)
	}
	defer tx.Rollback()

	query := fmt.Sprintf(
		"DELETE FROM categories WHERE id = $1",
		categoryTable,
	)
	_, err = tx.Exec(query, id)
	if err != nil {
		log.Error("error deleting category from database")
		return fmt.Errorf("%s %w", op, repository.ErrCategoryDelete)
	}

	query = fmt.Sprintf(
		"DELETE FROM %s WHERE category_id = $1",
		productCategoryTable,
	)
	_, err = tx.Exec(query, id)
	if err != nil {
		log.Error("error deleting related product categories from the database\n")
		return fmt.Errorf("%s %w", op, repository.ErrDeleteProductCategory)
	}

	err = tx.Commit()
	if err != nil {
		log.Error("transaction commit error\n")
		return fmt.Errorf("%s %w", op, ErrEndTransaction)
	}

	log.Info("category successfully deleted from the database")

	return nil
}

func (c *CategoryRepository) UpdateCategoryName(ctx context.Context, id int64, name string) (int64, error) {
	const op = "postgres.UpdateCategoryName"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
		slog.String("name", name),
	)

	log.Info("updating the category name in the database")

	query := fmt.Sprintf(
		"UPDATE %s SET name = $1 WHERE id = $2",
		categoryTable,
	)
	result, err := c.db.Exec(query, name, id)
	if err != nil {
		log.Error("error updating category name in database")
		return ErrCategoryID, fmt.Errorf("%s %w", op, repository.ErrUpdateCategory)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("error getting number of affected rows\n")
		return ErrCategoryID, fmt.Errorf("%s %w", op, err)
	}

	if rowsAffected == 0 {
		log.Warn("сategory with specified ID not found\n")
		return ErrCategoryID, fmt.Errorf("сategory with specified ID not found\n")
	}

	log.Info("category name successfully updated in database\n")

	return id, nil
}

func (c *CategoryRepository) GetAllCategoryies(ctx context.Context) ([]model.Category, error) {
	const op = "postgres.GetAllCategories"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("getting all categories from the database")

	var categories []model.Category

	query := fmt.Sprintf(
		"SELECT * FROM %s",
		categoryTable,
	)
	err := c.db.Select(&categories, query)
	if err != nil {
		log.Error("error getting categories from database\n")
		return nil, fmt.Errorf("%s %w", op, repository.ErrAllCategoryies)
	}

	log.Info("categories successfully retrieved from database")

	return categories, nil
}
