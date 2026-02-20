package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// HashPasswordMD5 hashes a password using MD5 to match the legacy PHP application.
// In a future update, this should be migrated to bcrypt.
func HashPasswordMD5(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
