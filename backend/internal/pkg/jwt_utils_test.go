package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTUtils(t *testing.T) {
	TEST_SECRET := []byte("test_secret")
	jwtUtils := NewJWTUtils()

	t.Run("Generate Token", func(t *testing.T) {
		data := map[string]any{"user": "test_user"}
		token, err := jwtUtils.GenerateToken(data, TEST_SECRET)
		assert.NoError(t, err, "Expected no error while generating token")
		assert.NotEmpty(t, token, "Generated token should not be empty")
	})

	t.Run("Parse Token", func(t *testing.T) {
		data := map[string]any{"user": "test_user"}
		token, err := jwtUtils.GenerateToken(data, TEST_SECRET)
		assert.NoError(t, err, "Expected no error while generating token")

		claims, err := jwtUtils.ParseToken(token, TEST_SECRET)
		assert.NoError(t, err, "Expected no error while parsing token")
		assert.Equal(t, data["user"], claims["data"].(map[string]any)["user"], "Parsed data should match original data")
	})

}
