package main

import (
	"context"
	"os"
	"path/filepath"
	"sync"

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

func (sc *stepCertificate) Runnable() bool {
	return true
}

func (sc *stepCertificate) Run(ctx context.Context) (err error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(sc.Certificates))
	for _, certificate := range sc.Certificates {
		go sc.run(ctx, certificate, wg, &err)
	}
	wg.Wait()

	return
}

func (sc *stepCertificate) run(ctx context.Context, certificate *certificate, wg *sync.WaitGroup, err *error) {
	defer wg.Done()

	// 统一域名配置
	if "" != certificate.Domain {
		certificate.Domains = append(certificate.Domains, certificate.Domain)
	}
	// 清理证书生成中间的过程文件
	certificate.id = rand.New().String().Build().Generate()

	if me := sc.mkdir(certificate); nil != me {
		*err = me
	} else if me := sc.make(ctx, certificate); nil != me {
		*err = me
	} else if ie := sc.install(ctx, certificate); nil != ie {
		*err = ie
	}

	return
}

func (sc *stepCertificate) make(_ context.Context, certificate *certificate) (err error) {
	ma := args.New().Build()
	ma.Flag("force") // 强制生成证书
	ma.Flag("issue")
	ma.Flag("log")                        // 生成日志
	ma.Option("dns", sc.dns(certificate)) // 使用DNS验证验证所有者
	ma.Option("email", sc.Email)
	ma.Option("server", "letsencrypt")
	ma.Flag("standalone").Option("httpport", sc.Port.Http)
	if abs, ae := filepath.Abs(sc.Dir); nil == ae {
		ma.Option("home", abs)
	}
	// 组装所有域名
	for _, domain := range certificate.Domains {
		ma.Option("domain", domain)
	}

	command := sc.Command(sc.Binary)
	command.Args(ma.Build())

	env := command.Environment()
	for key, value := range sc.Environments {
		env.Kv(sc.key(certificate, key), value)
	}
	for key, value := range certificate.Environments {
		env.Kv(sc.key(certificate, key), value)
	}
	env.Build()

	if _, err = command.Build().Exec(); nil != err {
		sc.Error("生成证书出错", field.New("certificate", certificate), field.Error(err))
	}

	return
}

func (sc *stepCertificate) install(_ context.Context, certificate *certificate) (err error) {
	ia := args.New().Build()
	ia.Flag("installcert")
	if abs, ae := filepath.Abs(sc.Dir); nil == ae {
		ia.Option("home", abs)
	}

	for _, domain := range certificate.Domains {
		ia.Option("domain", domain)
	}
	ia.Option("certpath", certificate.cert())
	ia.Option("key-file", certificate.key())
	ia.Option("fullchain-file", certificate.chain())
	if _, err = sc.Command(sc.Binary).Args(ia.Build()).Build().Exec(); nil != err {
		sc.Error("安装证书出错", field.New("certificate", certificate), field.Error(err))
	}

	return
}

func (sc *stepCertificate) mkdir(certificate *certificate) (err error) {
	if _, se := os.Stat(certificate.id); nil != se && os.IsNotExist(se) {
		err = os.MkdirAll(certificate.id, os.ModePerm)
	}
	if nil == err {
		sc.Cleanup().File(certificate.id).Build()
	}

	return
}