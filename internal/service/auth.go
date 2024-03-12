package service

import (
	"context"
	"errors"
	"fmt"
	"goapi/internal/lib/jwt"
	"goapi/internal/model"
	"goapi/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

const (
	ErrUserID = -1
)

var (
	ErrorInvalidCredentials = errors.New("invalid credential")
	ErrorUserExist          = errors.New("user already exist")
	ErrPasswordIsEmpty      = errors.New("password is empty")
	ErrEmailIsEmpty         = errors.New("email is empty")
	ErrFailedToSaveUser     = errors.New("failed to save user")
)

type AuthService struct {
	usrSaver    UserSaver
	usrProvider UserProvider
	log         *slog.Logger
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (model.User, error)
}

func NewAuthService(us UserSaver, up UserProvider, l *slog.Logger, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		usrSaver:    us,
		usrProvider: up,
		log:         l,
		tokenTTL:    tokenTTL,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	const op = "postgres.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("logging user")

	if err := s.validate(email, password); err != nil {
		log.Error("data is invalid ", err)
		return "", fmt.Errorf("data is invalid: %w", err)
	}

	user, err := s.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.log.Warn("user not found", err)
			return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		s.log.Warn("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Warn("pass", user.PassHash, password)
		s.log.Warn("invalid credential", err)
		return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
	}

	log.Info("user is login successfully")

	token, err := jwt.NewToken(user, s.tokenTTL)
	if err != nil {
		s.log.Warn("failed to generate token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *AuthService) Register(ctx context.Context, email, password string) (userID int64, err error) {
	const op = "postgres.RegisterNewUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	if err := s.validate(email, password); err != nil {
		log.Error("data is invalid ", err)
		return ErrUserID, fmt.Errorf("data is invalid: %w", err)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, repository.ErrUserExist) {
			s.log.Warn("app id not found", err)
			return ErrUserID, fmt.Errorf("%s: %w", op, ErrorUserExist)
		}
		log.Error("failed to get password hash: ", err)
		return ErrUserID, fmt.Errorf("%s %w", op, err)
	}

	id, err := s.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user: ", err)
		return ErrUserID, fmt.Errorf("%s %w", op, ErrFailedToSaveUser)
	}

	log.Info("user registered")

	return id, nil
}

func (s *AuthService) validate(email, password string) error {
	if email == "" {
		return ErrEmailIsEmpty
	}

	if password == "" {
		return ErrPasswordIsEmpty
	}

	return nil
}
