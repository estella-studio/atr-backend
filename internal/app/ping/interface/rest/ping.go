package rest

import (
	"net/http"

	"github.com/estella-studio/leon-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type PingHandler struct {
	Middleware middleware.MiddlewareItf
}

func NewPingHandler(routerGroup fiber.Router, middleware middleware.MiddlewareItf) {
	pingHandler := PingHandler{
		Middleware: middleware,
	}

	routerGroup = routerGroup.Group("/ping")

	routerGroup.Get("/", middleware.Authentication, middleware.UserStatus, pingHandler.Ping)
}

func (p *PingHandler) Ping(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).Context().Err()
}
