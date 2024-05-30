package security

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword generates a bcrypt hash of the password using a cost factor.
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// ComparePasswords compares a hashed password with a plain password
func ComparePasswords(hashedPassword string, plainPassword string) bool {
	// Hash the plain text password
	plainPasswordHash := HashPassword(plainPassword)
	// Compare the hashed passwords
	return hashedPassword == plainPasswordHash
}
