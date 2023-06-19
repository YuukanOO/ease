package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// Generates a unique prefix for a given string using the sha256 algorithm.
func Prefix(s string, prefixSize int) string {
	hash := sha256.Sum256([]byte(s))
	hashString := hex.EncodeToString(hash[:])
	return hashString[:prefixSize]
}
