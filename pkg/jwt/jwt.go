package jwt

import (
	"errors"
	"time"

	"github.com/brandonbraner/maas/pkg/permissions"
	"github.com/golang-jwt/jwt/v5"
)

var (
	tokenExpiration = time.Hour * 24            //TODO update later with something from config
	secretKey       = []byte("your-secret-key") //TODO Update with something from config that gets set from the env
	ErrInvalidToken = errors.New("invalid token")
)

type Claims struct {
	UserID                  string `json:"user_id"`
	Email                   string `json:"email"`
	permissions.Permissions `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, email string, permissions permissions.Permissions) (string, error) {
	claims := &Claims{
		UserID:      userID,
		Email:       email,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
