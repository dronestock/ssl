package core

type Loader interface {
	Cert(cert string)

	Key(key string)

	Chain(chain string)
}
