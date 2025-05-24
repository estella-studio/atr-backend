package usecase

import (
	"log"
	"time"

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
	NewEmailVerification(emailVerification *dto.EmailVerification) error
	ValidateEmail(validateEmail *dto.ValidateEmail) error
	GetEmailVerification(validateEmail *dto.EmailVerification) (uint, bool, error)
	CheckUsername(userName *dto.CheckUsername) error
	GetUserInfo(userID uuid.UUID) (dto.ResponseGetUserInfo, error)
	UpdateUserInfo(updateUserInfo dto.UpdateUserInfo, userID uuid.UUID) (dto.ResponseUpdateUserInfo, error)
	ResetPassword(resetPassword dto.ResetPassword) error
	ChangePassword(changePassword dto.ChangePassword, userID uuid.UUID) error
	CreatePasswordChangeEntry(changeID uuid.UUID, userID uuid.UUID) error
	CreatePasswordResetCode(changeID uuid.UUID, userID uuid.UUID, code uint) error
	UpdatePasswordChangeEntry(changeID uuid.UUID, userID uuid.UUID) error
	UpdatePasswordResetCode(changeID uuid.UUID, userID uuid.UUID, code uint) error
	GetPasswordResetCode(userID uuid.UUID, code uint) (uuid.UUID, uint, error)
	GetPasswordResetCodeValidity(userID uuid.UUID) (time.Time, error)
	GetPasswordChangeValidity(id uuid.UUID) (bool, time.Time, error)
	GetPasswordChangeEntry(id uuid.UUID) (uuid.UUID, error)
	GetUserID(getUserID dto.ResetPassword) (uuid.UUID, error)
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

	userDetail := entity.UserDetail{
		ID:           uuid.New(),
		UserID:       user.ID,
		ProfileIndex: register.ProfileIndex,
	}

	err = u.userRepo.Register(&user)
	if err != nil {
		return dto.ResponseRegister{},
			err
	}

	err = u.userRepo.RegisterUserDetail(&userDetail)
	if err != nil {
		return dto.ResponseRegister{},
			err
	}

	user.UserDetail = userDetail

	return user.ParseToDTOResponseRegister(), nil
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

	_ = u.userRepo.GetUserInfo(&user)

	log.Println(user)

	return user.ParseToDTOResponseLogin(), token, nil
}

func (u *UserUseCase) NewEmailVerification(emailVerification *dto.EmailVerification) error {
	verification := entity.Verification{
		ID:    emailVerification.ID,
		Email: emailVerification.Email,
		Code:  emailVerification.Code,
	}

	err := u.userRepo.NewEmailVerification(&verification)

	return err
}

func (u *UserUseCase) ValidateEmail(validateEmail *dto.ValidateEmail) error {
	verification := entity.Verification{
		Email: validateEmail.Email,
		Code:  validateEmail.Code,
	}

	err := u.userRepo.ValidateEmail(&verification)

	return err
}

func (u *UserUseCase) GetEmailVerification(validateEmail *dto.EmailVerification) (uint, bool, error) {
	verification := entity.Verification{
		Email: validateEmail.Email,
	}

	err := u.userRepo.GetEmailVerification(&verification)

	return verification.Code, verification.Success, err
}

func (u *UserUseCase) CheckUsername(userName *dto.CheckUsername) error {
	user := entity.User{
		Username: userName.Username,
	}

	err := u.userRepo.CheckUsername(&user)

	return err
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

	userDetail := entity.UserDetail{
		UserID:       userID,
		ProfileIndex: updateUserInfo.ProfileIndex,
	}

	err := u.userRepo.UpdateUserInfo(&user)
	if err != nil {
		return dto.ResponseUpdateUserInfo{},
			err
	}

	err = u.userRepo.UdpateUserDetail(&userDetail)
	if err != nil {
		log.Println(err)
	}

	return user.ParseToDTOResponseUpdateUserInfo(), nil
}

