package jwt

import (
	"time"

	"github.com/estella-studio/leon-backend/internal/infra/env"
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

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func NewJWT(env *env.Env) *JWT {
	secretKey := env.JWTSecretKey
	expiredTime := env.JWTExpiredDays

	return &JWT{
		secretKey:   secretKey,
		expiredTime: expiredTime,
	}
}

func (j *JWT) GenerateToken(userID uuid.UUID) (string, error) {
	claim := Claims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(j.expiredTime * 24)),
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
	var claim Claims

	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (any, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, err
	}

	userID := claim.ID

	return userID, nil
}
