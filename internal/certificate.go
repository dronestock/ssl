package internal

import (
	"os"
	"path"
	"path/filepath"

	"github.com/dronestock/ssl"
)

type Certificate struct {
	main.Manufacturer

	// 标题
	Title string `json:"title,omitempty" validate:"required"`
	// 域名
	Domain string `json:"domain,omitempty" validate:"required_without=Domains"`
	// 域名列表
	Domains []string `json:"domains,omitempty" validate:"required_without=Domain"`
	// 类型
	Type string `default:"dp" json:"type,omitempty"`
	// 环境变量
	Environments map[string]string `json:"environments,omitempty"`

	// 用于内部使用，确定一个证书的后续操作标识
	Id string
}

func (c *Certificate) Match(domain *Domain) (matched bool) {
	if "" != c.Domain {
		c.Domains = append(c.Domains, c.Domain)
	}
	for _, _domain := range c.Domains {
		if domain.Name == _domain {
			matched = true
		} else if match, me := path.Match(_domain, domain.Name); nil == me {
			matched = match
		}

		if matched {
			break
		}
	}

	return
}

func (c *Certificate) Load(loader Loader) (err error) {
	if ce := c.set(c.Cert(), loader.Cert); nil != ce {
		err = ce
	} else if ke := c.set(c.Key(), loader.Key); nil != ke {
		err = ke
	} else if fe := c.set(c.Chain(), loader.Chain); nil != fe {
		err = fe
	}

	return
}

func (c *Certificate) set(path string, setter Setter) (err error) {
	if bytes, re := os.ReadFile(path); nil != re {
		err = re
	} else {
		setter(string(bytes))
	}

	return
}

func (c *Certificate) Cert() string {
	return filepath.Join(c.Id, "cert.pem")
}

func (c *Certificate) Key() string {
	return filepath.Join(c.Id, "key.pem")
}

func (c *Certificate) Chain() string {
	return filepath.Join(c.Id, "chain.pem")
}
