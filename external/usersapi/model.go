package usersapi

import (
	"github.com/brandonbraner/maas/pkg/permissions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID      `bson:"_id,omitempty"`
	Username    string                  `json:"username" bson:"username"`
	Password    string                  `json:"password" bson:"password"`
	Firstname   string                  `json:"firstname" bson:"firstname"`
	Lastname    string                  `json:"lastname" bson:"lastname"`
	Permissions permissions.Permissions `json:"permissions" bson:"permissions"`
	Tokens      int                     `json:"tokens" bson:"tokens"`
}

// NewUser is a constructor function for creating a new User instance
func NewUser(username, password, firstname, lastname string, perms permissions.Permissions, tokens int) *User {
	return &User{
		ID:          primitive.NewObjectID(), // Automatically generate a new ObjectID
		Username:    username,
		Password:    password,
		Firstname:   firstname,
		Lastname:    lastname,
		Permissions: perms,
		Tokens:      tokens,
	}
}
