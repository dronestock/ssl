package core

type Loader interface {
	Key(key string)

	Chain(chain string)
}
