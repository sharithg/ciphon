package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	SecretKey []byte
}

type Claims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func New(s string) *Auth {
	return &Auth{SecretKey: []byte(s)}
}

func (a *Auth) CreateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(a.SecretKey)
	if err != nil {
		return "", fmt.Errorf("error signing jwt: %s", err)
	}

	return tokenString, nil
}

func (a *Auth) VerifyToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.SecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing jwt: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
