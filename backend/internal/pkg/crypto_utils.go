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

func (c *CryptoUtils) Hash256Password(email, username, password string) string {
	combined := fmt.Sprintf("%s|%s|%s", email, username, password)
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func (c *CryptoUtils) VerifyPassword(hashedPassword, email, username, password string) bool {
	expectedHash := c.Hash256Password(email, username, password)
	return hashedPassword == expectedHash
}
