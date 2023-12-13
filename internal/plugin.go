package internal

import (
	"github.com/dronestock/drone"
	"github.com/dronestock/ssl/internal/internel/config"
	"github.com/dronestock/ssl/internal/internel/step"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type Plugin struct {
	drone.Base
	config.Manufacturer

	// 执行程序
	Binary string `default:"${BINARY=acme.sh}"`
	// 目录
	Dir string `default:"${DIR=.}"`
	// 证书
	Certificate *config.Certificate `default:"${CERTIFICATE}" Validate:"required_without=Invalidates"`
	// 证书列表
	Certificates []*config.Certificate `default:"${CERTIFICATES}" Validate:"required_without=Certificate"`

	// 邮箱
	Email string `default:"${EMAIL}" Validate:"required,email"`
	// 环境变量
	Environments map[string]string `default:"${ENVIRONMENTS}" json:"environments,omitempty"`
	// 证书服务器
	// nolint: lll
	Server string `default:"${SERVER=letsencrypt}" Validate:"oneof=letsencrypt letsencrypt_test buypass buypass_test zerossl sslcom google googletest"`
}

func NewPlugin() drone.Plugin {
	return new(Plugin)
}

func (p *Plugin) Config() drone.Config {
	return p
}

func (p *Plugin) Setup() (err error) {
	if nil == p.Certificates {
		p.Certificates = make([]*config.Certificate, 0, 1)
	}
	if nil != p.Certificate {
		p.Certificates = append(p.Certificates, p.Certificate)
	}

	return
}

func (p *Plugin) Steps() drone.Steps {
	return drone.Steps{
		drone.NewStep(step.NewCertificate(
			&p.Base, p.Binary, p.Dir, p.Email, p.Environments, p.Server,
			p.Certificates,
		)).Name("证书").Build(),
		drone.NewStep(step.NewRefresh(&p.Base, &p.Manufacturer, p.Certificates)).Name("刷新").Build(),
		drone.NewStep(step.NewCleanup(&p.Base, &p.Manufacturer, p.Certificates)).Name("清理").Build(),
	}
}

func (p *Plugin) Fields() gox.Fields[any] {
	return gox.Fields[any]{
		field.New("certificates", p.Certificates),
		field.New("manufacturer", p.Manufacturer),
		field.New("email", p.Email),
		field.New("environments", p.Environments),
	}
}
