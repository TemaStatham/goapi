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
	productsTable = "product"
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
	const op = "productRepository.AddProduct"

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
	for _, categoryName := range categoryies {
		var categoryID int64
		err := tx.Get(&categoryID, "SELECT id FROM categories WHERE name = $1", categoryName)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Error("category %s does not exist in the database", categoryName)
				return ErrProductID, fmt.Errorf("%s category %s does not exist in the database", op, categoryName)
			}

			log.Error("error checking category in the database")
			return ErrProductID, fmt.Errorf("%s %w", op, err)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	var productID int64
	err = tx.Get(&productID, "INSERT INTO products (name) VALUES ($1) RETURNING id", name)
	if err != nil {
		log.Error("error inserting product into the database")
		return ErrProductID, fmt.Errorf("%s %w", op, err)
	}

	for _, categoryID := range categoryIDs {
		_, err := tx.Exec("INSERT INTO product_category (product_id, category_id) VALUES ($1, $2)", productID, categoryID)
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
	const op = "productRepository.DeleteProduct"

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

	_, err = tx.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		log.Error("error deleting a product from the database")
		return fmt.Errorf("%s %w", op, err)
	}

	_, err = tx.Exec("DELETE FROM product_category WHERE product_id = $1", id)
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
	const op = "productRepository.UpdateProductName"

	log := p.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
		slog.String("name", name),
	)

	log.Info("updating the product name in the database")

	result, err := p.db.Exec("UPDATE products SET name = $1 WHERE id = $2 RETURNING id", name, id)
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
	const op = "productRepository.UpdateProductCategoryies"

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

	_, err = tx.Exec("DELETE FROM product_category WHERE product_id = $1", id)
	if err != nil {
		log.Error("error deleting product-category links from the database\n")
		return ErrProductID, fmt.Errorf("%s %w", op, err)
	}

	for _, category := range categoryies {
		_, err := tx.Exec("INSERT INTO product_category (product_id, category_id) VALUES ($1, $2)", id, category.ID)
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
	const op = "productRepository.GetAllProducts"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("getting all products from the database\n")

	var products []model.Product
	err := p.db.Select(&products, "SELECT * FROM products")
	if err != nil {
		log.Error("error getting products from database\n")
		return nil, fmt.Errorf("%s %w", op, err)
	}

	log.Info("products successfully retrieved from database")

	return products, nil
}
