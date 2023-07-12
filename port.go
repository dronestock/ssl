package main

type port struct {
	Http int `default:"${PORT_HTTP=8080}" json:"http,omitempty" Validate:"max=65535"`
}
