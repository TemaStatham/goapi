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

func TestAddCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	categoryRepo := NewCategoryRepository(sqlxDB, logger)

	testName := "TestCategory"
	expectedID := int64(1)

	mock.ExpectQuery("^INSERT INTO categoryies (.+) RETURNING id$").
		WithArgs(testName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	categoryID, err := categoryRepo.AddCategory(context.Background(), testName)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, categoryID)
}

func TestDeleteCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	categoryRepo := NewCategoryRepository(sqlxDB, logger)

	testID := int64(1)

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM categoryies WHERE id = \\$1$").
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("^DELETE FROM product_category WHERE category_id = \\$1$").
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = categoryRepo.DeleteCategory(context.Background(), testID)
	assert.NoError(t, err)
}

func TestUpdateCategoryName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	categoryRepo := NewCategoryRepository(sqlxDB, logger)

	testID := int64(1)
	testName := "TestCategory"

	mock.ExpectExec("^UPDATE categoryies SET name = \\$1 WHERE id = \\$2$").
		WithArgs(testName, testID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	categoryID, err := categoryRepo.UpdateCategoryName(context.Background(), testID, testName)
	assert.NoError(t, err)
	assert.Equal(t, testID, categoryID)
}

func TestGetAllCategoryies(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	categoryRepo := NewCategoryRepository(sqlxDB, logger)

	expectedCategories := []model.Category{
		{ID: 1, Name: "Category1"},
		{ID: 2, Name: "Category2"},
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Category1").
		AddRow(2, "Category2")

	mock.ExpectQuery("^SELECT \\* FROM categoryies$").
		WillReturnRows(rows)

	categories, err := categoryRepo.GetAllCategoryies(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedCategories, categories)
}
