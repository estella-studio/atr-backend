package middleware

import (
	"github.com/estella-studio/atr-backend/internal/app/user/repository"
	"github.com/estella-studio/atr-backend/internal/infra/jwt"
	"github.com/gofiber/fiber/v2"
)

type MiddlewareItf interface {
	Authentication(ctx *fiber.Ctx) error
	UserStatus(ctx *fiber.Ctx) error
}

type Middleware struct {
	jwt      jwt.JWT
	userRepo repository.UserMySQLItf
}

func NewMiddleware(jwt jwt.JWT, userRepo repository.UserMySQLItf) MiddlewareItf {
	return &Middleware{
		jwt:      jwt,
		userRepo: userRepo,
	}
}
