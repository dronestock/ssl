package tencent

import (
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type UploadReq struct {
	ssl.UploadCertificateRequest
}

func (ur *UploadReq) Request() *ssl.UploadCertificateRequest {
	return &ur.UploadCertificateRequest
}

func (ur *UploadReq) Cert(cert string) {
	ur.CertificatePublicKey = &cert
}

func (ur *UploadReq) Key(key string) {
	ur.CertificatePrivateKey = &key
}

func (ur *UploadReq) Chain(_ string) {}
