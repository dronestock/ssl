package config

type Executable struct {
	// 执行程序
	Binary string `default:"${BINARY=acme.sh}"`
}
