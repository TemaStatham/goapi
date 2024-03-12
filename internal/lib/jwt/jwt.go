package jwt

import (
	"fmt"
	"goapi/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	key = "my_secret_key"
)

func NewToken(user model.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil // Ваш секретный ключ для подписи токена
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("failed to extract data from token\n")
	}

	userID, ok := claims["uid"].(int64)
	if !ok {
		return 0, fmt.Errorf("failed to extract user ID from token\n")
	}

	return userID, nil
}
