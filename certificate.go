package main

import (
	"path/filepath"
)

type certificate struct {
	// 域名
	Domain string `json:"domain,omitempty" validate:"required_without=Domains"`
	// 域名列表
	Domains []string `json:"domains,omitempty" validate:"required_without=Domain"`
	// 类型
	Type string `default:"dp" json:"type,omitempty"`
	// 环境变量
	Environments map[string]string `json:"environments,omitempty"`

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
