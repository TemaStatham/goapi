package service

import (
	"context"
	"goapi/internal/model"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_service "goapi/internal/service/mock"
)

func TestLoginOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := mock_service.NewMockUserProvider(ctrl)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	authService := NewAuthService(nil, mockUserProvider, mockLogger, time.Minute)

	testEmail := "test@example.com"
	testPassword := "1"
	mockUser := model.User{ID: 1, Email: testEmail, PassHash: []byte("$2a$10$H2R/kGmJtGZir7eYBSQPJO2Mfm3tlGY3C.3Wvt0N.HPsyYrtG0hUO")}
	mockUserProvider.EXPECT().User(gomock.Any(), testEmail).Return(mockUser, nil)

	token, err := authService.Login(context.Background(), testEmail, testPassword)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLoginFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := mock_service.NewMockUserProvider(ctrl)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	authService := NewAuthService(nil, mockUserProvider, mockLogger, time.Minute)

	testEmail := "test@example.com"
	testPassword := "0"
	mockUser := model.User{ID: 1, Email: testEmail, PassHash: []byte("$2a$10$H2R/kGmJtGZir7eYBSQPJO2Mfm3tlGY3C.3Wvt0N.HPsyYrtG0hUO")}
	mockUserProvider.EXPECT().User(gomock.Any(), testEmail).Return(mockUser, nil)

	token, err := authService.Login(context.Background(), testEmail, testPassword)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestRegisterOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSaver := mock_service.NewMockUserSaver(ctrl)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	authService := NewAuthService(mockUserSaver, nil, mockLogger, time.Minute)

	testEmail := "test@example.com"
	testPassword := "password"
	mockUserSaver.EXPECT().SaveUser(gomock.Any(), testEmail, gomock.Any()).Return(int64(1), nil)

	userID, err := authService.Register(context.Background(), testEmail, testPassword)
	assert.NoError(t, err)
	assert.NotEqual(t, ErrUserID, userID)
}
