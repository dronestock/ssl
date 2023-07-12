package chuangcache

import (
	"github.com/dronestock/ssl/internal"
)

var _ internal.TokenSetter = (*ListReq)(nil)

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

func (lr *ListReq) Token(token string) (req internal.TokenSetter) {
	lr.AccessToken = token
	req = lr

	return
}
