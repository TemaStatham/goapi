package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"goapi/internal/model"
	mock_service "goapi/internal/service/mock"
	"log/slog"
	"os"
	"testing"
)

func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdder := mock_service.NewMockAdderProduct(ctrl)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	productService := NewProductService(mockAdder, nil, nil, nil, mockLogger)

	testName := "Test Product"
	testCategories := []string{"Category1", "Category2"}
	mockAdder.EXPECT().AddProduct(gomock.Any(), testName, testCategories).Return(int64(1), nil)

	productID, err := productService.AddProduct(context.Background(), testName, testCategories)
	assert.NoError(t, err)
	assert.NotEqual(t, ErrProductId, productID)
}

func TestDeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeleter := mock_service.NewMockDeleterProduct(ctrl)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	productService := NewProductService(nil, mockDeleter, nil, nil, mockLogger)

	testProductID := int64(1)
	mockDeleter.EXPECT().DeleteProduct(gomock.Any(), testProductID).Return(nil)

	err := productService.DeleteProduct(context.Background(), testProductID)
	assert.NoError(t, err)
}

func TestEditProductName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mock_service.NewMockUpdaterProduct(ctrl)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	productService := NewProductService(nil, nil, mockUpdater, nil, mockLogger)

	testProductID := int64(1)
	testProductName := "Updated Product Name"

	mockUpdater.EXPECT().UpdateProductName(gomock.Any(), testProductID, testProductName).Return(testProductID, nil)

	productID, err := productService.EditProductName(context.Background(), testProductID, testProductName)
	assert.NoError(t, err)
	assert.Equal(t, testProductID, productID)
}

func TestGetCategoryProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetter := mock_service.NewMockGetterProduct(ctrl)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	productService := NewProductService(nil, nil, nil, mockGetter, mockLogger)

	testCategory := "Test Category"
	expectedProducts := []model.Product{
		{ID: 1, Name: "Product 1", Categoryies: []model.Category{}},
		{ID: 2, Name: "Product 2", Categoryies: []model.Category{}},
	}

	mockGetter.EXPECT().GetCategoryProducts(gomock.Any(), testCategory).Return(expectedProducts, nil)

	products, err := productService.GetCategoryProducts(context.Background(), testCategory)
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
}

func TestAddProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdder := mock_service.NewMockAdderProduct(ctrl)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	productService := NewProductService(mockAdder, nil, nil, nil, mockLogger)

	testProducts := []model.Product{
		{ID: 1, Name: "Product 1", Categoryies: []model.Category{}},
		{ID: 2, Name: "Product 2", Categoryies: []model.Category{}},
	}

	mockAdder.EXPECT().AddProducts(gomock.Any(), testProducts).Return(nil)

	err := productService.AddProducts(context.Background(), testProducts)
	assert.NoError(t, err)
}
