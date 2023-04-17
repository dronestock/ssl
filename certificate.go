package main

import (
	"os"
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

func (c *certificate) match(domain *domain) (matched bool) {
	return
}

func (c *certificate) load(loader loader) (err error) {
	if ce := c.set(c.cert(), loader.cert); nil != ce {
		err = ce
	} else if ke := c.set(c.key(), loader.key); nil != ke {
		err = ke
	} else if fe := c.set(c.chain(), loader.chain); nil != fe {
		err = fe
	}

	return
}

func (c *certificate) set(path string, setter setter) (err error) {
	if bytes, re := os.ReadFile(path); nil != re {
		err = re
	} else {
		setter(string(bytes))
	}

	return
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
