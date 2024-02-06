package apisix

type Response[T any] struct {
	Key   string `json:"key,omitempty"`
	Value T      `json:"value,omitempty"`
}
