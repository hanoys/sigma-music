package ports

import "github.com/hanoys/sigma-music/internal/domain"

type IHashPasswordProvider interface {
	EncodePassword(password string) domain.SaltedPassword
	ComparePasswordWithHash(password string, saltedPassword domain.SaltedPassword) bool
}
