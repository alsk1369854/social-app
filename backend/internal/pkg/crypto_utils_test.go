package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCryptoUtils(t *testing.T) {
	// Test singleton instance creation
	t.Run("Singleton Instance Creation", func(t *testing.T) {
		instance1 := NewCryptoUtils()
		instance2 := NewCryptoUtils()

		assert.Same(t, instance1, instance2, "Expected the same instance to be returned")
		assert.NotNil(t, instance1, "Instance should not be nil")
	})
}

func TestHash256Password(t *testing.T) {
	cryptoUtils := NewCryptoUtils()

	t.Run("Hashing Password", func(t *testing.T) {
		username := "testuser"
		email := "test@example.com"
		password := "password123"
		hashedPassword := cryptoUtils.Hash256Password(email, username, password)

		assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
	})

	t.Run("Verifying Password", func(t *testing.T) {
		username := "testuser"
		email := "test@example.com"
		password := "password123"
		hashedPassword := cryptoUtils.Hash256Password(email, username, password)

		assert.True(t, cryptoUtils.VerifyPassword(hashedPassword, email, username, password), "Expected password to be verified successfully")
	})

	t.Run("Verifying Incorrect Password", func(t *testing.T) {
		username := "testuser"
		email := "test@example.com"
		password := "password123"
		hashedPassword := cryptoUtils.Hash256Password(email, username, password)

		assert.False(t, cryptoUtils.VerifyPassword(hashedPassword, email, username, "wrongpassword"), "Expected password verification to fail")
	})
}
