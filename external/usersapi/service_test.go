package usersapi

import (
	"testing"

	"github.com/brandonbraner/maas/pkg/jwt"
	"github.com/brandonbraner/maas/pkg/permissions"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GenerateJwt(t *testing.T) {
	// Setup
	service := &UserService{}

	t.Run("successful token generation", func(t *testing.T) {
		userID := "123"
		username := "testuser"
		perms := permissions.Permissions{}

		token, err := service.GenerateJwt(userID, username, perms)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify token claims
		claims, err := jwt.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, username, claims.Email)
	})

	t.Run("empty user ID", func(t *testing.T) {
		_, err := service.GenerateJwt("", "testuser", permissions.Permissions{})
		assert.Error(t, err)
	})

	t.Run("empty username", func(t *testing.T) {
		_, err := service.GenerateJwt("123", "", permissions.Permissions{})
		assert.Error(t, err)
	})

}
