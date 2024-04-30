package ports

type IHashPasswordProvider interface {
	EncodePassword(string) string
	ComparePasswordWithHash(password string, hash string) bool
}
