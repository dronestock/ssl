package apisix

type Response struct {
	Key   string `json:"key,omitempty"`
	Value any    `json:"value,omitempty"`
}
