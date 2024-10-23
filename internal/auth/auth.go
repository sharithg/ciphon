package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	SecretKey            []byte
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func New(s string, accessTokenLifetime time.Duration, refreshTokenLifetime time.Duration) *Auth {
	return &Auth{
		SecretKey:            []byte(s),
		AccessTokenLifetime:  accessTokenLifetime,
		RefreshTokenLifetime: refreshTokenLifetime,
	}
}

func (a *Auth) CreateToken(userId string) (*TokenPair, error) {
	at, err := a.generateTokenWithUser(userId)
	if err != nil {
		return nil, fmt.Errorf("error signing access token: %v", err)
	}

	rt, err := a.generateRefreshToken(userId)
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
	refreshToken, err := a.parseAndVerifyToken(oldRefreshToken, a.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error parsing refresh token: %v", err)
	}

	if !refreshToken.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	sub, _ := claims["sub"].(float64)

	if !ok || int(sub) != 1 {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	at, err := a.generateTokenWithUser(claims["userId"].(string))
	if err != nil {
		return nil, fmt.Errorf("error signing new access token: %v", err)
	}

	rt, err := a.generateRefreshToken(claims["userId"].(string))
	if err != nil {
		return nil, fmt.Errorf("error signing new refresh token: %v", err)
	}

	return &TokenPair{AccessToken: at, RefreshToken: rt}, nil
}

func (a *Auth) generateTokenWithUser(userId string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(a.AccessTokenLifetime).Unix(),
	}
	return a.generateToken(claims, a.SecretKey)
}

func (a *Auth) generateRefreshToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"sub":    1,
		"userId": userId,
		"exp":    time.Now().Add(a.RefreshTokenLifetime).Unix(),
	}
	return a.generateToken(claims, a.SecretKey)
}

func (a *Auth) generateToken(claims jwt.MapClaims, secretKey []byte) (string, error) {
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
