package main

type port struct {
	Http int `default:"${PORT_HTTP=8080}" json:"http,omitempty" validate:"max=65535"`
}
