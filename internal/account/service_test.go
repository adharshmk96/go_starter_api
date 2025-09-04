package account_test

import (
	"servicehub_api/internal/account"
	"servicehub_api/pkg/domain"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAccountService_HashPassword(t *testing.T) {
	t.Run("should hash and compare password correctly", func(t *testing.T) {
		service := account.NewAccountService()

		password := "password"
		hash, err := service.HashPassword(password)
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}

		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		ok, err := service.ComparePassword(password, hash)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("should return error if password is empty", func(t *testing.T) {
		service := account.NewAccountService()

		password := ""
		hash, err := service.HashPassword(password)
		assert.ErrorIs(t, err, domain.ErrPasswordEmpty)
		assert.Empty(t, hash)
	})
}

func TestAccountService_GenerateAndValidateToken(t *testing.T) {
	// Set up test environment
	viper.Set("JWT_SECRET", "test_secret_key_for_jwt_validation")
	defer viper.Reset()

	service := account.NewAccountService()

	t.Run("should generate and validate token correctly", func(t *testing.T) {
		account := &domain.Account{ID: 123, Email: "test@example.com"}

		// Generate token
		token, err := service.GenerateToken(account)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate token
		accountID, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, uint(123), accountID)
	})

	t.Run("should return error if JWT secret is not set", func(t *testing.T) {
		// Temporarily unset JWT secret
		viper.Set("JWT_SECRET", "")

		account := &domain.Account{ID: 1, Email: "test@test.com"}
		token, err := service.GenerateToken(account)
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("should return error if token is invalid", func(t *testing.T) {
		invalidToken := "invalid_token"
		accountID, err := service.ValidateToken(invalidToken)
		assert.Error(t, err)
		assert.Equal(t, uint(0), accountID)
	})

	t.Run("should return error if token is malformed", func(t *testing.T) {
		malformedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid"
		accountID, err := service.ValidateToken(malformedToken)
		assert.Error(t, err)
		assert.Equal(t, uint(0), accountID)
	})
}
