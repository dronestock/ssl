package core

import (
	"fmt"
)

type Apisix struct {
	// 端点
	Endpoint string `json:"endpoint,omitempty"`
	// 前缀
	Prefix string `default:"apisix/admin" json:"prefix,omitempty"`
	// 授权
	Key string `json:"key,omitempty" validate:"required"`
	// 接口
	Api string `default:"ssls" json:"api,omitempty"`
}

func (a *Apisix) Url() string {
	return fmt.Sprintf("%s/%s/%s", a.Endpoint, a.Prefix, a.Api)
}

func (a *Apisix) Id(id string) string {
	return fmt.Sprintf("%s/%s", a.Url(), id)
}
