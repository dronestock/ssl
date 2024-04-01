package core

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dronestock/ssl/internal/internel/internal/key"
)

type Certificate struct {
	// 标题
	Title string `json:"title,omitempty" validate:"required"`
	// 域名
	Domain string `json:"domain,omitempty" validate:"required_without=Domains"`
	// 域名列表
	Domains []string `json:"domains,omitempty" validate:"required_without=Domain"`
	// 类型
	Type string `default:"dp" json:"type,omitempty"`

	// 用于内部使用，确定一个证书的后续操作标识
	Id string
}

func (c *Certificate) SniKey() key.Sni {
	return key.Sni(fmt.Sprintf("%s-sni", c.Id))
}

func (c *Certificate) Match(check *Domain) (matched bool) {
	if "" != c.Domain {
		c.Domains = append(c.Domains, c.Domain)
	}

	for _, domain := range c.Domains {
		if check.Name == domain {
			matched = true
		} else if check.Name == strings.ReplaceAll(domain, "*.", "") {
			matched = true
		} else if match, me := path.Match(domain, check.Name); nil == me {
			matched = match
		}

		if matched {
			break
		}
	}

	return
}

func (c *Certificate) Invalidate(chain string) (invalidate bool, err error) {
	if original, oe := c.load(); nil != oe {
		err = oe
	} else if check, ce := c.parse(chain); nil != ce {
		err = ce
	} else {
		invalidate = check.NotAfter.Before(original.NotAfter)
	}

	return
}

func (c *Certificate) Load(loader Loader) (err error) {
	if ke := c.set(c.Key(), loader.Key); nil != ke {
		err = ke
	} else if fe := c.set(c.Chain(), loader.Chain); nil != fe {
		err = fe
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

func (c *Certificate) set(path string, setter Setter) (err error) {
	if bytes, re := os.ReadFile(path); nil != re {
		err = re
	} else {
		setter(string(bytes))
	}

	return
}

func (c *Certificate) load() (certificate *x509.Certificate, err error) {
	if content, rfe := os.ReadFile(c.Chain()); nil != rfe {
		err = rfe
	} else {
		block, _ := pem.Decode(content)
		certificate, err = x509.ParseCertificate(block.Bytes)
	}

	return
}

func (c *Certificate) parse(chain string) (certificate *x509.Certificate, err error) {
	block, _ := pem.Decode([]byte(chain))
	certificate, err = x509.ParseCertificate(block.Bytes)

	return
}
