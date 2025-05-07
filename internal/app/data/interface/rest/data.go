package rest

import (
	"io"
	"net/http"

	"github.com/estella-studio/leon-backend/internal/app/data/usecase"
	"github.com/estella-studio/leon-backend/internal/domain/dto"
	"github.com/estella-studio/leon-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DataHandler struct {
	Middleware  middleware.MiddlewareItf
	DataUseCase usecase.DataUseCaseItf
}

func NewDataHandler(
	routerGroup fiber.Router, middleware middleware.MiddlewareItf,
	dataUseCase usecase.DataUseCaseItf,
) {
	dataHandler := DataHandler{
		Middleware:  middleware,
		DataUseCase: dataUseCase,
	}

	routerGroup = routerGroup.Group("/data")

	routerGroup.Post("/add", middleware.Authentication, dataHandler.Add)
}

func (d *DataHandler) Add(ctx *fiber.Ctx) error {
	var add dto.Add

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	file, err := ctx.FormFile("data")
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"cannot get data",
		)
	}

	fileContent, err := file.Open()
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to open file",
		)
	}

	byteContainer, err := io.ReadAll(fileContent)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to read file",
		)
	}

	add.UserID = userID
	add.Data = byteContainer

	res, err := d.DataUseCase.Add(add)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to save data",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "data saved",
		"payload": res,
	})
}
