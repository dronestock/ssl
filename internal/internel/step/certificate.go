package step

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dronestock/drone"
	"github.com/dronestock/ssl/internal/internel/config"
	"github.com/dronestock/ssl/internal/internel/internal/constant"
	"github.com/goexl/gox"
	"github.com/goexl/gox/args"
	"github.com/goexl/gox/field"
	"github.com/goexl/gox/rand"
)

type Certificate struct {
	*drone.Base

	binary       string
	dir          string
	email        string
	environments map[string]string
	server       string
	certificates []*config.Certificate
	aliases      map[string]string
	prefix       map[string]string
}

func NewCertificate(base *drone.Base,
	binary string, dir string, email string, environments map[string]string, server string,
	certificates []*config.Certificate,
) *Certificate {
	return &Certificate{
		Base: base,

		binary:       binary,
		dir:          dir,
		email:        email,
		environments: environments,
		server:       server,
		certificates: certificates,
		aliases: map[string]string{
			"aliyun": "ali",
			"dnspod": "dp",
		},
		prefix: map[string]string{
			"ali": "Ali",
		},
	}
}

func (c *Certificate) Runnable() bool {
	return true
}

func (c *Certificate) Run(ctx *context.Context) (err error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(c.certificates))
	for _, certificate := range c.certificates {
		go c.run(ctx, certificate, wg, &err)
	}
	wg.Wait()

	return
}

func (c *Certificate) run(ctx *context.Context, certificate *config.Certificate, wg *sync.WaitGroup, err *error) {
	defer wg.Done()

	// 统一域名配置
	if "" != certificate.Domain {
		certificate.Domains = append(certificate.Domains, certificate.Domain)
	}
	// 加入顶级域名
	self := ""
	for _, domain := range certificate.Domains {
		domains := strings.Split(domain, constant.Dot)
		length := len(domains)
		self = strings.Join(domains[length-2:length], constant.Dot)
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
	} else {
		// 注入域名
		*ctx = context.WithValue(*ctx, certificate.SniKey(), certificate.Domains)
	}
}

func (c *Certificate) make(ctx *context.Context, certificate *config.Certificate) (err error) {
	if "" == os.Getenv(constant.Simulate) {
		err = c.makeAcme(ctx, certificate)
	} else {
		err = c.makeSigned(ctx, certificate)
	}

	return
}

