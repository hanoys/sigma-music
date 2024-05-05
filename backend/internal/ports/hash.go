package ports

type SaltedPassword struct {
	HashPassword string
	Salt         string
}

type IHashPasswordProvider interface {
	EncodePassword(password string) SaltedPassword
	ComparePasswordWithHash(password string, saltedPassword SaltedPassword) bool
}
