package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"goapi/internal/model"
	"goapi/internal/repository"
	mock_service "goapi/internal/service/mock"
	"log/slog"
	"os"
	"testing"
)

func TestAddProduct(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAdderProduct, name string, categories []string)

	tests := []struct {
		name            string
		inputName       string
		inputCategories []string
		mockBehavior    mockBehavior
		expectedError   error
		expectedID      int64
	}{
		{
			name:            "Ok",
			inputName:       "Test Product",
			inputCategories: []string{"Category1", "Category2"},
			mockBehavior: func(r *mock_service.MockAdderProduct, name string, categories []string) {
				r.EXPECT().AddProduct(
					gomock.Any(),
					name,
					categories,
				).Return(int64(1), nil)
			},
			expectedError: nil,
			expectedID:    1,
		},
		{
			name:            "Empty Name",
			inputName:       "",
			inputCategories: []string{"Category1", "Category2"},
			mockBehavior:    func(r *mock_service.MockAdderProduct, name string, categories []string) {},
			expectedError:   fmt.Errorf("%s %w", "product.AddProduct", ErrProductNameIsEmpty),
			expectedID:      ErrProductId,
		},
		{
			name:            "Empty Categories",
			inputName:       "Test Product",
			inputCategories: []string{},
			mockBehavior:    func(r *mock_service.MockAdderProduct, name string, categories []string) {},
			expectedError:   fmt.Errorf("%s %w", "product.AddProduct", ErrCategoryiesEmpty),
			expectedID:      ErrProductId,
		},
		{
			name:            "Service Error",
			inputName:       "Test Product",
			inputCategories: []string{"Category1", "Category2"},
			mockBehavior: func(r *mock_service.MockAdderProduct, name string, categories []string) {
				r.EXPECT().AddProduct(
					gomock.Any(),
					name,
					categories,
				).Return(int64(-1), errors.New("something wrong"))
			},
			expectedError: fmt.Errorf("%s %w",
				"product.AddProduct", errors.New("something wrong")),
			expectedID: ErrProductId,
		},
		{
			name:            "Error Save Product",
			inputName:       "Test Product",
			inputCategories: []string{"Category1", "Category2"},
			mockBehavior: func(r *mock_service.MockAdderProduct, name string, categories []string) {
				r.EXPECT().AddProduct(
					gomock.Any(),
					name,
					categories,
				).Return(int64(-1), repository.ErrSaveProduct)
			},
			expectedError: fmt.Errorf("%s %w",
				"product.AddProduct", ErrorInvalidCredentials),
			expectedID: ErrProductId,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAdder := mock_service.NewMockAdderProduct(ctrl)
			mockLogger := slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)

			test.mockBehavior(mockAdder, test.inputName, test.inputCategories)

			productService := NewProductService(
				mockAdder,
				nil,
				nil,
				nil,
				mockLogger,
			)

			productID, err := productService.AddProduct(
				context.Background(),
				test.inputName,
				test.inputCategories,
			)

			assert.Equal(t, test.expectedID, productID)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	type mockBehavior func(r *mock_service.MockDeleterProduct, id int64)

	tests := []struct {
		name          string
		inputID       int64
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name:    "Ok",
			inputID: 1,
			mockBehavior: func(r *mock_service.MockDeleterProduct, id int64) {
				r.EXPECT().DeleteProduct(gomock.Any(), id).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:          "Invalid ID",
			inputID:       0,
			mockBehavior:  func(r *mock_service.MockDeleterProduct, id int64) {},
			expectedError: fmt.Errorf("%s %w", "product.DeleteProduct", ErrProductIDIsEmpty),
		},
		{
			name:    "Service Error",
			inputID: 1,
			mockBehavior: func(r *mock_service.MockDeleterProduct, id int64) {
				r.EXPECT().DeleteProduct(gomock.Any(), id).Return(errors.New("something error"))
			},
			expectedError: fmt.Errorf("%s %w", "product.DeleteProduct", errors.New("something error")),
		},
		{
			name:    "Error Delete Product",
			inputID: 1,
			mockBehavior: func(r *mock_service.MockDeleterProduct, id int64) {
				r.EXPECT().DeleteProduct(gomock.Any(), id).Return(repository.ErrDeleteProduct)
			},
			expectedError: fmt.Errorf("%s %w", "product.DeleteProduct", ErrorInvalidCredentials),
		},
		{
			name:    "Error Delete Product Category",
			inputID: 1,
			mockBehavior: func(r *mock_service.MockDeleterProduct, id int64) {
				r.EXPECT().DeleteProduct(gomock.Any(), id).Return(repository.ErrDeleteProductCategory)
			},
			expectedError: fmt.Errorf("%s %w", "product.DeleteProduct", ErrorInvalidCredentials),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDeleter := mock_service.NewMockDeleterProduct(ctrl)
			mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			test.mockBehavior(mockDeleter, test.inputID)

			productService := NewProductService(nil, mockDeleter, nil, nil, mockLogger)

			err := productService.DeleteProduct(context.Background(), test.inputID)

			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestEditProductName(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUpdaterProduct, id int64, name string)

	tests := []struct {
		name          string
		inputID       int64
		inputName     string
		mockBehavior  mockBehavior
		expectedID    int64
		expectedError error
	}{
		{
			name:      "Ok",
			inputID:   1,
			inputName: "Test",
			mockBehavior: func(r *mock_service.MockUpdaterProduct, id int64, name string) {
				r.EXPECT().UpdateProductName(gomock.Any(), id, name).Return(int64(1), nil)
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name:          "Invalid ID",
			inputID:       0,
			inputName:     "Test",
			mockBehavior:  func(r *mock_service.MockUpdaterProduct, id int64, name string) {},
			expectedID:    ErrProductId,
			expectedError: fmt.Errorf("%s %w", "product.EditProductName", ErrProductIDIsEmpty),
		},
		{
			name:          "Name Empty Error",
			inputID:       1,
			inputName:     "",
			mockBehavior:  func(r *mock_service.MockUpdaterProduct, id int64, name string) {},
			expectedID:    ErrProductId,
			expectedError: fmt.Errorf("%s %w", "product.EditProductName", ErrProductNameIsEmpty),
		},
		{
			name:      "Service Error",
			inputID:   1,
			inputName: "Test",
			mockBehavior: func(r *mock_service.MockUpdaterProduct, id int64, name string) {
				r.EXPECT().UpdateProductName(gomock.Any(), id, name).Return(int64(ErrProductId), errors.New("something error"))
			},
			expectedID:    ErrProductId,
			expectedError: fmt.Errorf("%s %w", "product.EditProductName", errors.New("something error")),
		},
		{
			name:      "Error Update Product",
			inputID:   1,
			inputName: "Test",
			mockBehavior: func(r *mock_service.MockUpdaterProduct, id int64, name string) {
				r.EXPECT().UpdateProductName(gomock.Any(), id, name).Return(int64(ErrProductId), repository.ErrUpdateProduct)
			},
			expectedID:    ErrProductId,
			expectedError: fmt.Errorf("%s %w", "product.EditProductName", ErrorInvalidCredentials),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUpdater := mock_service.NewMockUpdaterProduct(ctrl)
			mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			test.mockBehavior(mockUpdater, test.inputID, test.inputName)

			productService := NewProductService(nil, nil, mockUpdater, nil, mockLogger)

			id, err := productService.EditProductName(context.Background(), test.inputID, test.inputName)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedID, id)
		})
	}
}

func TestEditProductCategoryies(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUpdaterProduct, id int64, category []model.Category)

	tests := []struct {
		name          string
		inputID       int64
		inputCategory []model.Category
		mockBehavior  mockBehavior
		expectedID    int64
		expectedError error
	}{
		{
			name:          "Ok",
			inputID:       1,
			inputCategory: []model.Category{model.Category{Name: "Category1"}, model.Category{Name: "Category2"}},
			mockBehavior: func(r *mock_service.MockUpdaterProduct, id int64, category []model.Category) {
				r.EXPECT().UpdateProductCategoryies(gomock.Any(), id, category).Return(int64(1), nil)
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name:          "Invalid ID",
			inputID:       0,
			inputCategory: []model.Category{model.Category{Name: "Category1"}, model.Category{Name: "Category2"}},
			mockBehavior:  func(r *mock_service.MockUpdaterProduct, id int64, category []model.Category) {},
			expectedID:    ErrProductId,
			expectedError: fmt.Errorf("%s %w", "product.EditProductCategoryies", ErrProductIDIsEmpty),
		},
		{
			name:          "Categoryies Empty Error",
			inputID:       1,
			inputCategory: []model.Category{},
			mockBehavior:  func(r *mock_service.MockUpdaterProduct, id int64, category []model.Category) {},
			expectedID:    ErrProductId,
			expectedError: fmt.Errorf("%s %w", "product.EditProductCategoryies", ErrCategoryiesEmpty),
		},
		{
			name:          "Service Error",
			inputID:       1,
			inputCategory: []model.Category{model.Category{Name: "Category1"}, model.Category{Name: "Category2"}},
			mockBehavior: func(r *mock_service.MockUpdaterProduct, id int64, category []model.Category) {
				r.EXPECT().UpdateProductCategoryies(gomock.Any(), id, category).Return(int64(-1), errors.New("something error"))
			},
			expectedID:    ErrProductId,
			expectedError: fmt.Errorf("%s %w", "product.EditProductCategoryies", errors.New("something error")),
		},
		{
			name:          "Error Delete Product Categoryies",
			inputID:       1,
			inputCategory: []model.Category{model.Category{Name: "Category1"}, model.Category{Name: "Category2"}},
			mockBehavior: func(r *mock_service.MockUpdaterProduct, id int64, category []model.Category) {
				r.EXPECT().UpdateProductCategoryies(gomock.Any(), id, category).Return(int64(1), repository.ErrDeleteProductCategory)
			},
			expectedID:    ErrProductId,
			expectedError: fmt.Errorf("%s %w", "product.EditProductCategoryies", ErrorInvalidCredentials),
		},
		{
			name:          "Error Save Product Categoryies",
			inputID:       1,
			inputCategory: []model.Category{model.Category{Name: "Category1"}, model.Category{Name: "Category2"}},
			mockBehavior: func(r *mock_service.MockUpdaterProduct, id int64, category []model.Category) {
				r.EXPECT().UpdateProductCategoryies(gomock.Any(), id, category).Return(int64(1), repository.ErrSaveProductCategory)
			},
			expectedID:    ErrProductId,
			expectedError: fmt.Errorf("%s %w", "product.EditProductCategoryies", ErrorInvalidCredentials),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUpdater := mock_service.NewMockUpdaterProduct(ctrl)
			mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			test.mockBehavior(mockUpdater, test.inputID, test.inputCategory)

			productService := NewProductService(nil, nil, mockUpdater, nil, mockLogger)

			id, err := productService.EditProductCategory(context.Background(), test.inputID, test.inputCategory)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedID, id)
		})
	}
}

