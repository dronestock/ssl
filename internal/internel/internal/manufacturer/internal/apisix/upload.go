package apisix

import (
	"github.com/goexl/gox"
)

type (
	UploadReq struct {
		Content   string   `json:"cert,omitempty"`
		Private   string   `json:"key,omitempty"`
		SNIs      []string `json:"snis,omitempty"`
		Protocols []string `json:"ssl_protocols,omitempty"`
	}

	UploadRsp struct {
		Id     string `json:"id,omitempty"`
		Status int    `json:"status,omitempty"`
	}
)

func (ur *UploadReq) Key(key string) {
	ur.Private = key
}

func (ur *UploadReq) Chain(chain string) {
	ur.Content = chain
}

func (ur *UploadRsp) Code() int {
	return ur.Status
}

func (ur *UploadRsp) Message() string {
	return gox.ToString(ur.Status)
}
