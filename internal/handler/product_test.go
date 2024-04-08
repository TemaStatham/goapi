package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	service_mocks "goapi/internal/handler/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFailedAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductService := service_mocks.NewMockProductService(ctrl)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	h := NewHandler(nil, mockProductService, nil, logger)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	input := addProductType{
		Name:        "Test Product",
		Categoryies: []string{"Category1", "Category2"},
	}

	reqBody, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/product/add", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockProductService.EXPECT().AddProduct(gomock.Any(), input.Name, input.Categoryies).Return(int64(1), nil)

	h.addProduct(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
}

func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductService := service_mocks.NewMockProductService(ctrl)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	h := NewHandler(nil, mockProductService, nil, logger)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	userID := int64(123)
	c.Request = httptest.NewRequest("DELETE", "/api/product/delete", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(userCtx, userID)

	input := addProductType{
		Name:        "Test Product",
		Categoryies: []string{"Category1", "Category2"},
	}

	reqBody, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/product/add", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockProductService.EXPECT().AddProduct(gomock.Any(), input.Name, input.Categoryies).Return(int64(1), nil)

	h.addProduct(c)

	assert.Equal(t, http.StatusOK, c.Writer.Status())

}

func TestDeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductService := service_mocks.NewMockProductService(ctrl)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	h := NewHandler(nil, mockProductService, nil, logger)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	userID := int64(123)
	c.Request = httptest.NewRequest("DELETE", "/api/product/delete", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(userCtx, userID)

	input := deleteProductType{
		ID: 1,
	}

	reqBody, _ := json.Marshal(input)
	req, _ := http.NewRequest("DELETE", "/api/product/delete", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockProductService.EXPECT().DeleteProduct(gomock.Any(), input.ID).Return(nil)

	h.deleteProduct(c)

	assert.Equal(t, http.StatusOK, c.Writer.Status())
}
