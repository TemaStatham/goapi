package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_service "goapi/internal/service/mock"
	"log/slog"
	"os"
	"testing"
)

func TestCategoryService_AddCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdder := mock_service.NewMockAdderCategory(ctrl)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	categoryService := NewCategoryService(mockAdder, nil, nil, nil, mockLogger)

	testName := "Test Category"
	expectedID := int64(1)

	mockAdder.EXPECT().AddCategory(gomock.Any(), testName).Return(expectedID, nil)

	categoryID, err := categoryService.AddCategory(context.Background(), testName)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, categoryID)
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeleter := mock_service.NewMockDeleterCategory(ctrl)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	categoryService := NewCategoryService(nil, mockDeleter, nil, nil, mockLogger)

	testID := int64(1)

	mockDeleter.EXPECT().DeleteCategory(gomock.Any(), testID).Return(nil)

	err := categoryService.DeleteCategory(context.Background(), testID)
	assert.NoError(t, err)
}

func TestCategoryService_EditCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mock_service.NewMockUpdaterCategory(ctrl)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	categoryService := NewCategoryService(nil, nil, mockUpdater, nil, mockLogger)

	testID := int64(1)
	testName := "Updated Category"
	expectedID := int64(1)

	mockUpdater.EXPECT().UpdateCategoryName(gomock.Any(), testID, testName).Return(expectedID, nil)

	categoryID, err := categoryService.EditCategory(context.Background(), testID, testName)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, categoryID)
}
