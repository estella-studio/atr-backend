package repository

import (
	"errors"
	"log"
	"time"

	"github.com/estella-studio/leon-backend/internal/domain/dto"
	"github.com/estella-studio/leon-backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMySQLItf interface {
	Register(user *entity.User) error
	RegisterUserDetail(userDetail *entity.UserDetail) error
	Login(user *entity.User) error
	CheckFriendRequestExist(friendRequest *entity.FriendRequest) error
	CheckFriendRequestFromFriend(friendRequest *entity.FriendRequest) error
	NewFriendRequest(friendRequest *entity.FriendRequest) error
	GetFriendRequestSent(friendRequest *[]entity.FriendRequest, userParam dto.GetFriendRequest) error
	GetFriendRequestSentPaged(friendRequest *[]entity.FriendRequest, userParam dto.GetFriendRequest, offset int, limit int) error
	AcceptFriendRequest(friendRequest *entity.FriendRequest) error
	GetFriendList(user *[]entity.User, userID uuid.UUID) error
	AddFriendList(friend *entity.Friend) error
	CheckUserID(checkUserID *entity.User) error
	ChangePassword(user *entity.User) error
	NewEmailVerification(verification *entity.Verification) error
	ValidateEmail(verification *entity.Verification) error
	GetEmailVerification(verification *entity.Verification) error
	GetEmail(user *entity.User, userParam dto.ResetPassword) error
	CheckUsername(user *entity.User) error
	GetUserInfo(user *entity.User) error
	GetUserInfoPublic(user *entity.User) error
	GetUserIDFromEmail(user *entity.User, userParam dto.ResetPassword) error
	GetUserIDFromUsername(user *entity.User) error
	GetUsername(user *entity.User, userParam dto.Login) error
	GetPasswordChangeID(passwordChange *entity.PasswordChange, userParam dto.ResetPassword) error
	GetPasswordResetCode(passwordResetcode *entity.PasswordResetCode) error
	GetPasswordChangeValidity(passwordChange *entity.PasswordChange) error
	GetPasswordResetCodeValidity(passwordChange *entity.PasswordResetCode) error
	GetPasswordChangeEntry(passwordChange *entity.PasswordChange, userParam dto.ResetPasswordWithID) error
	CreatePasswordChangeEntry(passwordChange *entity.PasswordChange) error
	CreatePasswordResetCode(passwordResetCode *entity.PasswordResetCode) error
	UpdatePasswordResetCode(passwordResetCode *entity.PasswordResetCode) error
	UpdatePasswordChangeEntry(passwordChange *entity.PasswordChange) error
	UpdateUserInfo(user *entity.User) error
	UdpateUserDetail(userDetail *entity.UserDetail) error
	UpdateLastActivity(userID uuid.UUID) error
	CheckReportUser(userReporting *entity.UserReporting) error
	ReportUser(userReporting *entity.UserReporting) error
	SoftDelete(user *entity.User) error
}

type UserMySQL struct {
	db *gorm.DB
}

func NewUserMySQL(db *gorm.DB) UserMySQLItf {
	return &UserMySQL{
		db: db,
	}
}

func (r *UserMySQL) Register(user *entity.User) error {
	return r.db.Debug().
		Create(user).
		Error
}

func (r *UserMySQL) RegisterUserDetail(userDetail *entity.UserDetail) error {
	return r.db.Debug().
		Create(userDetail).
		Error
}

func (r *UserMySQL) Login(user *entity.User) error {
	return r.db.Debug().
		First(user).
		Error
}

func (r *UserMySQL) CheckFriendRequestExist(friendRequest *entity.FriendRequest) error {
	return r.db.Debug().
		Select("accepted").
		Where("user_id = ?", friendRequest.UserID).
		Where("friend_id = ?", friendRequest.FriendID).
		First(friendRequest).
		Error
}

func (r *UserMySQL) CheckFriendRequestFromFriend(friendRequest *entity.FriendRequest) error {
	return r.db.Debug().
		Select("accepted").
		Where("user_id = ?", friendRequest.FriendID).
		First(friendRequest).
		Error
}

