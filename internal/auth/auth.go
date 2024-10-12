package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	SecretKey []byte
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func New(s string) *Auth {
	return &Auth{SecretKey: []byte(s)}
}

func (a *Auth) CreateToken(userId string) (*TokenPair, error) {
	accessTokenClaims := jwt.MapClaims{
		"userId": userId,
	}
	at, err := a.generateToken(accessTokenClaims, a.SecretKey, time.Minute*15) // 15-minute expiration
	if err != nil {
		return nil, fmt.Errorf("error signing access token: %v", err)
	}

	refreshTokenClaims := jwt.MapClaims{
		"sub": 1,
	}
	rt, err := a.generateToken(refreshTokenClaims, []byte("secret"), time.Hour*24) // 24-hour expiration
	if err != nil {
		return nil, fmt.Errorf("error signing refresh token: %v", err)
	}

	return &TokenPair{AccessToken: at, RefreshToken: rt}, nil
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

func (a *Auth) RefreshToken(oldRefreshToken string) (*TokenPair, error) {
	refreshToken, err := a.parseAndVerifyToken(oldRefreshToken, []byte("secret"))
	if err != nil {
		return nil, fmt.Errorf("error parsing refresh token: %v", err)
	}

	if !refreshToken.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] != 1 {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	newAccessTokenClaims := jwt.MapClaims{
		"userId": claims["userId"],
	}
	at, err := a.generateToken(newAccessTokenClaims, a.SecretKey, time.Minute*15)
	if err != nil {
		return nil, fmt.Errorf("error signing new access token: %v", err)
	}

	newRefreshTokenClaims := jwt.MapClaims{
		"sub":    1,
		"userId": claims["userId"],
	}
	rt, err := a.generateToken(newRefreshTokenClaims, []byte("secret"), time.Hour*24)
	if err != nil {
		return nil, fmt.Errorf("error signing new refresh token: %v", err)
	}

	return &TokenPair{AccessToken: at, RefreshToken: rt}, nil
}

func (a *Auth) generateToken(claims jwt.MapClaims, secretKey []byte, duration time.Duration) (string, error) {
	claims["exp"] = time.Now().Add(duration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (a *Auth) parseAndVerifyToken(tokenStr string, secretKey []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
}
