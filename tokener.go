package main

type tokener interface {
	token(token string) tokener
}
