package usecase

import (
	"github.com/estella-studio/leon-backend/internal/app/user/repository"
	"github.com/estella-studio/leon-backend/internal/domain/dto"
	"github.com/estella-studio/leon-backend/internal/domain/entity"
	"github.com/estella-studio/leon-backend/internal/infra/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCaseItf interface {
	Register(register dto.Register) (dto.ResponseRegister, error)
	Login(login dto.Login) (dto.ResponseLogin, string, error)
	GetUserInfo(userID uuid.UUID) (dto.ResponseGetUserInfo, error)
	UpdateUserInfo(updateUserInfo dto.UpdateUserInfo, userID uuid.UUID) (dto.ResponseUpdateUserInfo, error)
	SoftDelete(userID uuid.UUID) error
}

type UserUseCase struct {
	userRepo repository.UserMySQLItf
	jwt      jwt.JWTItf
}

func NewUserUseCase(userRepo repository.UserMySQLItf, jwt *jwt.JWT) UserUseCaseItf {
	return &UserUseCase{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (u *UserUseCase) Register(register dto.Register) (dto.ResponseRegister, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(register.Password),
		bcrypt.DefaultCost)
	if err != nil {
		return dto.ResponseRegister{},
			err
	}

	user := entity.User{
		ID:       uuid.New(),
		Email:    register.Email,
		Username: register.Username,
		Password: string(hashedPassword),
		Name:     register.Name,
	}

	err = u.userRepo.Register(&user)
	if err != nil {
		return dto.ResponseRegister{},
			err
	}

	return user.ParseToDTOResponseRegister(),
		nil
}

func (u *UserUseCase) Login(login dto.Login) (dto.ResponseLogin, string, error) {
	var user entity.User

	err := u.userRepo.GetUsername(&user, dto.Login{Username: login.Username})
	if err != nil {
		return dto.ResponseLogin{},
			"",
			err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return dto.ResponseLogin{},
			"",
			err
	}

	token, err := u.jwt.GenerateToken(user.ID)
	if err != nil {
		return dto.ResponseLogin{},
			"",
			err
	}

	return user.ParseToDTOResponseLogin(), token, nil
}

func (u *UserUseCase) GetUserInfo(userID uuid.UUID) (dto.ResponseGetUserInfo, error) {
	user := entity.User{
		ID: userID,
	}

	err := u.userRepo.GetUserInfo(&user)
	if err != nil {
		return dto.ResponseGetUserInfo{},
			err
	}

	return user.ParseToDTOResponseGetUserInfo(), nil
}

func (u *UserUseCase) UpdateUserInfo(updateUserInfo dto.UpdateUserInfo, userID uuid.UUID) (dto.ResponseUpdateUserInfo, error) {
	user := entity.User{
		ID:       userID,
		Email:    updateUserInfo.Email,
		Username: updateUserInfo.Username,
		Name:     updateUserInfo.Name,
	}

	err := u.userRepo.UpdateUserInfo(&user)
	if err != nil {
		return dto.ResponseUpdateUserInfo{},
			err
	}

	return user.ParseToDTOResponseUpdateUserInfo(), nil
}

func (u *UserUseCase) SoftDelete(userID uuid.UUID) error {
	user := entity.User{
		ID: userID,
	}

	err := u.userRepo.SoftDelete(&user)

	return err
}
