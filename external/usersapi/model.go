package usersapi

import (
	"fmt"

	"github.com/brandonbraner/maas/pkg/permissions"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type User struct {
	ID          primitive.ObjectID      `json:"_id" bson:"_id,omitempty"`
	Username    string                  `json:"username" bson:"username" validate:"email"`
	Password    string                  `json:"password" bson:"password"`
	Firstname   string                  `json:"firstname" bson:"firstname"`
	Lastname    string                  `json:"lastname" bson:"lastname"`
	Permissions permissions.Permissions `json:"permissions" bson:"permissions"`
	Tokens      int                     `json:"tokens" bson:"tokens"`
}

// NewUser is a constructor function for creating a new User instance
func NewUserModel(username, password, firstname, lastname string, tokens int, permissions permissions.Permissions) (*User, error) {

	user := &User{
		Username:    username,
		Password:    password,
		Firstname:   firstname,
		Lastname:    lastname,
		Permissions: permissions,
		Tokens:      tokens,
	}

	err := validate.Struct(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

type UserPermission func(*User) error

// WithPermissions returns a UserO that sets the specified permission to true
func WithPermission(permission string) UserPermission {
	return func(u *User) error {
		switch permission {
		case "generate_llm_meme":
			u.Permissions.GenerateLlmMeme = true
		default:
			return fmt.Errorf("unknown permission: %s", permission)
		}
		return nil
	}
}
