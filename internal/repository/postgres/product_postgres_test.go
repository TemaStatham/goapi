package postgres

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"goapi/internal/model"
)

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	productRepo := NewProductRepository(sqlxDB, logger)

	testID := int64(1)

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM products WHERE id = \\$1$").
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("^DELETE FROM product_category WHERE product_id = \\$1$").
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = productRepo.DeleteProduct(context.Background(), testID)
	assert.NoError(t, err)
}

func TestUpdateProductName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	productRepo := NewProductRepository(sqlxDB, logger)

	testID := int64(1)
	testName := "TestProduct"

	mock.ExpectExec("^UPDATE products SET name = \\$1 WHERE id = \\$2 RETURNING id$").
		WithArgs(testName, testID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	productID, err := productRepo.UpdateProductName(context.Background(), testID, testName)
	assert.NoError(t, err)
	assert.Equal(t, testID, productID)
}

func TestUpdateProductCategoryies(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	productRepo := NewProductRepository(sqlxDB, logger)

	testID := int64(1)
	testCategories := []model.Category{
		{ID: 1, Name: "Category1"},
		{ID: 2, Name: "Category2"},
	}

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM product_category WHERE product_id = \\$1$").
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("^INSERT INTO product_category (.+)").
		WithArgs(testID, testCategories[0].ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("^INSERT INTO product_category (.+)").
		WithArgs(testID, testCategories[1].ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	productID, err := productRepo.UpdateProductCategoryies(context.Background(), testID, testCategories)
	assert.NoError(t, err)
	assert.Equal(t, testID, productID)
}

func TestGetAllProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	productRepo := NewProductRepository(sqlxDB, logger)

	expectedProducts := []model.Product{
		{ID: 1, Name: "Product1"},
		{ID: 2, Name: "Product2"},
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Product1").
		AddRow(2, "Product2")

	mock.ExpectQuery("^SELECT \\* FROM products$").
		WillReturnRows(rows)

	products, err := productRepo.GetAllProducts(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
}

func TestGetCategoryProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	productRepo := NewProductRepository(sqlxDB, logger)

	testCategory := "TestCategory"

	expectedProducts := []model.Product{
		{ID: 1, Name: "Product1"},
		{ID: 2, Name: "Product2"},
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Product1").
		AddRow(2, "Product2")

	mock.ExpectQuery("^SELECT (.+) FROM products p (.+)").
		WithArgs(testCategory).
		WillReturnRows(rows)

	products, err := productRepo.GetCategoryProducts(context.Background(), testCategory)
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
}
