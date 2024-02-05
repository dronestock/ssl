package chuangcache

import (
	"github.com/dronestock/ssl/internal/internel/internal/core"
)

var _ core.TokenSetter = (*ListReq)(nil)

type (
	ListReq struct {
		RequestV2

		PageSize int
		PageNo   int
	}

	ListRsp struct {
		Certificates []*Certificate `json:"DataSet,omitempty"`
		Total        int
		Page         int
		Size         int
	}
)

func (lr *ListReq) Token(token string) (req core.TokenSetter) {
	lr.AccessToken = token
	req = lr

	return
}