func (u *UserUseCase) ResetPassword(resetPassword dto.ResetPassword) error {
	user := entity.User{
		Email: resetPassword.Email,
	}

	err := u.userRepo.GetEmail(&user, dto.ResetPassword{Email: resetPassword.Email})

	return err
}

func (u *UserUseCase) ChangePassword(changePassword dto.ChangePassword, userID uuid.UUID) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(changePassword.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := entity.User{
		ID:       userID,
		Password: string(hashedPassword),
	}

	err = u.userRepo.ChangePassword(&user)

	return err
}

func (u *UserUseCase) CreatePasswordChangeEntry(changeID uuid.UUID, userID uuid.UUID) error {
	passwordChange := entity.PasswordChange{
		ID:      changeID,
		UserID:  userID,
		Success: false,
	}

	err := u.userRepo.CreatePasswordChangeEntry(&passwordChange)

	return err
}

func (u *UserUseCase) CreatePasswordResetCode(changeID uuid.UUID, userID uuid.UUID, code uint) error {
	passwordReset := entity.PasswordResetCode{
		ID:               uuid.New(),
		PasswordChangeID: changeID,
		UserID:           userID,
		Code:             code,
	}

	err := u.userRepo.CreatePasswordResetCode(&passwordReset)

	return err
}

func (u *UserUseCase) UpdatePasswordChangeEntry(changeID uuid.UUID, userID uuid.UUID) error {
	passwordChange := entity.PasswordChange{
		ID:      changeID,
		UserID:  userID,
		Success: true,
	}

	err := u.userRepo.UpdatePasswordChangeEntry(&passwordChange)

	return err
}

func (u *UserUseCase) UpdatePasswordResetCode(changeID uuid.UUID, userID uuid.UUID, code uint) error {
	passwordResetCode := entity.PasswordResetCode{
		PasswordChangeID: changeID,
		UserID:           userID,
		Code:             code,
	}

	err := u.userRepo.UpdatePasswordResetCode(&passwordResetCode)

	return err
}

func (u *UserUseCase) GetPasswordResetCode(userID uuid.UUID, code uint) (uuid.UUID, uint, error) {
	passwordResetCode := entity.PasswordResetCode{
		UserID: userID,
	}

	err := u.userRepo.GetPasswordResetCode(
		&passwordResetCode,
	)

	return passwordResetCode.PasswordChangeID, passwordResetCode.Code, err
}

func (u *UserUseCase) GetPasswordResetCodeValidity(userID uuid.UUID) (time.Time, error) {
	passwordResetCode := entity.PasswordResetCode{
		UserID: userID,
	}

	err := u.userRepo.GetPasswordResetCodeValidity(&passwordResetCode)

	return passwordResetCode.CreatedAt, err
}

func (u *UserUseCase) GetPasswordChangeValidity(id uuid.UUID) (bool, time.Time, error) {
	passwordChange := entity.PasswordChange{
		ID: id,
	}

	err := u.userRepo.GetPasswordChangeValidity(&passwordChange)

	return passwordChange.Success, passwordChange.CreatedAt, err
}

func (u *UserUseCase) GetPasswordChangeEntry(id uuid.UUID) (uuid.UUID, error) {
	passwordChange := entity.PasswordChange{
		ID: id,
	}

	err := u.userRepo.GetPasswordChangeEntry(&passwordChange, dto.ResetPasswordWithID{ID: id})

	return passwordChange.UserID, err
}

func (u *UserUseCase) GetUserID(getUserID dto.ResetPassword) (uuid.UUID, error) {
	user := entity.User{
		Email: getUserID.Email,
	}

	err := u.userRepo.GetUserID(&user, getUserID)

	return user.ID, err
}

func (u *UserUseCase) SoftDelete(userID uuid.UUID) error {
	user := entity.User{
		ID: userID,
	}

	err := u.userRepo.SoftDelete(&user)

	return err
}
