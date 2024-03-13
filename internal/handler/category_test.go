package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	service_mocks "goapi/internal/handler/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAddCategoryFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCategoryService := service_mocks.NewMockCategoryService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	h := NewHandler(nil, nil, mockCategoryService, logger)

	r := gin.Default()
	r.POST("/api/category/add", h.signIn)

	w := httptest.NewRecorder()
	reqBody := `{"name":"test"}`
	req, _ := http.NewRequest("POST", "/api/category/add", strings.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer YOUR_ACCESS_TOKEN")
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "{\"message\":\"invalid input body\"}", w.Body.String())
}
