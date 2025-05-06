package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type PingHandler struct{}

func NewPingHandler(routerGroup fiber.Router) {
	pingHandler := PingHandler{}

	routerGroup = routerGroup.Group("/ping")

	routerGroup.Get("/", pingHandler.Ping)
}

func (p *PingHandler) Ping(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).Context().Err()
}
