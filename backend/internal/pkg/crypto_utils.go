package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

type CryptoUtils struct{}

var cryptoUtilsOnce sync.Once
var cryptoUtils *CryptoUtils

func NewCryptoUtils() *CryptoUtils {
	cryptoUtilsOnce.Do(func() {
		cryptoUtils = &CryptoUtils{}
	})
	return cryptoUtils
}

type CryptoUtilsPasswordHashInput struct {
	Email    string
	Username string
	Password string
}

func (c *CryptoUtils) GeneratePasswordHash(input *CryptoUtilsPasswordHashInput) string {
	combined := fmt.Sprintf("%s|%s|%s", input.Email, input.Username, input.Password)
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func (c *CryptoUtils) VerifyPasswordHash(hashedPassword string, input *CryptoUtilsPasswordHashInput) bool {
	expectedHash := c.GeneratePasswordHash(input)
	return hashedPassword == expectedHash
}