func (r *UserMySQL) NewFriendRequest(friendRequest *entity.FriendRequest) error {
	return r.db.Debug().
		Create(friendRequest).
		Error
}

func (r *UserMySQL) GetFriendRequestSent(
	friendRequest *[]entity.FriendRequest, userParam dto.GetFriendRequest,
) error {
	return r.db.Debug().
		Select("user_id, friend_id").
		Where("accepted = ? ", false).
		Find(friendRequest, userParam).
		Error
}

func (r *UserMySQL) GetFriendRequestSentPaged(
	friendRequest *[]entity.FriendRequest, userParam dto.GetFriendRequest,
	offset int, limit int,
) error {
	return r.db.Debug().
		Select("user_id, friend_id").
		Where("accepted = ? ", false).
		Limit(limit).
		Offset(offset).
		Find(friendRequest, userParam).
		Error
}

func (r *UserMySQL) AcceptFriendRequest(friendRequest *entity.FriendRequest) error {
	if r.db.Debug().
		Model(&friendRequest).
		Where("user_id = ?", friendRequest.UserID).
		Where("friend_id = ?", friendRequest.FriendID).
		Where("accepted = ?", false).
		Update("accepted", true).RowsAffected == 0 {
		return errors.New("invalid friend id")
	}

	return nil
}

func (r *UserMySQL) GetFriendList(user *[]entity.User, userID uuid.UUID) error {
	return r.db.Debug().
		Model(&user).
		Raw(`
		SELECT username, name FROM users WHERE id IN
		(
			SELECT 
				CASE 
					WHEN user_id = ? THEN friend_id
					ELSE user_id
				END AS user_id
			FROM friends
			WHERE user_id = ? OR friend_id = ?
		)
		`,
			userID, userID, userID,
		).
		Scan(&user).
		Error
}

func (r *UserMySQL) AddFriendList(friend *entity.Friend) error {
	return r.db.Debug().
		Create(friend).
		Error
}

func (r *UserMySQL) CheckUserID(user *entity.User) error {
	return r.db.Debug().
		First(user).
		Error
}

func (r *UserMySQL) ChangePassword(user *entity.User) error {
	return r.db.Debug().
		Model(&user).
		Update("password", user.Password).
		Error
}

func (r *UserMySQL) NewEmailVerification(verification *entity.Verification) error {
	return r.db.Debug().
		Create(verification).
		Error
}

func (r *UserMySQL) ValidateEmail(verification *entity.Verification) error {
	return r.db.Debug().
		Model(&verification).
		Where(
			entity.Verification{
				Email: verification.Email,
				Code:  verification.Code,
			},
		).
		Update("success", true).
		Error
}

func (r *UserMySQL) GetEmailVerification(verification *entity.Verification) error {
	return r.db.Debug().
		Select("code, success").
		Where("email = ?", verification.Email).
		First(verification).
		Error
}

func (r *UserMySQL) GetEmail(user *entity.User, userParam dto.ResetPassword) error {
	return r.db.Debug().
		Select("email").
		First(user, userParam).
		Error
}

func (r *UserMySQL) CheckUsername(user *entity.User) error {
	return r.db.Debug().
		Raw("SELECT `username` FROM `users` WHERE username = ?", user.Username).
		First(&user).
		Error
}

func (r *UserMySQL) GetUserInfo(user *entity.User) error {
	return r.db.Debug().
		Model(&user).
		Preload("UserDetail").
		Select("users.id, users.email, users.username, users.name, users.created_at, users.updated_at, user_details.profile_index, user_details.last_activity").
		Joins("LEFT JOIN user_details ON user_details.user_id = users.id").
		First(&user).
		Error
}

func (r *UserMySQL) GetUserInfoPublic(user *entity.User) error {
	return r.db.Debug().
		Model(&user).
		Preload("UserDetail").
		Select("users.username, users.name, user_details.profile_index, user_details.last_activity").
		Joins("LEFT JOIN user_details ON user_details.user_id = users.id").
		First(&user).
		Error
}

func (r *UserMySQL) GetUserIDFromEmail(user *entity.User, userParam dto.ResetPassword) error {
	return r.db.Debug().
		Select("id").
		First(user, userParam).
		Error
}

