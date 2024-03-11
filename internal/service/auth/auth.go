package auth

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

var (
	ErrorInvalidCredentials = errors.New("invalid credential")
	ErrorUserExist          = errors.New("user already exist")
)

type Service struct {
	usrSaver    UserSaver
	usrProvider UserProvider
	log         *slog.Logger
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (ui int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (model.User, error)
}

func NewService(us UserSaver, up UserProvider, l *slog.Logger, tokenTTL time.Duration) *Service {
	return &Service{
		usrSaver:    us,
		usrProvider: up,
		log:         l,
		tokenTTL:    tokenTTL,
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	const op = "auth.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("logging user")

	user, err := s.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.log.Warn("user not found", err)
			return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		s.log.Warn("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
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

func (s *Service) Register(ctx context.Context, email, password string) (userID int64, err error) {
	const op = "auth.RegisterNewUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, repository.ErrUserExist) {
			s.log.Warn("app id not found", err)
			return 0, fmt.Errorf("%s: %w", op, ErrorUserExist)
		}
		log.Error("failed to get password hash: ", err)
		return 0, fmt.Errorf("%s %w", op, err)
	}

	id, err := s.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user: ", err)
	}

	log.Info("user registered")

	return id, nil
}
