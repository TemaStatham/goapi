package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"goapi/internal/model"
	"goapi/internal/repository"
)

func TestSaveUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	authRepo := NewAuthPostgres(sqlxDB, logger)

	testEmail := "test@example.com"
	testPassHash := []byte("password_hash")
	expectedID := int64(1)

	mock.ExpectQuery("^INSERT INTO users (.+) RETURNING id$").
		WithArgs(testEmail, testPassHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	userID, err := authRepo.SaveUser(context.Background(), testEmail, testPassHash)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, userID)
}

func TestSaveUserError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	authRepo := NewAuthPostgres(sqlxDB, logger)

	testEmail := "test@example.com"
	testPassHash := []byte("password_hash")

	mock.ExpectQuery("^INSERT INTO users (.+) RETURNING id$").
		WithArgs(testEmail, testPassHash).
		WillReturnError(errors.New("some error"))

	userID, err := authRepo.SaveUser(context.Background(), testEmail, testPassHash)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserExist, err)
	assert.Equal(t, int64(0), userID)
}

func TestUserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	authRepo := NewAuthPostgres(sqlxDB, logger)

	testEmail := "test@example.com"

	mock.ExpectQuery("^SELECT id, email, passHash FROM users WHERE email=\\$1$").
		WithArgs(testEmail).
		WillReturnError(sql.ErrNoRows)

	user, err := authRepo.User(context.Background(), testEmail)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)
	assert.Equal(t, model.User{}, user)
}
