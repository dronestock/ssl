package chuangcache

import (
	"github.com/dronestock/ssl/internal"
)

var _ internal.TokenSetter = (*BindReq)(nil)

type BindReq struct {
	Request

	Id     string `json:"ssl_key,omitempty"`
	Domain string `json:"domain_id,omitempty"`
}

func (br *BindReq) Token(token string) (req internal.TokenSetter) {
	br.AccessToken = token
	req = br

	return
}
