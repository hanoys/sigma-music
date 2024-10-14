package domain

type SaltedPassword struct {
	HashPassword string
	Salt         string
}
