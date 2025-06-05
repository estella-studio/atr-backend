package rest

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
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
	Mailer      mailer.MailerItf
}

func NewUserHandler(
	routerGroup fiber.Router, validator *validator.Validate,
	middleware middleware.MiddlewareItf, userUseCase usecase.UserUseCaseItf,
	config *env.Env, mailer mailer.MailerItf,
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
	routerGroup.Get("/renewtoken", middleware.Authentication, middleware.UserStatus, userHandler.RenewToken)
	routerGroup.Post("/friendrequest", middleware.Authentication, middleware.UserStatus, userHandler.SendFriendRequest)
	routerGroup.Get("/friendrequestsent", middleware.Authentication, middleware.UserStatus, userHandler.GetFriendRequestSent)
	routerGroup.Patch("/friendrequest", middleware.Authentication, middleware.UserStatus, userHandler.AcceptFriendRequest)
	routerGroup.Get("/friends", middleware.Authentication, middleware.UserStatus, userHandler.GetFriendList)
	routerGroup.Post("/emailverification", userHandler.NewEmailVerification)
	routerGroup.Post("/validateemail", userHandler.ValidateEmail)
	routerGroup.Get("/checkusername", userHandler.CheckUsername)
	routerGroup.Get("/info", middleware.Authentication, middleware.UserStatus, userHandler.GetUserInfo)
	routerGroup.Get("/publicinfo", userHandler.GetUserInfoPublic)
	routerGroup.Patch("/update", middleware.Authentication, middleware.UserStatus, userHandler.UpdateUserInfo)
	routerGroup.Get("/resetpassword", userHandler.ResetPassword)
	routerGroup.Post("/resetpassword", userHandler.ResetPasswordWithID)
	routerGroup.Get("/checkpasswordresetcode", userHandler.CheckPasswordResetCode)
	routerGroup.Post("/resetpasswordwithcode", userHandler.ResetPasswordWithCode)
	routerGroup.Post("/changepassword", middleware.Authentication, middleware.UserStatus, userHandler.ChangePassword)
	routerGroup.Post("/report", middleware.Authentication, middleware.UserStatus, userHandler.ReportUser)
	routerGroup.Delete("/delete", middleware.Authentication, middleware.UserStatus, userHandler.SoftDelete)
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

	_, validated, err := u.UserUseCase.GetEmailVerification(
		&dto.EmailVerification{
			Email: register.Email,
		},
	)
	if err != nil {
		log.Println(err)
	}

	if !validated {
		return fiber.NewError(
			http.StatusBadRequest,
			"email not validated",
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

func (u *UserHandler) RenewToken(ctx *fiber.Ctx) error {
	var renewToken dto.RenewToken
	var err error

	renewToken.ID, err = uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	token, err := u.UserUseCase.RenewToken(renewToken)
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "token renewed",
		"token":   token,
	})
}

func (u *UserHandler) SendFriendRequest(ctx *fiber.Ctx) error {
	var sendFriendRequest dto.SendFriendRequest

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	sendFriendRequest.UserID = userID

	err = ctx.BodyParser(&sendFriendRequest)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = u.Validator.Struct(sendFriendRequest)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	sendFriendRequest.FriendID, err = u.UserUseCase.GetUserIDFromUsername(sendFriendRequest.Username)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid username",
		)
	}

	if sendFriendRequest.UserID == sendFriendRequest.FriendID {
		return fiber.NewError(
			http.StatusBadRequest,
			"cannot add own id as friend",
		)
	}

	err = u.UserUseCase.CheckUserID(&dto.CheckUserID{ID: sendFriendRequest.FriendID})
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"can't find friend id",
		)
	}

	accepted, err := u.UserUseCase.CheckFriendRequestExist(
		&dto.CheckFriendRequestExist{
			UserID:   userID,
			FriendID: sendFriendRequest.FriendID,
		})
	if err == nil {
		if accepted {
			return fiber.NewError(
				http.StatusBadRequest,
				"user is already a friend",
			)
		} else if !accepted {
			return fiber.NewError(
				http.StatusConflict,
				"already requested",
			)
		}
	}

	accepted, err = u.UserUseCase.CheckFriendRequestFromFriend(sendFriendRequest.FriendID)
	if err == nil {
		_ = u.UserUseCase.AcceptFriendRequest(
			&dto.AcceptFriendRequest{
				UserID:   sendFriendRequest.UserID,
				FriendID: sendFriendRequest.FriendID,
			})

		if accepted {
			return fiber.NewError(
				http.StatusBadRequest,
				"user is already a friend",
			)
		} else {
			return ctx.Status(http.StatusOK).JSON(fiber.Map{
				"message": "user already sent friend request, this user will be added as friend",
			})
		}

	}

	err = u.UserUseCase.NewFriendRequest(&sendFriendRequest)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to send friend request",
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "friend request sent",
	})
}

