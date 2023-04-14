package main

import (
	"context"
	"sync"

	"github.com/goexl/gox/args"
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

func (sc *stepCertificate) Run(_ context.Context) (err error) {
	return
}

func (sc *stepCertificate) make(ctx context.Context, certificate *certificate, wg *sync.WaitGroup) (err error) {
	wg.Add(1)
	defer wg.Done()

	ma := args.New().Build()
	ma.Flag("force") // 强制生成证书
	ma.Flag("issue")
	ma.Flag("log") // 生成日志

	return
}
