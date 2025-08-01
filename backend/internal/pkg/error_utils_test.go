package pkg

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorUtils(t *testing.T) {
	errorUtils := NewErrorUtils()

	t.Run("Server Internal Error", func(t *testing.T) {
		err := errorUtils.ServerInternalError("Test error")
		assert.Error(t, err, "Expected an error to be returned")
		assert.Contains(t, err.Error(), "Internal Server Error: Test error", "Error message should contain the provided message")
	})

	t.Run("Is Server Internal Error", func(t *testing.T) {
		err := errorUtils.ServerInternalError("Test error")
		assert.True(t, errorUtils.IsServerInternalError(err.Error()), "Expected the error to be recognized as a server internal error")

		otherErr := errors.New("Some other error")
		assert.False(t, errorUtils.IsServerInternalError(otherErr.Error()), "Expected the error to not be recognized as a server internal error")
	})

}