func (r *UserMySQL) GetUserIDFromUsername(user *entity.User) error {
	return r.db.Debug().
		Select("id").
		Where("username = ?", user.Username).
		First(user).
		Error
}

func (r *UserMySQL) GetUsername(user *entity.User, userParam dto.Login) error {
	return r.db.Debug().
		First(user, userParam).
		Error
}

func (r *UserMySQL) GetPasswordChangeID(passwordChange *entity.PasswordChange, userParam dto.ResetPassword) error {
	return r.db.Debug().
		Select("id").
		First(passwordChange, userParam).
		Error
}

func (r *UserMySQL) GetPasswordResetCode(passwordResetcode *entity.PasswordResetCode) error {
	return r.db.Debug().
		Order("created_at desc").
		Select("password_change_id, code").
		Where("user_id = ?", passwordResetcode.UserID).
		First(passwordResetcode).
		Error
}

func (r *UserMySQL) GetPasswordResetCodeValidity(passwordChange *entity.PasswordResetCode) error {
	return r.db.Debug().
		Order("created_at desc").
		Find(&passwordChange).
		Select("created_at").
		Where("user_id = ?", passwordChange.UserID).
		Error
}

func (r *UserMySQL) GetPasswordChangeValidity(passwordChange *entity.PasswordChange) error {
	return r.db.Debug().
		Select("id, created_at, success").
		First(passwordChange).
		Error
}

func (r *UserMySQL) GetPasswordChangeEntry(passwordChange *entity.PasswordChange, userParam dto.ResetPasswordWithID) error {
	return r.db.Debug().
		Select("user_id").
		First(passwordChange).
		Error
}

func (r *UserMySQL) CreatePasswordChangeEntry(passwordChange *entity.PasswordChange) error {
	return r.db.Debug().
		Create(passwordChange).
		Error
}

func (r *UserMySQL) CreatePasswordResetCode(passwordResetCode *entity.PasswordResetCode) error {
	return r.db.Debug().
		Create(passwordResetCode).
		Error
}

func (r *UserMySQL) UpdatePasswordResetCode(passwordResetCode *entity.PasswordResetCode) error {
	return r.db.Debug().
		Model(&passwordResetCode).
		Where("password_change_id = ?", passwordResetCode.PasswordChangeID).
		Delete(passwordResetCode).
		Error
}

func (r *UserMySQL) UpdatePasswordChangeEntry(passwordChange *entity.PasswordChange) error {
	return r.db.Debug().
		Model(&passwordChange).
		Where("id = ?", passwordChange.ID).
		Update("success", passwordChange.Success).
		Error
}

func (r *UserMySQL) UpdateUserInfo(user *entity.User) error {
	return r.db.Debug().
		Updates(user).
		Error
}

func (r *UserMySQL) UdpateUserDetail(userDetail *entity.UserDetail) error {
	var err error

	if r.db.Debug().
		Where("user_id = ?", userDetail.UserID).
		Updates(userDetail).RowsAffected == 0 {
		userDetail.ID = uuid.New()

		err = r.db.Create(userDetail).Error
	}

	return err
}

func (r *UserMySQL) UpdateLastActivity(userID uuid.UUID) error {
	userDetail := entity.UserDetail{
		LastActivity: time.Now().UTC(),
	}
	var err error

	log.Println(userID)

	if r.db.Debug().
		Table("user_details").
		Where("user_id = ?", userID).
		Update("last_activity", userDetail.LastActivity).RowsAffected == 0 {
		userDetail.ID = uuid.New()
		userDetail.UserID = userID

		err = r.db.Create(&userDetail).Error
	}

	return err
}

func (r *UserMySQL) CheckReportUser(userReporting *entity.UserReporting) error {
	return r.db.Debug().
		Select("id").
		Where("user_id = ?", userReporting.UserID).
		Where("reporter_id = ?", userReporting.ReporterID).
		Take(userReporting).
		Error
}

func (r *UserMySQL) ReportUser(userReporting *entity.UserReporting) error {
	return r.db.Debug().
		Create(userReporting).
		Error
}

func (r *UserMySQL) SoftDelete(user *entity.User) error {
	return r.db.Debug().
		Delete(user).
		Error
}
