package main

type loader interface {
	cert(cert string)
	key(key string)
	chain(chain string)
}
