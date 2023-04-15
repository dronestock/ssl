package main

import (
	"path/filepath"
	"time"
)

type certificate struct {
	// 域名
	Domain string `json:"domain,omitempty" validate:"required_without=Domains"`
	// 域名列表
	Domains []string `json:"domains,omitempty" validate:"required_without=Domain,dive,domain"`
	// 类型
	Type string `json:"type,omitempty" validate:"required,oneof=ali"`
	// 超时时间
	Timeout time.Duration `default:"15s" json:"timeout,omitempty"`

	// 用于内部使用，确定一个证书的后续操作标识
	id string
}

func (c *certificate) cert() string {
	return filepath.Join(c.id, "cert.pem")
}

func (c *certificate) key() string {
	return filepath.Join(c.id, "key.pem")
}

func (c *certificate) chain() string {
	return filepath.Join(c.id, "chain.pem")
}
