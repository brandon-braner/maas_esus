package usersapi

import (
	"context"
	"errors"

	"github.com/brandonbraner/maas/pkg/jwt"
	"github.com/brandonbraner/maas/pkg/permissions"
)

type UserService struct {
	Repo *UserRepository
}

func NewUserService() (*UserRepository, error) {
	repo, err := NewUserRepository(context.TODO())
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

func (us *UserService) GenerateJwt(userid, username string, permissions permissions.Permissions) (string, error) {
	if userid == "" {
		return "", errors.New("user ID cannot be empty")
	}
	if username == "" {
		return "", errors.New("username cannot be empty")
	}

	token, err := jwt.GenerateToken(userid, username, permissions)
	if err != nil {
		return "", err
	}

	return token, nil
}
