package accounts

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	UUID      string `json:"uuid" bson:"_id,omitempty"`
	Email     string `json:"email" bson:"email,omitempty"`
	Password  string `json:"-" bson:"password,omitempty"`
	AvatarURL string `json:"avatarURL" bson:"avatar,omitempty"`
	Username  string `json:"username" bson:"username,omitempty"`
	Sex       string `json:"sex" bson:"sex,omitempty"`
	Country   string `json:"country" bson:"country,omitempty"`
	Language  string `json:"language" bson:"lang,omitempty"`
	Birthday  int64  `json:"birthday" bson:"birthday,omitempty"`
	CreatedAt int64  `json:"-" bson:"created_at,omitempty"`
	LoginAt   int64  `json:"-" bson:"login_at,omitempty"`
	LogoutAt  int64  `json:"-" bson:"logout_at,omitempty"`
	IsActive  bool   `json:"-" bson:"is_active,omitempty"`
	IsAdmin   bool   `json:"-" bson:"is_admin,omitempty"`
	IsDeleted bool   `json:"-" bson:"is_deleted,omitempty"`
}

func (u *Account) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match")
	}
	return nil
}

func (u *Account) GeneratePasswordHash() error {
	pwd, err := generatePasswordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}

type CreateAccountDTO struct {
	Email          string `json:"email" bson:"email"`
	Password       string `json:"password" bson:"password"`
	RepeatPassword string `json:"repeat_password" bson:"-"`
}

type CredentialsDTO struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UpdateCredentialsDTO struct {
	UUID        string `json:"uuid,omitempty" bson:"_id,omitempty"`
	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	Password    string `json:"password,omitempty" bson:"password,omitempty"`
	OldPassword string `json:"old_password,omitempty" bson:"-"`
	NewPassword string `json:"new_password,omitempty" bson:"-"`
}

type UpdateAccountDTO struct {
	UUID      string `json:"uuid,omitempty" bson:"_id,omitempty"`
	AvatarURL string `json:"avatarURL,omitempty" bson:"avatar,omitempty"`
	Username  string `json:"username,omitempty" bson:"username,omitempty"`
	Sex       string `json:"sex,omitempty" bson:"sex,omitempty"`
	Country   string `json:"country,omitempty" bson:"country,omitempty"`
	Language  string `json:"language,omitempty" bson:"lang,omitempty"`
	Birthday  int64  `json:"birthday,omitempty" bson:"birthday,omitempty"`
}

func NewAccount(dto CreateAccountDTO) Account {
	tNow := time.Now().UnixNano()
	return Account{
		Email:     dto.Email,
		Password:  dto.Password,
		CreatedAt: tNow,
		IsActive:  true,
		IsAdmin:   false,
	}
}
func NewAdmin(dto CreateAccountDTO) Account {
	tNow := time.Now().UnixNano()
	return Account{
		Email:     dto.Email,
		Password:  dto.Password,
		CreatedAt: tNow,
		IsActive:  true,
		IsAdmin:   true,
	}
}

func UpdatedCredentials(dto UpdateCredentialsDTO) Account {
	return Account{
		UUID:     dto.UUID,
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func UpdatedAccount(dto UpdateAccountDTO) Account {
	return Account{
		UUID:      dto.UUID,
		AvatarURL: dto.AvatarURL,
		Username:  dto.Username,
		Sex:       dto.Sex,
		Country:   dto.Country,
		Language:  dto.Language,
		Birthday:  dto.Birthday,
	}
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password due to error %w", err)
	}
	return string(hash), nil
}
