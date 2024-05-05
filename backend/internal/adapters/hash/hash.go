package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/hanoys/sigma-music/internal/ports"
	"time"
)

type HashPasswordProvider struct {
}

func NewHashPasswordProvider() *HashPasswordProvider {
	return &HashPasswordProvider{}
}

func (h *HashPasswordProvider) genSalt() string {
	salt := sha256.Sum256([]byte(time.Now().String()))
	return hex.EncodeToString(salt[:])
}

func (h *HashPasswordProvider) EncodePassword(password string) ports.SaltedPassword {
	salt := h.genSalt()
	hash := sha256.Sum256([]byte(password + salt))
	return ports.SaltedPassword{
		HashPassword: hex.EncodeToString(hash[:]),
		Salt:         salt,
	}
}
func (h *HashPasswordProvider) ComparePasswordWithHash(password string, saltedPassword ports.SaltedPassword) bool {
	passwordHash := sha256.Sum256([]byte(password + saltedPassword.Salt))
	return hex.EncodeToString(passwordHash[:]) == password
}
