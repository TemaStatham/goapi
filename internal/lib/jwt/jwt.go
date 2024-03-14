package jwt

import (
	"errors"
	"goapi/internal/model"
	"time"

	//"github.com/golang-jwt/jwt/v5"
	"github.com/dgrijalva/jwt-go"
)

const (
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int64 `json:"user_id"`
}

func NewToken(user model.User, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		int64(user.ID),
	})

	return token.SignedString([]byte(signingKey))
}

func ParseToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}