func TestGetAllProducts(t *testing.T) {
	type mockBehavior func(r *mock_service.MockGetterProduct)

	tests := []struct {
		name             string
		tag              string
		mockBehavior     mockBehavior
		expectedProducts []model.Product
		expectedError    error
	}{
		{
			name: "Ok",
			tag:  TagGetAllProducts,
			mockBehavior: func(r *mock_service.MockGetterProduct) {
				r.EXPECT().GetAllProducts(gomock.Any()).Return([]model.Product{model.Product{Name: "Test"}}, nil)
			},
			expectedProducts: []model.Product{model.Product{Name: "Test"}},
			expectedError:    nil,
		},
		{
			name:             "Invalid Tag",
			tag:              "Random Tag",
			mockBehavior:     func(r *mock_service.MockGetterProduct) {},
			expectedProducts: []model.Product{},
			expectedError:    fmt.Errorf("%s %w", "product.GetAllProducts", ErrProductUnknownTag),
		},
		{
			name: "Product Not Found",
			tag:  TagGetAllProducts,
			mockBehavior: func(r *mock_service.MockGetterProduct) {
				r.EXPECT().GetAllProducts(gomock.Any()).Return([]model.Product{}, repository.ErrProductNotFound)
			},
			expectedProducts: []model.Product{},
			expectedError:    fmt.Errorf("%s %w", "product.GetAllProducts", ErrorInvalidCredentials),
		},
		{
			name: "Service Error",
			tag:  TagGetAllProducts,
			mockBehavior: func(r *mock_service.MockGetterProduct) {
				r.EXPECT().GetAllProducts(gomock.Any()).Return([]model.Product{}, errors.New("something error"))
			},
			expectedProducts: []model.Product{},
			expectedError:    fmt.Errorf("%s %w", "product.GetAllProducts", errors.New("something error")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGetter := mock_service.NewMockGetterProduct(ctrl)
			mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			test.mockBehavior(mockGetter)

			productService := NewProductService(nil, nil, nil, mockGetter, mockLogger)

			products, err := productService.GetAllProducts(context.Background(), test.tag)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedProducts, products)
		})
	}
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
