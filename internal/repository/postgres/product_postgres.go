package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"goapi/internal/model"
	"log/slog"
)

const (
	productsTable = "products"
	ErrProductID  = 0
)

var (
	ErrStartTransaction = errors.New("error starting transaction")
	ErrEndTransaction   = errors.New("error committing transaction")
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

func (p *ProductRepository) AddProduct(
	ctx context.Context,
	name string,
	categoryies []string,
) (int64, error) {
	const op = "postgres.AddProduct"

	log := p.log.With(
		slog.String("op", op),
		slog.String("name", name),
	)

	log.Info("save product in db")

	tx, err := p.db.Beginx()
	if err != nil {
		log.Error(ErrStartTransaction.Error())
		return 0, fmt.Errorf("%s %w", op, ErrStartTransaction)
	}
	defer tx.Rollback()

	var categoryIDs []int64
	query := fmt.Sprintf(
		"SELECT id FROM %s WHERE name = $1",
		categoryTable,
	)
	for _, categoryName := range categoryies {
		var categoryID int64
		err := tx.Get(&categoryID, query, categoryName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Error("category %s does not exist in the database", categoryName)
				return ErrProductID, fmt.Errorf("%s category %s does not exist in the database", op, categoryName)
			}

			log.Error("error checking category in the database")
			return ErrProductID, fmt.Errorf("%s %w", op, err)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	query = fmt.Sprintf(
		"INSERT INTO %s (name) VALUES ($1) RETURNING id",
		productsTable,
	)
	var productID int64
	err = tx.Get(&productID, query, name)
	if err != nil {
		log.Error("error inserting product into the database")
		return ErrProductID, fmt.Errorf("%s %w", op, err)
	}

	query = fmt.Sprintf(
		"INSERT INTO %s (product_id, category_id) VALUES ($1, $2)",
		productCategoryTable,
	)
	for _, categoryID := range categoryIDs {
		_, err := tx.Exec(query, productID, categoryID)
		if err != nil {
			log.Error("error inserting product-category relationship into the database")
			return ErrProductID, fmt.Errorf("%s %w", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error(ErrEndTransaction.Error())
		return ErrProductID, fmt.Errorf("%s %w", op, ErrEndTransaction)
	}

	log.Info("product is saved in db successfully")

	return productID, nil
}

func (p *ProductRepository) DeleteProduct(
	ctx context.Context,
	id int64,
) error {
	const op = "postgres.DeleteProduct"

	log := p.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("removing a product from the database")

	tx, err := p.db.Beginx()
	if err != nil {
		log.Error(ErrStartTransaction.Error())
		return fmt.Errorf("%s %w", op, ErrStartTransaction)
	}
	defer tx.Rollback()

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id = $1",
		productsTable,
	)
	_, err = tx.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		log.Error("error deleting a product from the database")
		return fmt.Errorf("%s %w", op, err)
	}

	query = fmt.Sprintf(
		"DELETE FROM %s WHERE product_id = $1",
		productCategoryTable,
	)
	_, err = tx.Exec(query, id)
	if err != nil {
		log.Error("error deleting product-category links from the database")
		return fmt.Errorf("%s %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		log.Error(ErrEndTransaction.Error())
		return fmt.Errorf("%s %w", op, ErrEndTransaction)
	}

	log.Info("the product was successfully removed from the database")

	return nil
}

func (p *ProductRepository) UpdateProductName(
	ctx context.Context,
	id int64,
	name string,
) (int64, error) {
	const op = "postgres.UpdateProductName"

	log := p.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
		slog.String("name", name),
	)

	log.Info("updating the product name in the database")

	query := fmt.Sprintf(
		"UPDATE %s SET name = $1 WHERE id = $2 RETURNING id",
		productsTable,
	)
	result, err := p.db.Exec(query, name, id)
	if err != nil {
		log.Error("error updating product name in database\n")
		return ErrProductID, fmt.Errorf("%s %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("error getting number of affected rows")
		return ErrProductID, fmt.Errorf("%s %w", op, err)
	}

	if rowsAffected == 0 {
		log.Warn("no product found with the specified ID\n")
		return ErrProductID, fmt.Errorf("%s no product found with the specified ID\n", op)
	}

	log.Info("product name was successfully updated in the database")

	return id, nil
}

func (p *ProductRepository) UpdateProductCategoryies(
	ctx context.Context,
	id int64,
	categoryies []model.Category,
) (int64, error) {
	const op = "postgres.UpdateProductCategoryies"

	log := p.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("updating product categories in the database")

	tx, err := p.db.Beginx()
	if err != nil {
		log.Error(ErrStartTransaction.Error())
		return ErrProductID, fmt.Errorf("%s %w", op, ErrStartTransaction)
	}
	defer tx.Rollback()

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE product_id = $1",
		productCategoryTable,
	)
	_, err = tx.Exec("DELETE FROM product_category WHERE product_id = $1", id)
	if err != nil {
		log.Error("error deleting product-category links from the database\n")
		return ErrProductID, fmt.Errorf("%s %w", op, err)
	}

	query = fmt.Sprintf(
		"INSERT INTO %s (product_id, category_id) VALUES ($1, $2)",
		productCategoryTable,
	)
	for _, category := range categoryies {
		_, err := tx.Exec(query, id, category.ID)
		if err != nil {
			log.Error("error adding a new product-category link to the database\n")
			return ErrProductID, fmt.Errorf("%s %w", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error(ErrEndTransaction.Error())
		return ErrProductID, fmt.Errorf("%s %w", op, ErrEndTransaction)
	}

	log.Info("product categories have been successfully updated in the database\n")

	return id, nil
}

func (p *ProductRepository) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	const op = "postgres.GetAllProducts"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("getting all products from the database\n")

	var products []model.Product
	query := fmt.Sprintf(
		"SELECT * FROM $s",
		productsTable,
	)
	err := p.db.Select(&products, query)
	if err != nil {
		log.Error("error getting products from database\n")
		return nil, fmt.Errorf("%s %w", op, err)
	}

	log.Info("products successfully retrieved from database")

	return products, nil
}

func (p *ProductRepository) GetCategoryProducts(ctx context.Context, category string) ([]model.Product, error) {
	const op = "postgres.GetCategoryProducts"

	log := p.log.With(
		slog.String("op", op),
		slog.String("category", category),
	)

	log.Info("getting products by category from db")

	var products []model.Product
	query := fmt.Sprintf(
		" SELECT p.id, p.name\n"+
			"FROM %s p\n"+
			"INNER JOIN product_category pc ON p.id = pc.product_id\n"+
			"INNER JOIN categories c ON pc.category_id = c.id\n"+
			"WHERE c.name = $1",
		productsTable,
	)

	if err := p.db.SelectContext(ctx, &products, query, category); err != nil {
		log.Error("failed to get products by category from db")
		return nil, fmt.Errorf("%s %w", op, err)
	}

	log.Info("products by category retrieved from db")

	return products, nil
}