func (u *UserHandler) GetFriendRequestSent(ctx *fiber.Ctx) error {
	var res *[]dto.ResponseGetFriendRequest

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	q := ctx.Queries()

	offset, _ := strconv.Atoi(q["offset"])

	limit, _ := strconv.Atoi(q["limit"])

	res, err = u.UserUseCase.GetFriendRequestSent(userID, offset, limit)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to retrieve friend request list",
		)
	}

	if len(*res) == 0 {
		return fiber.NewError(
			http.StatusNotFound,
			"no friend request",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved friend request list",
		"payload": res,
	})
}

func (u *UserHandler) AcceptFriendRequest(ctx *fiber.Ctx) error {
	var acceptFriendRequest dto.AcceptFriendRequest

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&acceptFriendRequest)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = u.Validator.Struct(acceptFriendRequest)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	acceptFriendRequest.UserID = userID

	acceptFriendRequest.FriendID, err = u.UserUseCase.GetUserIDFromUsername(acceptFriendRequest.Username)
	if err != nil {
		return fiber.NewError(
			http.StatusNotFound,
			"user not found",
		)
	}

	err = u.UserUseCase.AcceptFriendRequest(&acceptFriendRequest)
	if err != nil {
		return fiber.NewError(
			fiber.StatusNotFound,
			"no friend request found with current id",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "friend request accepted",
	})
}

func (u *UserHandler) GetFriendList(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	res, err := u.UserUseCase.GetFriendList(userID)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to get friend list",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "retrieved friend list",
		"payload": res,
	})
}

func (u *UserHandler) NewEmailVerification(ctx *fiber.Ctx) error {
	var emailVerification dto.EmailVerification

	err := ctx.BodyParser(&emailVerification)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = u.Validator.Struct(emailVerification)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	go func() {
		var codeString string

		for i := uint(0); i < u.Config.AccountRegistrationCodeDigitCount; i++ {
			codeString += (string)(rand.Intn(10) + 48)
		}

		newCode, _ := strconv.Atoi(codeString)

		code, success, err := u.UserUseCase.GetEmailVerification(&emailVerification)
		if err == nil && success {
			return
		} else if err != nil {
			emailVerification.ID = uuid.New()
			emailVerification.Code = uint(newCode)

			err = u.UserUseCase.NewEmailVerification(&emailVerification)
			if err != nil {
				log.Println(err)
			}
		} else {
			newCode = int(code)
		}

		err = u.Mailer.AccountRegistration(emailVerification.Email, uint(newCode))
		if err != nil {
			log.Println(err)
		}
	}()

	return ctx.Status(http.StatusOK).Context().Err()
}

func (u *UserHandler) ValidateEmail(ctx *fiber.Ctx) error {
	var validateEmail dto.ValidateEmail

	err := ctx.BodyParser(&validateEmail)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = u.Validator.Struct(validateEmail)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	err = u.UserUseCase.ValidateEmail(&validateEmail)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid code",
		)
	}

	code, success, err := u.UserUseCase.GetEmailVerification(
		&dto.EmailVerification{
			Email: validateEmail.Email,
		},
	)
	if err != nil ||
		!success ||
		(code != validateEmail.Code) {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid code",
		)
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "email confirmed",
	})
}

func (u *UserHandler) CheckUsername(ctx *fiber.Ctx) error {
	var user dto.CheckUsername

	q := ctx.Queries()

	if len(ctx.Get("Username")) != 0 {
		user.Username = ctx.Get("Username")
	} else if len(q["username"]) != 0 {
		user.Username = q["username"]
	}

	err := u.Validator.Struct(user)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	err = u.UserUseCase.CheckUsername(&user)
	if err == nil {
		return fiber.NewError(
			http.StatusConflict,
		)
	}

	return ctx.Status(http.StatusOK).Context().Err()
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

func (u *UserHandler) GetUserInfoPublic(ctx *fiber.Ctx) error {
	var getUserInfoPublic dto.GetUserInfoPublic

	q := ctx.Queries()

	getUserInfoPublic.Username = q["Username"]

	err := u.Validator.Struct(getUserInfoPublic)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	userID, err := u.UserUseCase.GetUserIDFromUsername(getUserInfoPublic.Username)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid username",
		)
	}

	res, err := u.UserUseCase.GetUserInfoPublic(userID)
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

	go func() {
		userID, err := u.UserUseCase.GetUserIDFromEmail(user)
		if err != nil {
			log.Println(err)
		}

		lastRequest, err := u.UserUseCase.GetPasswordResetCodeValidity(userID)
		if err != nil {
			log.Println(err)
		}

		timeSecond := time.Since(lastRequest)

		if timeSecond < time.Duration(u.Config.PasswordChangeCodeRetrySeconds*int(time.Second)) {
			log.Printf("time since last: %d\n", timeSecond)
			return
		}

		changeID := uuid.New()
		var codeString string

		for i := uint(0); i < u.Config.PasswordChangeCodeDigitcount; i++ {
			codeString += (string)(rand.Intn(10) + 48)
		}

		code, _ := strconv.Atoi(codeString)

		err = u.UserUseCase.CreatePasswordChangeEntry(changeID, userID)
		if err != nil {
			log.Println(err)
		}

		err = u.UserUseCase.CreatePasswordResetCode(changeID, userID, uint(code))
		if err != nil {
			log.Println(err)
		}

		err = u.UserUseCase.ResetPassword(user)
		if err == nil {
			err := u.Mailer.PasswordReset(user.Email, changeID, uint(code))
			if err != nil {
				log.Println(err)
			}
		}
	}()

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
		time.Since(createdAt) > time.Duration(u.Config.PasswordChangeExpiryMinutes*int(time.Minute)) {
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

	go func() {
		err := u.UserUseCase.UpdatePasswordChangeEntry(id, userID)
		if err != nil {
			log.Println(err)
		}
	}()

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "password changed",
	})
}

