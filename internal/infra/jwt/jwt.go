package jwt

import (
	"time"

	"github.com/estella-studio/atr-backend/internal/infra/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTItf interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
}

type JWT struct {
	secretKey   string
	expiredTime uint
}

type Claims struct {
	ID uuid.UUID
	jwt.RegisteredClaims
}

func NewJWT(env *env.Env) *JWT {
	return &JWT{
		secretKey:   env.JWTSecretKey,
		expiredTime: env.JWTExpiredDays,
	}
}

func (j *JWT) GenerateToken(userID uuid.UUID) (string, error) {
	claim := Claims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Hour * 24 * time.Duration(j.expiredTime)),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)

	tokenStirng, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenStirng, nil
}

func (j *JWT) ValidateToken(tokenString string) (uuid.UUID, error) {
	var claims Claims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, err
	}

	userID := claims.ID

	return userID, nil
}
