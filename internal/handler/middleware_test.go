package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserIdentityFailed(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer test-failed-token")

	h := &Handler{}

	h.userIdentity(c)

	userId := c.GetInt64("userId")
	assert.Empty(t, userId)
}

func TestGetUserId(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set("userId", 123)

	id, err := getUserId(c)

	assert.NoError(t, err)
	assert.Equal(t, 123, id)
}
