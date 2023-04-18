package main

var _ tokenSetter = (*chuangcacheBindReq)(nil)

type chuangcacheBindReq struct {
	chuangcacheReq

	Id     string `json:"ssl_key,omitempty"`
	Domain string `json:"domain_id,omitempty"`
}

func (cbr *chuangcacheBindReq) token(token string) (req tokenSetter) {
	cbr.AccessToken = token
	req = cbr

	return
}
