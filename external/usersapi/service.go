package usersapi

import (
	"context"
	"errors"

	"github.com/brandonbraner/maas/pkg/jwt"
	"github.com/brandonbraner/maas/pkg/permissions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	Repo *userRepository
}

func NewUserService() (*UserService, error) {
	repo, err := NewUserRepository(context.TODO())
	if err != nil {
		return nil, err
	}

	service := UserService{
		Repo: repo,
	}
	return &service, nil
}

func (us *UserService) NewUser(username, password, firstname, lastname string, tokens int, permissions permissions.Permissions) (*User, error) {
	user, err := NewUserModel(username, password, firstname, lastname, tokens, permissions)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) CreateUser(user *User) (*User, error) {
	insertid, err := us.Repo.Create(context.TODO(), user)

	if err != nil {
		return nil, err
	}
	user.ID = insertid.InsertedID.(primitive.ObjectID)
	return user, err
}

func (us *UserService) DeleteUser(userid string) {
	us.Repo.Delete(context.TODO(), userid)
}

func (us *UserService) GenerateJwt(username string) (string, error) {
	if username == "" {
		return "", errors.New("username cannot be empty")
	}

	// use the repository so we can get everything we need from the database
	// it will also run it through validation
	user, err := us.Repo.GetByUserName(context.TODO(), username)

	if err != nil {
		return "", err
	}
	userid := user.ID.Hex()
	token, err := jwt.GenerateToken(userid, user.Username, user.Permissions)
	if err != nil {
		return "", err
	}

	return token, nil
}
