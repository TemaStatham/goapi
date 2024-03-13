package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestSignUpOk(t *testing.T) {
	inputParams := []struct {
		testName             string
		in                   signUpType
		expectedStatusCode   int
		expectedResponseBody string
		expectedID           int64
	}{
		{
			"Ok",
			signUpType{"test1@example.com", "password1"},
			http.StatusOK,
			"{\"id\":0}",
			0,
		},
		{
			"Ok",
			signUpType{"test2@example.com", "password2"},
			http.StatusOK,
			"{\"id\":0}",
			0,
		},
	}

	for _, params := range inputParams {
		t.Run(params.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := service_mocks.NewMockAuthService(ctrl)
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			h := NewHandler(mockAuthService, nil, nil, logger)

			r := gin.Default()
			r.POST("/auth/sign-up", h.signUp)

			mockAuthService.EXPECT().
				Register(gomock.Any(), params.in.Email, params.in.Password).
				Return(params.expectedID, nil)

			w := httptest.NewRecorder()
			reqBody := `{"email": "` + params.in.Email + `", "password": "` + params.in.Password + `"}`
			req := httptest.NewRequest("POST", "/auth/sign-up",
				bytes.NewBufferString(reqBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, params.expectedStatusCode)
			assert.Equal(t, w.Body.String(), params.expectedResponseBody)
		})
	}
}

func TestSignUpFailed(t *testing.T) {
	inputParams := []struct {
		testName             string
		in                   signUpType
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			"Ok",
			signUpType{},
			http.StatusBadRequest,
			"{\"message\":\"invalid input body\"}",
		},
	}
	for _, params := range inputParams {
		t.Run(params.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := service_mocks.NewMockAuthService(ctrl)
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			h := NewHandler(mockAuthService, nil, nil, logger)

			r := gin.Default()
			r.POST("/auth/sign-up", h.signUp)

			w := httptest.NewRecorder()
			reqBody := `{"email": "` + params.in.Email + `", "password": "` + params.in.Password + `"}`
			req := httptest.NewRequest("POST", "/auth/sign-up",
				bytes.NewBufferString(reqBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, params.expectedStatusCode)
			assert.Equal(t, w.Body.String(), params.expectedResponseBody)
		})
	}
}

func TestSignInOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service_mocks.NewMockAuthService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	h := NewHandler(mockAuthService, nil, nil, logger)

	r := gin.Default()
	r.POST("/auth/sign-in", h.signIn)

	expectedToken := "mockToken"
	mockAuthService.EXPECT().Login(gomock.Any(), "test@example.com", "password").Return(expectedToken, nil)

	w := httptest.NewRecorder()
	reqBody := `{"email": "test@example.com", "password": "password"}`
	req, _ := http.NewRequest("POST", "/auth/sign-in", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	fmt.Println(response["token"])
	assert.Equal(t, expectedToken, response["token"])
}

func TestSignInFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := service_mocks.NewMockAuthService(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	h := NewHandler(mockAuthService, nil, nil, logger)

	r := gin.Default()
	r.POST("/auth/sign-in", h.signIn)

	w := httptest.NewRecorder()
	reqBody := `{}`
	req, _ := http.NewRequest("POST", "/auth/sign-in", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Equal(t, "{\"message\":\"invalid input body\"}", w.Body.String())
}
