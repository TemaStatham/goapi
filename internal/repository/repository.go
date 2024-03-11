package repository

import "errors"

var (
	ErrUserExist    = errors.New("user already exist")
	ErrUserNotFound = errors.New("user not found")
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}
