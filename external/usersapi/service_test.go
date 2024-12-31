package usersapi

import (
	"testing"

	"github.com/brandonbraner/maas/config"
	"github.com/brandonbraner/maas/pkg/jwt"
	"github.com/brandonbraner/maas/pkg/permissions"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GenerateJwt(t *testing.T) {
	// Setup
	config.SetupMongoTestConfig()
	service, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Run("successful token generation", func(t *testing.T) {
		username := "testuser@example.com"

		usermodel, err := service.NewUser(username, "password", "Bob", "Smith", 100, permissions.Permissions{})
		if err != nil {
			t.Fatal(err.Error())
		}

		user, err := service.CreateUser(usermodel)
		if err != nil {
			t.Fatal(err.Error())
		}

		token, err := service.GenerateJwt(user.Username)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify token claims
		claims, err := jwt.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, user.ID.Hex(), claims.UserID)
		assert.Equal(t, username, claims.Email)

	})

	// t.Run("empty user ID", func(t *testing.T) {
	// 	_, err := service.GenerateJwt("", "testuser", permissions.Permissions{})
	// 	assert.Error(t, err)
	// })

	// t.Run("empty username", func(t *testing.T) {
	// 	_, err := service.GenerateJwt("123", "", permissions.Permissions{})
	// 	assert.Error(t, err)
	// })

}
