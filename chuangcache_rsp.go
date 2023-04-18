package main

var _ statusCoder = (*chuangcacheRsp[bool])(nil)

type chuangcacheRsp[T any] struct {
	// 接口返回码
	// 0：操作失败
	// 1：操作成功
	Status int `json:"status"`
	// 接口返回信息
	Info string `json:"info"`
	// 数据
	Data T `json:"data"`
}

func (cr *chuangcacheRsp[T]) code() int {
	return cr.Status
}

func (cr *chuangcacheRsp[T]) message() string {
	return cr.Info
}