func (c *Certificate) makeAcme(_ *context.Context, certificate *config.Certificate) (err error) {
	ma := args.New().Build()
	// 强制生成证书
	ma.Flag("force")
	ma.Flag("issue")
	// 生成日志
	ma.Flag("log")
	// 使用DNS验证验证所有者
	ma.Option("dns", c.dns(certificate))
	ma.Option("email", c.email)
	ma.Option("server", c.server)
	if abs, ae := filepath.Abs(c.dir); nil == ae {
		ma.Option("home", abs)
	}
	// 组装所有域名
	for _, domain := range certificate.Domains {
		ma.Option("domain", domain)
	}

	command := c.Command(c.binary)
	command.Args(ma.Build())

	env := command.Environment()
	for key, value := range c.environments {
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

func (c *Certificate) makeSigned(_ *context.Context, certificate *config.Certificate) (err error) {
	key := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIDj4TQ0rDAKGA/sRbus/lXRS01LtKM1BPbB+LZlMW4zZoAoGCCqGSM49
AwEHoUQDQgAE6irAQyPM78k5C17qM+KYF4XauZqV+IuBzQV2iBvIfzRMpTtFmqxN
jPUbWKuAM8mZgQ1sBSHf1Noeaqbbvx3NjQ==
-----END EC PRIVATE KEY-----
`
	chain := `-----BEGIN CERTIFICATE-----
MIIEMDCCAxigAwIBAgISBLDiU7O/tJ+LqmGUo+mcY+HzMA0GCSqGSIb3DQEBCwUA
MDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD
EwJSMzAeFw0yNDAxMzExNTA5MTdaFw0yNDA0MzAxNTA5MTZaMBoxGDAWBgNVBAMM
DyouaXRjb3Vyc2VlLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABOoqwEMj
zO/JOQte6jPimBeF2rmalfiLgc0FdogbyH80TKU7RZqsTYz1G1irgDPJmYENbAUh
39TaHmqm278dzY2jggIhMIICHTAOBgNVHQ8BAf8EBAMCB4AwHQYDVR0lBBYwFAYI
KwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFCDx0pag
vUU7YpK+zeT2bX/haK2YMB8GA1UdIwQYMBaAFBQusxe3WFbLrlAJQOYfr52LFMLG
MFUGCCsGAQUFBwEBBEkwRzAhBggrBgEFBQcwAYYVaHR0cDovL3IzLm8ubGVuY3Iu
b3JnMCIGCCsGAQUFBzAChhZodHRwOi8vcjMuaS5sZW5jci5vcmcvMCkGA1UdEQQi
MCCCDyouaXRjb3Vyc2VlLmNvbYINaXRjb3Vyc2VlLmNvbTATBgNVHSAEDDAKMAgG
BmeBDAECATCCAQUGCisGAQQB1nkCBAIEgfYEgfMA8QB3AEiw42vapkc0D+VqAvqd
MOscUgHLVt0sgdm7v6s52IRzAAABjWBIJusAAAQDAEgwRgIhAKr0cSSdUQMSCM+M
7KAc2d3lk7cVUPFoXorSQUsbbLwbAiEAty+bxGtxv5h7bPNIYUrAkcS+B1ea1aAh
eNrpUc+QVwgAdgA7U3d1Pi25gE6LMFsG/kA7Z9hPw/THvQANLXJv4frUFwAAAY1g
SCboAAAEAwBHMEUCIG1GvTVDAnV3qZSRFysSPFx+NZzhEDMJjFOnRJc/gB3pAiEA
1v2Tz6lX3ub1mXU7NYUY5S1xMuGmdD2RtRFIeDFNmlowDQYJKoZIhvcNAQELBQAD
ggEBAHeuARNdbwj8e4nIEwWrIM4XrGUDKxPe4kYMqYE019BF20sn8gNfpdus/ShG
IbeqHOAXYbqTdsWrRMyEpod1vfx5h//gOemOySH8v2kaiohWqYpKaktSqhGltEKM
UQ8oO9U31drMVRYLD7kMmsM/nkbY2VSwiWUryIe+F5rQrj7vrg09Obztm7fDER9h
uclLX126YDnRmEhDa8rxjyIaMPf+FoCxIU5MkUeE5yogNYhdLZxgLbqb4DErWqPS
kWDX80Y0TUcbELIDEup9+BX1NiJ9Pna/ZzYB2j/KU7Dq7B2mDaJbYKtR5TuOhIul
uPFfzxrJJdGnTeRJ72Rt3apUFPY=
-----END CERTIFICATE-----`
	if wke := os.WriteFile(certificate.Key(), []byte(key), os.ModePerm); nil != wke {
		err = wke
	} else if wce := os.WriteFile(certificate.Chain(), []byte(chain), os.ModePerm); nil != wce {
		err = wce
	}

	return
}

func (c *Certificate) install(ctx *context.Context, certificate *config.Certificate) (err error) {
	if "" == os.Getenv(constant.Simulate) {
		err = c.installAcme(ctx, certificate)
	} else {
		err = c.installSigned(ctx, certificate)
	}

	return
}

func (c *Certificate) installAcme(_ *context.Context, certificate *config.Certificate) (err error) {
	ia := args.New().Build()
	ia.Flag("installcert")
	if abs, ae := filepath.Abs(c.dir); nil == ae {
		ia.Option("home", abs)
	}

	for _, domain := range certificate.Domains {
		ia.Option("domain", domain)
	}
	ia.Option("certpath", certificate.Cert())
	ia.Option("key-file", certificate.Key())
	ia.Option("fullchain-file", certificate.Chain())
	if _, err = c.Command(c.binary).Args(ia.Build()).Build().Exec(); nil != err {
		c.Error("安装证书出错", field.New("certificate", certificate), field.Error(err))
	}

	return
}

func (c *Certificate) installSigned(_ *context.Context, _ *config.Certificate) (err error) {
	return
}

func (c *Certificate) mkdir(certificate *config.Certificate) (err error) {
	if _, se := os.Stat(certificate.Id); nil != se && os.IsNotExist(se) {
		err = os.MkdirAll(certificate.Id, os.ModePerm)
	}
	if nil == err {
		c.Cleanup().File(certificate.Id).Build()
	}

	return
}

func (c *Certificate) provider(certificate *config.Certificate) (provider string) {
	if dp, ok := c.aliases[certificate.Type]; ok {
		provider = dp
	} else {
		provider = certificate.Type
	}

	return
}

func (c *Certificate) dns(certificate *config.Certificate) string {
	return gox.StringBuilder(constant.Dns, constant.Underline, c.provider(certificate)).String()
}

func (c *Certificate) key(certificate *config.Certificate, key string) string {
	provider := c.provider(certificate)
	prefix := strings.ToUpper(provider)
	if cached, ok := c.prefix[provider]; ok {
		prefix = cached
	}

	return gox.StringBuilder(prefix, constant.Underline, key).String()
}
