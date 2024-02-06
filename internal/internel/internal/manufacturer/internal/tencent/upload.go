package tencent

import (
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type UploadReq struct {
	request *ssl.UploadCertificateRequest
}

func NewUploadReq() *UploadReq {
	return &UploadReq{
		request: ssl.NewUploadCertificateRequest(),
	}
}

func (ur *UploadReq) Request() *ssl.UploadCertificateRequest {
	return ur.request
}

func (ur *UploadReq) Cert(cert string) {
	ur.request.CertificatePublicKey = &cert
}

func (ur *UploadReq) Key(key string) {
	ur.request.CertificatePrivateKey = &key
}

func (ur *UploadReq) Chain(_ string) {
	// 空
}

func (ur *UploadReq) Alias(alias string) {
	ur.request.Alias = &alias
}
