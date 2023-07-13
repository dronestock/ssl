package chuangcache

import (
	"github.com/dronestock/ssl/internal/core"
)

var _ core.StatusCoder = (*Response[bool])(nil)

type Response[T any] struct {
	// 接口返回码
	// 0：操作失败
	// 1：操作成功
	Status int `json:"status"`
	// 接口返回信息
	Info string `json:"info"`
	// 数据
	Data T `json:"data"`
}

func (cr *Response[T]) Code() int {
	return cr.Status
}

func (cr *Response[T]) Message() string {
	return cr.Info
}
