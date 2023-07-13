package core

type TokenSetter interface {
	Token(token string) TokenSetter
}
