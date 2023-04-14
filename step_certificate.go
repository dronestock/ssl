package main

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/goexl/gox/args"
	"github.com/goexl/gox/field"
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
	for _, certificate := range sc.Certificates {
		go sc.run(ctx, certificate, wg, &err)
	}
	wg.Wait()

	return
}

func (sc *stepCertificate) run(ctx context.Context, certificate *certificate, wg *sync.WaitGroup, err *error) {
	wg.Add(1)
	defer wg.Done()

	if "" != certificate.Domain {
		certificate.Domains = append(certificate.Domains, certificate.Domain)
	}
	if me := sc.make(ctx, certificate); nil != me {
		*err = me
	} else if ie := sc.install(ctx, certificate); nil != ie {
		*err = ie
	}

	return
}

func (sc *stepCertificate) make(ctx context.Context, certificate *certificate) (err error) {
	ma := args.New().Build()
	ma.Flag("force") // 强制生成证书
	ma.Flag("issue")
	ma.Flag("log")                             // 生成日志
	ma.Option("dns", "")                       // 使用DNS验证验证所有者
	ma.Option("dnssleep", certificate.Timeout) // 超时时间，在给定的时间后，验证DNS的设置是否正确

	for _, domain := range certificate.Domains {
		ma.Option("domain", domain)
	}
	if _, err = sc.Command(sc.Binary).Args(ma.Build()).Build().Exec(); nil != err {
		sc.Error("生成证书出错", field.New("certificate", certificate), field.Error(err))
	}

	return
}

func (sc *stepCertificate) install(ctx context.Context, certificate *certificate) (err error) {
	ia := args.New().Build()
	ia.Flag("installcert")

	for _, domain := range certificate.Domains {
		ia.Option("domain", domain)
	}
	ia.Option("certpath", filepath.Join(certificate.id, "cert.pem"))
	ia.Option("key-file", filepath.Join(certificate.id, "privkey.pem"))
	ia.Option("fullchain-file", filepath.Join(certificate.id, "fullchain.pem"))
	if _, err = sc.Command(sc.Binary).Args(ia.Build()).Build().Exec(); nil != err {
		sc.Error("安装证书出错", field.New("certificate", certificate), field.Error(err))
	}

	return
}
