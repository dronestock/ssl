package chuangcache

import (
	"github.com/dronestock/ssl/internal"
)

var (
	_ internal.Loader      = (*UploadReq)(nil)
	_ internal.TokenSetter = (*UploadReq)(nil)
)

type (
	UploadReq struct {
		Request

		Title       string `json:"ssl_title,omitempty"`
		Private     string `json:"private_key,omitempty"`
		Certificate string `json:"certificate,omitempty"`
	}

	UploadRsp struct {
		Id string `json:"ssl_key,omitempty"`
	}
)

func (ur *UploadReq) Token(token string) (req internal.TokenSetter) {
	ur.AccessToken = token
	req = ur

	return
}

func (ur *UploadReq) Cert(cert string) {
	ur.Certificate = cert
}

func (ur *UploadReq) Key(key string) {
	ur.Private = key
}

func (ur *UploadReq) Chain(_ string) {}
