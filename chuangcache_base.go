package main

var _ tokener = (*chuangcacheReq)(nil)

type (
	chuangcacheReq struct {
		// 密钥
		AccessToken string `json:"access_token"`
	}

	// 响应基类
	chuangcacheRsp[T any] struct {
		// 接口返回码
		// 0：操作失败
		// 1：操作成功
		Status int `json:"status"`
		// 接口返回信息
		Info string `json:"info"`
		// 数据
		Data T `json:"data"`
	}
)

func (cr *chuangcacheReq) token(token string) (req tokener) {
	cr.AccessToken = token
	req = cr

	return
}
