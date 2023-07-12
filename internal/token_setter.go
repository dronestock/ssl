package internal

type TokenSetter interface {
	Token(token string) TokenSetter
}
