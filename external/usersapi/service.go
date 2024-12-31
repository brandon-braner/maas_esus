package usersapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/brandonbraner/maas/pkg/jwtservice"
	"github.com/brandonbraner/maas/pkg/permissions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	Repo *userRepository
	mu   sync.Mutex
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

func (us *UserService) NewUser(username, password, firstname, lastname string, tokens int, perms permissions.Permissions) (*User, error) {
	user, err := NewUserModel(username, password, firstname, lastname, tokens, perms)
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
	return user, nil
}

func (us *UserService) DeleteUser(userid string) {
	us.Repo.Delete(context.TODO(), userid)
}

func (us *UserService) GenerateJwt(username string) (string, error) {
	user, err := us.Repo.GetByUserName(context.TODO(), username)
	if err != nil {
		return "", err
	}
	userid := user.ID.Hex()
	token, err := jwtservice.GenerateToken(userid, user.Username, user.Permissions)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *UserService) GetUserByUsername(username string) (*User, error) {
	user, err := us.Repo.GetByUserName(context.TODO(), username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) UpdatePermission(username string, permission string, val bool) error {

	ok := permissions.ValidatePermission(permission)
	if !ok {
		return fmt.Errorf("invalid permission %s", permission)
	}

	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"permissions." + permission: val}}

	_, err := us.Repo.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("error updating permission for user %s: %v", username, err)
		return fmt.Errorf("error updating permission for user %s: %v", username, err)
	}
	return nil
}

func (us *UserService) DeleteAllUsers() (int64, error) {
	count, err := us.Repo.DeleteAll(context.TODO())
	return count, err
}

func (us *UserService) UpdateTokens(username string, amount int) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	user, err := us.Repo.GetByUserName(context.TODO(), username)
	if err != nil {
		return err
	}

	user.Tokens += amount
	err = us.Repo.Update(context.TODO(), user.ID.Hex(), user)
	return err
}

func (us *UserService) GetTokenCount(username string) (int, error) {
	if username == "" {
		return 0, errors.New("username cannot be empty")
	}

	user, err := us.Repo.GetByUserName(context.TODO(), username)
	if err != nil {
		return 0, err
	}

	return user.Tokens, nil

}