func (u *UserHandler) CheckPasswordResetCode(ctx *fiber.Ctx) error {
	var checkPasswordResetCode dto.CheckPasswordResetCode

	err := ctx.BodyParser(&checkPasswordResetCode)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = u.Validator.Struct(checkPasswordResetCode)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	userID, err := u.UserUseCase.GetUserIDFromEmail(
		dto.ResetPassword{
			Email: checkPasswordResetCode.Email,
		},
	)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid email",
		)
	}

	_, code, err := u.UserUseCase.GetPasswordResetCode(userID, checkPasswordResetCode.Code)
	if err != nil ||
		code != checkPasswordResetCode.Code {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid code",
		)
	}

	return ctx.Status(http.StatusOK).Context().Err()
}

func (u *UserHandler) ResetPasswordWithCode(ctx *fiber.Ctx) error {
	var user dto.ResetPasswordWithCode

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

	userID, err := u.UserUseCase.GetUserIDFromEmail(dto.ResetPassword{Email: user.Email})
	if err != nil {
		return fiber.NewError(
			http.StatusNoContent,
		)
	}

	passwordChagneID, code, err := u.UserUseCase.GetPasswordResetCode(userID, user.Code)
	if err != nil ||
		code != user.Code {
		return fiber.NewError(
			http.StatusUnauthorized,
			"invalid code",
		)
	}

	user.PasswordChangeId = passwordChagneID

	err = u.UserUseCase.ChangePassword(dto.ChangePassword{Password: user.Password}, userID)
	if err != nil {
		return fiber.NewError(
			http.StatusInternalServerError,
			"failed to change passsword",
		)
	}

	go func() {
		log.Println(user.PasswordChangeId)
		err := u.UserUseCase.UpdatePasswordResetCode(user.PasswordChangeId, userID, user.Code)
		if err != nil {
			log.Println(err)
		}

		err = u.UserUseCase.UpdatePasswordChangeEntry(user.PasswordChangeId, userID)
		if err != nil {
			log.Println(err)
		}
	}()

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

	go func() {
		err := u.UserUseCase.CreatePasswordChangeEntry(changeID, userID)
		if err != nil {
			log.Println(err)
		}
	}()

	err = u.UserUseCase.ChangePassword(user, userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to change password",
		})
	}

	go func() {
		err := u.UserUseCase.UpdatePasswordChangeEntry(changeID, userID)
		if err != nil {
			log.Println(err)
		}
	}()

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "password changed",
	})
}

func (u *UserHandler) ReportUser(ctx *fiber.Ctx) error {
	var reportUser dto.ReportUser

	userID, err := uuid.Parse(ctx.Locals("userID").(string))
	if err != nil {
		return fiber.NewError(
			http.StatusUnauthorized,
			"user unauthorized",
		)
	}

	err = ctx.BodyParser(&reportUser)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"failed to parse request body",
		)
	}

	err = u.Validator.Struct(reportUser)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	reportUser.ReporterID = userID

	reportUser.UserID, err = u.UserUseCase.GetUserIDFromUsername(reportUser.Username)
	if err != nil {
		return fiber.NewError(
			http.StatusBadRequest,
			"invalid username",
		)
	}

	if reportUser.UserID == reportUser.ReporterID {
		return fiber.NewError(
			http.StatusBadRequest,
			"cannot report own id",
		)
	}

	err = u.UserUseCase.ReportUser(reportUser)
	if err != nil {
		return fiber.NewError(
			http.StatusConflict,
			err.Error(),
		)
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "user reported",
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

	return ctx.Status(http.StatusNoContent).Context().Err()
}
