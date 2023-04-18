package main

var _ tokenSetter = (*chuangcacheReq)(nil)

type (
	chuangcacheReq struct {
		AccessToken string `json:"access_token"`
	}

	chuangcacheReqV2 struct {
		AccessToken string
	}
)

func (cr *chuangcacheReq) token(token string) (req tokenSetter) {
	cr.AccessToken = token
	req = cr

	return
}

func (cr *chuangcacheReqV2) token(token string) (req tokenSetter) {
	cr.AccessToken = token
	req = cr

	return
}
