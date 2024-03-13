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

func TestAddProduct(t *testing.T) {
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

	h.addProduct(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())

}
