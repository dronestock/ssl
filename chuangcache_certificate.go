package main

var (
	_ loader      = (*chuangcacheUploadReq)(nil)
	_ tokenSetter = (*chuangcacheUploadReq)(nil)
)

type (
	chuangcacheUploadReq struct {
		chuangcacheReq

		Title       string `json:"ssl_title,omitempty"`
		Key         string `json:"private_key,omitempty"`
		Certificate string `json:"certificate,omitempty"`
	}

	chuangcacheUploadRsp struct {
		Id string `json:"ssl_key,omitempty"`
	}
)

func (cur *chuangcacheUploadReq) token(token string) (req tokenSetter) {
	cur.AccessToken = token
	req = cur

	return
}

func (cur *chuangcacheUploadReq) cert(cert string) {
	cur.Certificate = cert
}

func (cur *chuangcacheUploadReq) key(key string) {
	cur.Key = key
}

func (cur *chuangcacheUploadReq) chain(_ string) {}
