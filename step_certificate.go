package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dronestock/ssl/internal/config"
	"github.com/goexl/gox/args"
	"github.com/goexl/gox/field"
	"github.com/goexl/gox/rand"
)

type stepCertificate struct {
	*plugin
}

func newStepCertificate(plugin *plugin) *stepCertificate {
	return &stepCertificate{
		plugin: plugin,
	}
}

func (c *stepCertificate) Runnable() bool {
	return true
}

func (c *stepCertificate) Run(ctx context.Context) (err error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(c.Certificates))
	for _, _certificate := range c.Certificates {
		go c.run(ctx, _certificate, wg, &err)
	}
	wg.Wait()

	return
}

func (c *stepCertificate) run(ctx context.Context, certificate *config.Certificate, wg *sync.WaitGroup, err *error) {
	defer wg.Done()

	// 统一域名配置
	if "" != certificate.Domain {
		certificate.Domains = append(certificate.Domains, certificate.Domain)
	}
	// 加入顶级域名
	self := ""
	for _, domain := range certificate.Domains {
		domains := strings.Split(domain, dot)
		length := len(domains)
		self = strings.Join(domains[length-2:length], dot)
		if "" != self && self != domain {
			break
		}
	}
	certificate.Domains = append(certificate.Domains, self)

	// 清理证书生成中间的过程文件
	certificate.Id = rand.New().String().Build().Generate()

	if mke := c.mkdir(certificate); nil != mke {
		*err = mke
	} else if mae := c.make(ctx, certificate); nil != mae {
		*err = mae
	} else if ie := c.install(ctx, certificate); nil != ie {
		*err = ie
	}
}

func (c *stepCertificate) make(_ context.Context, certificate *config.Certificate) (err error) {
	ma := args.New().Build()
	// 强制生成证书
	ma.Flag("force")
	ma.Flag("issue")
	// 生成日志
	ma.Flag("log")
	// 使用DNS验证验证所有者
	ma.Option("dns", c.dns(certificate))
	ma.Option("email", c.Email)
	ma.Option("server", c.Server)
	if abs, ae := filepath.Abs(c.Dir); nil == ae {
		ma.Option("home", abs)
	}
	// 组装所有域名
	for _, domain := range certificate.Domains {
		ma.Option("domain", domain)
	}

	command := c.Command(c.Binary)
	command.Args(ma.Build())

	env := command.Environment()
	for key, value := range c.Environments {
		env.Kv(c.key(certificate, key), value)
	}
	for key, value := range certificate.Environments {
		env.Kv(c.key(certificate, key), value)
	}
	env.Build()

	if _, err = command.Build().Exec(); nil != err {
		c.Error("生成证书出错", field.New("certificate", certificate), field.Error(err))
	}

	return
}

func (c *stepCertificate) install(_ context.Context, certificate *config.Certificate) (err error) {
	ia := args.New().Build()
	ia.Flag("installcert")
	if abs, ae := filepath.Abs(c.Dir); nil == ae {
		ia.Option("home", abs)
	}

	for _, domain := range certificate.Domains {
		ia.Option("domain", domain)
	}
	ia.Option("certpath", certificate.Cert())
	ia.Option("key-file", certificate.Key())
	ia.Option("fullchain-file", certificate.Chain())
	if _, err = c.Command(c.Binary).Args(ia.Build()).Build().Exec(); nil != err {
		c.Error("安装证书出错", field.New("certificate", certificate), field.Error(err))
	}

	return
}

func (c *stepCertificate) mkdir(certificate *config.Certificate) (err error) {
	if _, se := os.Stat(certificate.Id); nil != se && os.IsNotExist(se) {
		err = os.MkdirAll(certificate.Id, os.ModePerm)
	}
	if nil == err {
		c.Cleanup().File(certificate.Id).Build()
	}

	return
}
