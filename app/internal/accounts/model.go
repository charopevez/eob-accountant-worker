package accounts

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	UUID     string `json:"uuid" bson:"_id,omitempty"`
	Email    string `json:"email" bson:"email,omitempty"`
	Password string `json:"-" bson:"password,omitempty"`
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

type UpdateCredentialsDTO struct {
	UUID        string `json:"uuid,omitempty" bson:"_id,omitempty"`
	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	Password    string `json:"password,omitempty" bson:"password,omitempty"`
	OldPassword string `json:"old_password,omitempty" bson:"-"`
	NewPassword string `json:"new_password,omitempty" bson:"-"`
}
type UpdateBioDTO struct {
	UUID        string `json:"uuid,omitempty" bson:"_id,omitempty"`
	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	Password    string `json:"password,omitempty" bson:"password,omitempty"`
	OldPassword string `json:"old_password,omitempty" bson:"-"`
	NewPassword string `json:"new_password,omitempty" bson:"-"`
}

func NewAccount(dto CreateAccountDTO) Account {
	return Account{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func UpdatedCredentials(dto UpdateCredentialsDTO) Account {
	return Account{
		UUID:     dto.UUID,
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password due to error %w", err)
	}
	return string(hash), nil
}
