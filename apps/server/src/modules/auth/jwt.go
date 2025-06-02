package auth

import (
	"errors"
	"peekaping/src/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type TokenMaker struct {
	accessTokenSecretKey  string
	refreshTokenSecretKey string
	accessTokenDuration   time.Duration
	refreshTokenDuration  time.Duration
}

func NewTokenMaker(cfg *config.Config) *TokenMaker {
	return &TokenMaker{
		accessTokenSecretKey:  cfg.AccessTokenSecretKey,
		refreshTokenSecretKey: cfg.RefreshTokenSecretKey,
		accessTokenDuration:   cfg.AccessTokenExpiresIn,
		refreshTokenDuration:  cfg.RefreshTokenExpiresIn,
	}
}

func (maker *TokenMaker) CreateAccessToken(user *Model) (string, error) {
	return maker.createToken(user, "access", maker.accessTokenDuration, maker.accessTokenSecretKey)
}

func (maker *TokenMaker) CreateRefreshToken(user *Model) (string, error) {
	return maker.createToken(user, "refresh", maker.refreshTokenDuration, maker.refreshTokenSecretKey)
}

func (maker *TokenMaker) createToken(user *Model, tokenType string, duration time.Duration, secretKey string) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func (maker *TokenMaker) VerifyToken(tokenString string, tokenType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		if tokenType == "access" {
			return []byte(maker.accessTokenSecretKey), nil
		}
		if tokenType == "refresh" {
			return []byte(maker.refreshTokenSecretKey), nil
		}
		return nil, ErrInvalidToken
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
