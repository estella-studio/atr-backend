package rest

import (
	"net/http"

	"github.com/estella-studio/leon-backend/internal/app/user/usecase"
	"github.com/estella-studio/leon-backend/internal/domain/dto"
	"github.com/estella-studio/leon-backend/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	Validator   *validator.Validate
	Middleware  middleware.MiddlewareItf
	userUseCase usecase.UserUseCaseItf
}

func NewUserHandler(
	routerGroup fiber.Router, validator *validator.Validate,
	middleware middleware.MiddlewareItf, userUseCase usecase.UserUseCaseItf,
) {
	userHandler := UserHandler{
		Validator:   validator,
		Middleware:  middleware,
		userUseCase: userUseCase,
	}

	routerGroup = routerGroup.Group("/users")

	routerGroup.Post("/register", userHandler.Register)
	routerGroup.Post("/login", userHandler.Login)
	routerGroup.Get("/info", middleware.Authentication, userHandler.GetUserInfo)
	routerGroup.Patch("/update", middleware.Authentication, userHandler.UpdateUserInfo)
}

func (u *UserHandler) Register(ctx *fiber.Ctx) error {
	var register dto.Register

	err := ctx.BodyParser(&register)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = u.Validator.Struct(register)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, err := u.userUseCase.Register(register)
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"please use another email / username",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "user registered",
		"payload": res,
	})
}

func (u *UserHandler) Login(ctx *fiber.Ctx) error {
	var login dto.Login

	err := ctx.BodyParser(&login)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = u.Validator.Struct(login)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	res, token, err := u.userUseCase.Login(login)
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"invalid username or password",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user authenticated",
		"token":   token,
		"payload": res,
	})
}

func (u *UserHandler) GetUserInfo(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, "user unauthorized")
	}

	res, err := u.userUseCase.GetUserInfo(userID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to get user info")
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved user info",
		"payload": res,
	})
}

func (u *UserHandler) UpdateUserInfo(ctx *fiber.Ctx) error {
	var user dto.UpdateUserInfo

	err := ctx.BodyParser(&user)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = u.Validator.Struct(user)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	_, err = u.userUseCase.UpdateUserInfo(user, userID)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to update user info",
		)
	}

	res, err := u.userUseCase.GetUserInfo(userID)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"user info updated but failed to retrieve updated content")
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user updated",
		"payload": res,
	})
}
