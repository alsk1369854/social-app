package pkg

import (
	"regexp"
	"sync"
)

type AuthUtils struct{}

var BEARER_TOKEN_REGEX = regexp.MustCompile(`^Bearer\s+(?P<tokenString>.+)$`)

var authUtilsOnce sync.Once
var authUtils *AuthUtils

func NewAuthUtils() *AuthUtils {
	authUtilsOnce.Do(func() {
		authUtils = &AuthUtils{}
	})
	return authUtils
}

func (u *AuthUtils) ExtractBearerToken(authHeader string) (string, bool) {
	matches := BEARER_TOKEN_REGEX.FindStringSubmatch(authHeader)
	if len(matches) < 2 {
		return "", false
	}
	tokenString := matches[1]
	if tokenString == "" {
		return "", false
	}
	return tokenString, true
}
