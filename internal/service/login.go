package service

import (
	"context"
	"github.com/masihyeganeh/audit-log/internal/auth"
	"github.com/pkg/errors"
)

// Login user login
func (s *service) Login(ctx context.Context, request LoginRequest) (string, error) {
	storedHashedPassword, salt, hasReadAccess, hasWriteAccess, err := s.datastoreRepository.FindUser(ctx, request.Username)
	if err != nil {
		return "", err
	}

	generatedHashedPassword, err := s.authentication.EncryptPasswordWithSalt(request.Password, salt)
	if err != nil {
		return "", err
	}

	// TODO: is vulnerable to timing attack. user constant-time comparison instead
	if storedHashedPassword != generatedHashedPassword {
		return "", errors.New("wrong password")
	}

	return s.authentication.Login(auth.User{
		Username:       request.Username,
		HasReadAccess:  hasReadAccess,
		HasWriteAccess: hasWriteAccess,
	}), nil
}
