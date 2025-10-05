package rest

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	datausecase "github.com/estella-studio/atr-backend/internal/app/data/usecase"
	userusecase "github.com/estella-studio/atr-backend/internal/app/user/usecase"
	"github.com/estella-studio/atr-backend/internal/domain/dto"
	"github.com/estella-studio/atr-backend/internal/infra/env"
	"github.com/estella-studio/atr-backend/internal/infra/s3"
	"github.com/estella-studio/atr-backend/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DataHandler struct {
	Validator   *validator.Validate
	Middleware  middleware.MiddlewareItf
	DataUseCase datausecase.DataUseCaseItf
	UserUseCase userusecase.UserUseCaseItf
	Env         *env.Env
	S3          s3.S3Itf
}

func NewDataHandler(
	routerGroup fiber.Router, validator *validator.Validate,
	middleware middleware.MiddlewareItf, dataUseCase datausecase.DataUseCaseItf,
	userUseCase userusecase.UserUseCaseItf, env *env.Env,
	s3 s3.S3Itf,
) {
	dataHandler := DataHandler{
		Validator:   validator,
		Middleware:  middleware,
		DataUseCase: dataUseCase,
		UserUseCase: userUseCase,
		Env:         env,
		S3:          s3,
	}

	routerGroup = routerGroup.Group("/data")

	routerGroup.Post("/add", middleware.Authentication, middleware.UserStatus, dataHandler.Add)
	routerGroup.Get("/get", middleware.Authentication, middleware.UserStatus, dataHandler.Retrieve)
	routerGroup.Get("/list", middleware.Authentication, middleware.UserStatus, dataHandler.List)
	routerGroup.Get("/listpublic", middleware.Authentication, middleware.UserStatus, dataHandler.ListPublic)
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

	add.Type, err = strconv.ParseBool(ctx.Get("X-Type"))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	err = d.Validator.Struct(add)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"cannot get data",
		)
	}

	fileContent, err := file.Open()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to open file")
	}

	byteContainer, err := io.ReadAll(fileContent)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to read file",
		)
	}

	add.ID = uuid.New()
	add.UserID = userID
	add.Data = fmt.Sprintf("%s/%v", d.Env.S3BucketURLPrefix, add.ID)

	res, err := d.DataUseCase.Add(add)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to save data",
		)
	}

	go func() {
		err = d.S3.Upload(context.Background(), add.ID.String(), byteContainer)
		if err != nil {
			log.Println(err)
		}
	}()

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "data saved",
		"payload": res,
	})
}

func (d *DataHandler) Retrieve(ctx *fiber.Ctx) error {
	var retrieve dto.Retrieve

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&retrieve)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = d.Validator.Struct(retrieve)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	retrieve.UserID = userID

	res, err := d.DataUseCase.Retrieve(retrieve)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return fiber.NewError(
				http.StatusNotFound,
				"no save data found with current id",
			)
		}

		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to retrieve save data",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved save data",
		"payload": res,
	})
}

func (d *DataHandler) List(ctx *fiber.Ctx) error {
	var res *[]dto.ResponseList

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	offset, _ := strconv.Atoi(ctx.Get("X-Offset"))

	limit, _ := strconv.Atoi(ctx.Get("X-Limit"))

	res, err = d.DataUseCase.List(userID, offset, limit)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to retrieve save data list",
		)
	}

	if len(*res) == 0 {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "no save data found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved save data list",
		"payload": res,
	})
}

func (d *DataHandler) ListPublic(ctx *fiber.Ctx) error {
	var res *[]dto.ResponseList

	userID, err := d.UserUseCase.GetUserIDFromUsername(ctx.Get("X-Username"))
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid username",
		)
	}

	offset, _ := strconv.Atoi(ctx.Get("X-Offset"))

	limit, _ := strconv.Atoi(ctx.Get("X-Limit"))

	res, err = d.DataUseCase.ListPublic(userID, offset, limit)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to retrieve save data list",
		)
	}

	if len(*res) == 0 {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "no save data found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved save data list",
		"payload": res,
	})
}
