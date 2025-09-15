package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *Middleware) UserStatus(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		log.Println(err)
	}

	err = m.userRepo.UpdateLastActivity(userID)
	if err != nil {
		log.Println(err)
	}

	return ctx.Next()
}
