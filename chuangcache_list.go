package main

var _ tokenSetter = (*chuangcacheListReq)(nil)

type (
	chuangcacheListReq struct {
		chuangcacheReqV2

		PageSize int
		PageNo   int
	}

	chuangcacheListRsp struct {
		Certificates []*chuangcacheCertificate `json:"DataSet,omitempty"`
		Total        int
		Page         int
		Size         int
	}
)

func (clr *chuangcacheListReq) token(token string) (req tokenSetter) {
	clr.AccessToken = token
	req = clr

	return
}
