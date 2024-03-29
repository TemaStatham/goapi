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
	usersTable = "users"
)

type AuthRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewAuthPostgres(db *sqlx.DB, l *slog.Logger) *AuthRepository {
	return &AuthRepository{
		db:  db,
		log: l,
	}
}

func (a *AuthRepository) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "AuthRepository.SaveUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("save user in db")

	var id int64

	query := fmt.Sprintf(
		"INSERT INTO %s (email, passHash) VALUES ($1, $2) RETURNING id",
		usersTable,
	)

	row := a.db.QueryRow(query, email, passHash)
	if err := row.Scan(&id); err != nil {
		log.Error("error insert user in db")
		return id, repository.ErrUserExist
	}

	log.Info("user is saved in db successfully")

	return id, nil
}

func (a *AuthRepository) User(ctx context.Context, email string) (model.User, error) {
	const op = "AuthRepository.User"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("get user in db")

	var user model.User

	query := fmt.Sprintf(
		"SELECT id, email, passHash FROM %s WHERE email=$1",
		usersTable,
	)

	err := a.db.Get(&user, query, email)
	if err != nil {
		log.Error("error get user from db")
		return user, repository.ErrUserNotFound
	}

	log.Info("user get from db successfully")

	return user, nil
}
