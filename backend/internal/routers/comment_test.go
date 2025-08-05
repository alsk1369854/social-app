package routers

import (
	"backend/internal/tests"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommentRouter(t *testing.T) {
	server, apiRouter, _, _, cleanup := tests.SetupTestServer("test_comment_router")
	defer cleanup()

	NewCommentRouter().Bind(apiRouter)
	NewPostRouter().Bind(apiRouter)
	NewUserRouter().Bind(apiRouter)

	userData, accessToken, err := tests.SetupTestUser(server)
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	assert.NotEmpty(t, accessToken)

}
