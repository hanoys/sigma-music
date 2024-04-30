package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

type HashPasswordProvider struct {
}

func NewHashPasswordProvider() *HashPasswordProvider {
	return &HashPasswordProvider{}
}

func (h *HashPasswordProvider) EncodePassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
func (h *HashPasswordProvider) ComparePasswordWithHash(password string, hash string) bool {
	passwordHash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(passwordHash[:]) == hash
}
