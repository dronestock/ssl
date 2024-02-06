package apisix

type Page[T any] struct {
	Total int `json:"total,omitempty"`
	List  []T `json:"list,omitempty"`
}
