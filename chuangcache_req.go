package main

var _ tokenSetter = (*chuangcacheReq)(nil)

type chuangcacheReq struct {
	AccessToken string `json:"access_token"`
}

func (cr *chuangcacheReq) token(token string) (req tokenSetter) {
	cr.AccessToken = token
	req = cr

	return
}
