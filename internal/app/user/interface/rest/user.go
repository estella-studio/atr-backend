package rest

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/estella-studio/leon-backend/internal/app/user/usecase"
	"github.com/estella-studio/leon-backend/internal/domain/dto"
	"github.com/estella-studio/leon-backend/internal/infra/env"
	"github.com/estella-studio/leon-backend/internal/infra/mailer"
	"github.com/estella-studio/leon-backend/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	Validator   *validator.Validate
	Middleware  middleware.MiddlewareItf
	UserUseCase usecase.UserUseCaseItf
	Config      *env.Env
	Mailer      *mailer.Mailer
}

func NewUserHandler(
	routerGroup fiber.Router, validator *validator.Validate,
	middleware middleware.MiddlewareItf, userUseCase usecase.UserUseCaseItf,
	config *env.Env, mailer *mailer.Mailer,
) {
	userHandler := UserHandler{
		Config:      config,
		Validator:   validator,
		Middleware:  middleware,
		UserUseCase: userUseCase,
		Mailer:      mailer,
	}

	routerGroup = routerGroup.Group("/users")

	routerGroup.Post("/register", userHandler.Register)
	routerGroup.Post("/login", userHandler.Login)
	routerGroup.Get("/info", middleware.Authentication, userHandler.GetUserInfo)
	routerGroup.Patch("/update", middleware.Authentication, userHandler.UpdateUserInfo)
	routerGroup.Get("/resetpassword", userHandler.ResetPassword)
	routerGroup.Post("/resetpassword", userHandler.ResetPasswordWithID)
	routerGroup.Post("/changepassword", middleware.Authentication, userHandler.ChangePassword)
	routerGroup.Delete("/delete", middleware.Authentication, userHandler.SoftDelete)
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

	res, err := u.UserUseCase.Register(register)
	if err != nil {
		return fiber.NewError(
			http.StatusConflict,
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

	res, token, err := u.UserUseCase.Login(login)
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
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := u.UserUseCase.GetUserInfo(userID)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get user info",
		)
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

	_, err = u.UserUseCase.UpdateUserInfo(user, userID)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fiber.NewError(
				http.StatusConflict,
				"please use another email / username",
			)
		}

		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to update user info",
		)
	}

	res, err := u.UserUseCase.GetUserInfo(userID)
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

func (u *UserHandler) ResetPassword(ctx *fiber.Ctx) error {
	var user dto.ResetPassword

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

	userID, err := u.UserUseCase.GetUserID(user)
	if err != nil {
		log.Println(err)
	}

	changeID := uuid.New()

	go u.UserUseCase.CreatePasswordChangeEntry(changeID, userID)

	err = u.UserUseCase.ResetPassword(user)
	if err == nil {
		go u.Mailer.PasswordReset(user.Email, changeID)
	}

	return ctx.Status(http.StatusOK).Context().Err()
}

func (u *UserHandler) ResetPasswordWithID(ctx *fiber.Ctx) error {
	var user dto.ChangePassword

	q := ctx.Queries()

	id, err := uuid.Parse(q["id"])
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid query",
		)
	}

	err = ctx.BodyParser(&user)
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

	alreadyChanged, createdAt, err := u.UserUseCase.GetPasswordChangeValidity(id)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to change password",
		)
	}

	userID, err := u.UserUseCase.GetPasswordChangeEntry(id)
	if err != nil ||
		alreadyChanged ||
		time.Since(createdAt) > time.Duration(u.Config.PasswordChangeExpiryMinute*int(time.Minute)) {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid link",
		)
	}

	err = u.UserUseCase.ChangePassword(user, userID)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to change password",
		)
	}

	go u.UserUseCase.UpdatePasswordChangeEntry(id, userID)

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "password changed",
	})
}

func (u *UserHandler) ChangePassword(ctx *fiber.Ctx) error {
	var user dto.ChangePassword

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&user)
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

	changeID := uuid.New()

	go u.UserUseCase.CreatePasswordChangeEntry(changeID, userID)

	err = u.UserUseCase.ChangePassword(user, userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to change password",
		})
	}

	go u.UserUseCase.UpdatePasswordChangeEntry(changeID, userID)

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "password changed",
	})
}

func (u *UserHandler) SoftDelete(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = u.UserUseCase.SoftDelete(userID)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to delete user",
		)
	}

	return ctx.Status(http.StatusNoContent).JSON(fiber.Map{
		"message": "user deleted",
	})
}
