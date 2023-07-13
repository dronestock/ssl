package main

import (
	"strings"

	"github.com/dronestock/drone"
	"github.com/dronestock/ssl/internal/config"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type plugin struct {
	drone.Base
	config.Manufacturer

	// 执行程序
	Binary string `default:"${BINARY=acme.sh}"`
	// 目录
	Dir string `default:"${DIR=.}"`
	// 证书
	Certificate *config.Certificate `default:"${CERTIFICATE}" Validate:"required_without=Certificates"`
	// 证书列表
	Certificates []*config.Certificate `default:"${CERTIFICATES}" Validate:"required_without=Certificate"`

	// 邮箱
	Email string `default:"${EMAIL}" Validate:"required,email"`
	// 环境变量
	Environments map[string]string `default:"${ENVIRONMENTS}" json:"environments,omitempty"`
	// 端口
	Port config.Port `default:"${PORT}"`
	// 证书服务器
	// nolint: lll
	Server string `default:"${SERVER=zerossl}" Validate:"oneof=letsencrypt letsencrypt_test buypass buypass_test zerossl sslcom google googletest"`

	// 别名
	aliases map[string]string
}

func newPlugin() drone.Plugin {
	return &plugin{
		aliases: map[string]string{
			"aliyun": "ali",
			"dnspod": "dp",
		},
	}
}

func (p *plugin) Config() drone.Config {
	return p
}

func (p *plugin) Setup() (err error) {
	if nil == p.Certificates {
		p.Certificates = make([]*config.Certificate, 0, 1)
	}
	if nil != p.Certificate {
		p.Certificates = append(p.Certificates, p.Certificate)
	}

	return
}

func (p *plugin) Steps() drone.Steps {
	return drone.Steps{
		drone.NewStep(newStepCertificate(p)).Name("证书").Build(),
		drone.NewStep(newStepRefresh(p)).Name("创世云").Build(),
	}
}

func (p *plugin) Fields() gox.Fields[any] {
	return gox.Fields[any]{
		field.New("certificates", p.Certificates),
		field.New("manufacturer", p.Manufacturer),
		field.New("email", p.Email),
		field.New("environments", p.Environments),
	}
}

func (p *plugin) provider(certificate *config.Certificate) (provider string) {
	if dp, ok := p.aliases[certificate.Type]; ok {
		provider = dp
	} else {
		provider = certificate.Type
	}

	return
}

func (p *plugin) dns(certificate *config.Certificate) string {
	return gox.StringBuilder(dns, underline, p.provider(certificate)).String()
}

func (p *plugin) key(certificate *config.Certificate, key string) string {
	return gox.StringBuilder(strings.ToUpper(p.provider(certificate)), underline, key).String()
}
