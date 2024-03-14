package jwt

import (
	"github.com/stretchr/testify/assert"
	"goapi/internal/model"
	"testing"
	"time"
)

func TestNewTokenAndParseToken(t *testing.T) {
	user := model.User{
		ID:       123,
		Email:    "test@example.com",
		PassHash: []byte("test"),
	}

	duration := time.Hour
	tokenString, err := NewToken(user, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	userID, err := ParseToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, int64(user.ID), userID)
}
