package pkg

import (
	"regexp"
	"sync"

	"github.com/pkg/errors"
)

const SERVER_INTERNAL_ERROR_PREFIX = "Internal Server Error: "

type ErrorUtils struct{}

type ServerInternalError error

var ServerInternalErrorRegex = regexp.MustCompile("^" + regexp.QuoteMeta(SERVER_INTERNAL_ERROR_PREFIX))

var errorUtilsOnce sync.Once
var errorUtils *ErrorUtils

func NewErrorUtils() *ErrorUtils {
	errorUtilsOnce.Do(func() {
		errorUtils = &ErrorUtils{}
	})
	return errorUtils
}

func (u *ErrorUtils) ServerInternalError(msg string) ServerInternalError {
	if u.IsServerInternalError(msg) {
		return errors.New(msg)
	}
	return errors.New(SERVER_INTERNAL_ERROR_PREFIX + msg)
}

func (u *ErrorUtils) IsServerInternalError(msg string) bool {
	return ServerInternalErrorRegex.MatchString(msg)
}
