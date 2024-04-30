package ports

type IHashPasswordProvider interface {
	EncodePassword(string) string
	DecodePassword(string) string
}
