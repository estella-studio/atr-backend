package middleware

import (
	"github.com/estella-studio/leon-backend/internal/infra/jwt"
	"github.com/gofiber/fiber/v2"
)

type MiddlewareItf interface {
	Authentication(ctx *fiber.Ctx) error
}

type Middleware struct {
	jwt jwt.JWT
}

func NewMiddleware(jwt jwt.JWT) MiddlewareItf {
	return &Middleware{
		jwt: jwt,
	}
}
