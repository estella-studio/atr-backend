package usecase

import (
	"errors"
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
	RenewToken(renewToken dto.RenewToken) (string, error)
	CheckUserID(checkUserID *dto.CheckUserID) error
	CheckFriendRequestExist(CheckFriendRequestExist *dto.CheckFriendRequestExist) (bool, error)
	CheckFriendRequestFromFriend(friendID uuid.UUID) (bool, error)
	NewFriendRequest(sendFriendRequest *dto.SendFriendRequest) error
	GetFriendRequestSent(userID uuid.UUID, offset int, limit int) (*[]dto.ResponseGetFriendRequest, error)
	AcceptFriendRequest(acceptFriendRequest *dto.AcceptFriendRequest) error
	GetFriendList(userID uuid.UUID) (*[]dto.ResponseFriendList, error)
	NewEmailVerification(emailVerification *dto.EmailVerification) error
	ValidateEmail(validateEmail *dto.ValidateEmail) error
	GetEmailVerification(validateEmail *dto.EmailVerification) (uint, bool, error)
	CheckUsername(userName *dto.CheckUsername) error
	GetUserInfo(userID uuid.UUID) (dto.ResponseGetUserInfo, error)
	GetUserInfoPublic(userID uuid.UUID) (dto.ResponseGetUserInfoPublic, error)
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
	GetUserIDFromEmail(getUserID dto.ResetPassword) (uuid.UUID, error)
	GetUserIDFromUsername(username string) (uuid.UUID, error)
	ReportUser(reportUser dto.ReportUser) error
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

	return user.ParseToDTOResponseLogin(), token, nil
}

func (u *UserUseCase) RenewToken(renewToken dto.RenewToken) (string, error) {
	var user entity.User

	err := u.userRepo.CheckUserID(&user)
	if err != nil {
		return "",
			err
	}

	token, err := u.jwt.GenerateToken(user.ID)

	return token, err
}

func (u *UserUseCase) CheckUserID(checkUserID *dto.CheckUserID) error {
	user := entity.User{
		ID: checkUserID.ID,
	}

	err := u.userRepo.CheckUserID(&user)

	return err
}

func (u *UserUseCase) CheckFriendRequestExist(checkFriendRequestExist *dto.CheckFriendRequestExist) (bool, error) {
	friendRequest := entity.FriendRequest{
		UserID:   checkFriendRequestExist.UserID,
		FriendID: checkFriendRequestExist.FriendID,
	}

	err := u.userRepo.CheckFriendRequestExist(&friendRequest)

	return friendRequest.Accepted, err
}

func (u *UserUseCase) CheckFriendRequestFromFriend(friendID uuid.UUID) (bool, error) {
	friendRequest := entity.FriendRequest{
		FriendID: friendID,
	}

	err := u.userRepo.CheckFriendRequestFromFriend(&friendRequest)

	return friendRequest.Accepted, err
}

func (u *UserUseCase) NewFriendRequest(sendFriendRequest *dto.SendFriendRequest) error {
	friendRequest := entity.FriendRequest{
		ID:       uuid.New(),
		UserID:   sendFriendRequest.UserID,
		FriendID: sendFriendRequest.FriendID,
		Accepted: false,
	}

	err := u.userRepo.NewFriendRequest(&friendRequest)

	return err
}

func (u *UserUseCase) GetFriendRequestSent(userID uuid.UUID, offset int, limit int) (*[]dto.ResponseGetFriendRequest, error) {
	friendRequest := new([]entity.FriendRequest)

	if offset == 0 && limit == 0 {
		err := u.userRepo.GetFriendRequestSent(friendRequest, dto.GetFriendRequest{UserID: userID})
		if err != nil {
			return nil, err
		}
	} else {
		err := u.userRepo.GetFriendRequestSentPaged(friendRequest, dto.GetFriendRequest{UserID: userID}, offset, limit)
		if err != nil {
			return nil, err
		}
	}

	res := make([]dto.ResponseGetFriendRequest, len(*friendRequest))

	for i, friendRequest := range *friendRequest {
		res[i] = friendRequest.ParseToDTOResponseGetFriendRequest()
	}

	return &res, nil
}

func (u *UserUseCase) AcceptFriendRequest(acceptFriendRequest *dto.AcceptFriendRequest) error {
	friendRequest := entity.FriendRequest{
		UserID:   acceptFriendRequest.FriendID,
		FriendID: acceptFriendRequest.UserID,
	}

	friend := entity.Friend{
		UserID:   acceptFriendRequest.UserID,
		FriendID: acceptFriendRequest.FriendID,
	}

	err := u.userRepo.AcceptFriendRequest(&friendRequest)
	if err != nil {
		log.Println(err)
		return err
	}

	err = u.userRepo.AddFriendList(&friend)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (u *UserUseCase) GetFriendList(userID uuid.UUID) (*[]dto.ResponseFriendList, error) {
	user := new([]entity.User)

	err := u.userRepo.GetFriendList(user, userID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ResponseFriendList, len(*user))

	for i, friend := range *user {
		res[i] = friend.ParseToDTOResponseGetFriendList()
	}

	return &res, nil
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

func (u *UserUseCase) GetUserInfoPublic(userID uuid.UUID) (dto.ResponseGetUserInfoPublic, error) {
	user := entity.User{
		ID: userID,
	}

	err := u.userRepo.GetUserInfoPublic(&user)
	if err != nil {
		return dto.ResponseGetUserInfoPublic{},
			err
	}

	return user.ParseToDTOResponseGetUserInfoPublic(), nil
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

func (u *UserUseCase) GetUserIDFromEmail(getUserID dto.ResetPassword) (uuid.UUID, error) {
	user := entity.User{
		Email: getUserID.Email,
	}

	err := u.userRepo.GetUserIDFromEmail(&user, getUserID)

	return user.ID, err
}

func (u *UserUseCase) GetUserIDFromUsername(username string) (uuid.UUID, error) {
	user := entity.User{
		Username: username,
	}

	err := u.userRepo.GetUserIDFromUsername(&user)
	if err != nil {
		return uuid.Nil,
			err
	}

	return user.ID, nil
}

func (u *UserUseCase) ReportUser(reportUser dto.ReportUser) error {
	userReporting := entity.UserReporting{
		UserID:     reportUser.UserID,
		ReporterID: reportUser.ReporterID,
	}

	err := u.userRepo.CheckReportUser(&userReporting)
	if err == nil {
		return errors.New("user already reported")
	}

	userReporting.ID = uuid.New()

	err = u.userRepo.ReportUser(&userReporting)

	return err
}

func (u *UserUseCase) SoftDelete(userID uuid.UUID) error {
	user := entity.User{
		ID: userID,
	}

	err := u.userRepo.SoftDelete(&user)

	return err
}
