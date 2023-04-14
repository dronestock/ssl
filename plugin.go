package main

import (
	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type plugin struct {
	drone.Base

	// 执行程序
	Binary string `default:"${BINARY=acme.sh}"`
	// 证书
	Certificate *certificate `default:"${CERTIFICATE}" validate:"required_without=Certificates"`
	// 证书列表
	Certificates []*certificate `default:"${CERTIFICATES}" validate:"required_without=Certificate"`

	// 别名
	aliases map[string]string
}

func newPlugin() drone.Plugin {
	return &plugin{
		aliases: map[string]string{
			"aliyun": "ali",
		},
	}
}

func (p *plugin) Config() drone.Config {
	return p
}

func (p *plugin) Setup() (err error) {
	if nil != p.Certificate {
		p.Certificates = append(p.Certificates, p.Certificate)
	}

	return
}

func (p *plugin) Steps() drone.Steps {
	return drone.Steps{
		drone.NewStep(newStepCertificate(p)).Name("证书").Build(),
		drone.NewStep(newStepChuangcache(p)).Name("创世云").Build(),
	}
}

func (p *plugin) Fields() gox.Fields[any] {
	return gox.Fields[any]{
		field.New("certificates", p.Certificates),
	}
}
