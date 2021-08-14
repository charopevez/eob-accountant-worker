package accounts

import (
	"context"
	"errors"
	"fmt"

	"github.com/charopevez/eob-accountant-worker/internal/apperror"
	"github.com/charopevez/eob-accountant-worker/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

var _ Service = &service{}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(accountStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: accountStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	Create(ctx context.Context, dto CreateAccountDTO) (string, error)
	AuthenticateAccount(ctx context.Context, email, password string) (Account, error)
	GetAccount(ctx context.Context, uuid string) (Account, error)
	UpdateCredentials(ctx context.Context, dto UpdateCredentialsDTO) error
	Delete(ctx context.Context, uuid string) error
}

//?register new user
func (s service) Create(ctx context.Context, dto CreateAccountDTO) (accUUID string, err error) {
	s.logger.Debug("check password and repeat password")
	if dto.Password != dto.RepeatPassword {
		return accUUID, apperror.BadRequestError("password does not match repeat password")
	}

	acc := NewAccount(dto)

	s.logger.Debug("generate password hash")
	err = acc.GeneratePasswordHash()
	if err != nil {
		s.logger.Errorf("failed to create user account due to error %v", err)
		return
	}

	accUUID, err = s.storage.Create(ctx, acc)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return accUUID, err
		}
		return accUUID, fmt.Errorf("failed to create user. error: %w", err)
	}

	return accUUID, nil
}

//? authenticate user by mail and password
func (s service) AuthenticateAccount(ctx context.Context, email, password string) (u Account, err error) {
	u, err = s.storage.FindByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return u, err
		}
		return u, fmt.Errorf("failed to find user by email. error: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return u, apperror.ErrNotFound
	}

	return u, nil
}

func (s service) GetAccount(ctx context.Context, uuid string) (u Account, err error) {
	u, err = s.storage.FindOne(ctx, uuid)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return u, err
		}
		return u, fmt.Errorf("failed to find user by uuid. error: %w", err)
	}
	return u, nil
}

//? update user credentials
func (s service) UpdateCredentials(ctx context.Context, dto UpdateCredentialsDTO) error {
	var updatedAccount Account
	s.logger.Debug("compare old and new passwords")
	if dto.OldPassword != dto.NewPassword || dto.NewPassword == "" {
		s.logger.Debug("get account by uuid")
		account, err := s.GetAccount(ctx, dto.UUID)
		if err != nil {
			return err
		}

		s.logger.Debug("compare hash current password and old password")
		err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(dto.OldPassword))
		if err != nil {
			return apperror.BadRequestError("old password does not match current password")
		}

		dto.Password = dto.NewPassword
	}

	updatedAccount = UpdatedCredentials(dto)

	s.logger.Debug("generate password hash")
	err := updatedAccount.GeneratePasswordHash()
	if err != nil {
		return fmt.Errorf("failed to update account credentials. error %w", err)
	}

	err = s.storage.UpdateCredentials(ctx, updatedAccount)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update user. error: %w", err)
	}
	return nil
}

func (s service) Delete(ctx context.Context, uuid string) error {
	err := s.storage.Delete(ctx, uuid)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return err
}
