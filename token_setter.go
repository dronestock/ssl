package main

type tokenSetter interface {
	token(token string) tokenSetter
}
