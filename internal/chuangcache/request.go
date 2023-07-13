package chuangcache

import (
	"github.com/dronestock/ssl/internal/core"
)

var _ core.TokenSetter = (*Request)(nil)

type (
	Request struct {
		AccessToken string `json:"access_token"`
	}

	RequestV2 struct {
		AccessToken string
	}
)

func (r *Request) Token(token string) (req core.TokenSetter) {
	r.AccessToken = token
	req = r

	return
}

func (r *RequestV2) Token(token string) (req core.TokenSetter) {
	r.AccessToken = token
	req = r

	return
}
