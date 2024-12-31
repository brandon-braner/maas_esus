package usersapi

import (
	"testing"

	"github.com/brandonbraner/maas/config"
	"github.com/brandonbraner/maas/pkg/jwtservice"
	"github.com/brandonbraner/maas/pkg/permissions"
	"github.com/stretchr/testify/assert"
)

func Setup(t *testing.T, service *UserService) *User {
	username := "testuser@example.com"

	usermodel, err := service.NewUser(username, "password", "Bob", "Smith", 100, permissions.Permissions{})
	if err != nil {
		t.Fatal(err.Error())
	}

	user, err := service.CreateUser(usermodel)
	if err != nil {
		t.Fatal(err.Error())
	}

	return user
}

func TearDown(service *UserService, userid string) {
	service.DeleteUser(userid)
}

func TestUserService_GenerateJwt(t *testing.T) {
	// Setup
	config.SetupMongoTestConfig()
	service, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}
	user := Setup(t, service)
	t.Run("successful token generation", func(t *testing.T) {

		token, err := service.GenerateJwt(user.Username)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		useridString := user.ID.Hex()
		// Verify token claims
		claims, err := jwtservice.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, useridString, claims.UserID)
		assert.Equal(t, user.Username, claims.Email)

		TearDown(service, useridString)
	})

	t.Run("empty username", func(t *testing.T) {
		_, err := service.GenerateJwt("")
		assert.Error(t, err)
	})

}
