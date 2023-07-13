package config

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/dronestock/drone"
	"github.com/dronestock/ssl/internal/core"
	"github.com/dronestock/ssl/internal/feature"
	"github.com/dronestock/ssl/internal/manufacturer"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type Manufacturer struct {
	Chuangcache *core.Chuangcache `default:"${CHUANGCACHE}" json:"chuangcache,omitempty"`
	Tencent     *core.Tencent     `default:"${TENCENT}" json:"tencent,omitempty"`
}

func (m *Manufacturer) Refresh(ctx context.Context, base drone.Base, certificate *Certificate) (err error) {
	refreshers := make([]feature.Refresher, 0, 1)
	if nil != m.Chuangcache {
		refreshers = append(refreshers, manufacturer.NewChuangcache())
	}
	if nil != m.Tencent {
		if tencent, te := manufacturer.NewTencent(m.Tencent, base.Logger); nil != te {
			err = te
		} else {
			refreshers = append(refreshers, tencent)
		}
	}
	if nil != err {
		return
	}

	total := len(refreshers)
	for _, _refresher := range refreshers {
		if re := m.refresh(ctx, base, _refresher, certificate); nil != re {
			err = re
			total++
		}
	}
	if 0 == total {
		err = nil
	}

	return
}

func (m *Manufacturer) Clean(ctx context.Context, base drone.Base) (err error) {
	cleaners := make([]feature.Cleaner, 0, 1)
	if nil != m.Chuangcache {
		cleaners = append(cleaners, manufacturer.NewChuangcache())
	}
	if nil != m.Tencent {
		if tencent, te := manufacturer.NewTencent(m.Tencent, base.Logger); nil != te {
			err = te
		} else {
			cleaners = append(cleaners, tencent)
		}
	}
	if nil != err {
		return
	}

	for _, cleaner := range cleaners {
		_ = m.cleanup(ctx, cleaner)
	}

	return
}

func (m *Manufacturer) refresh(
	ctx context.Context,
	base drone.Base,
	refresher feature.Refresher,
	local *Certificate,
) (err error) {
	if certificate, ue := m.upload(ctx, base, refresher, &local.Certificate); nil != ue {
		err = ue
	} else if domains, me := refresher.Domains(ctx); nil != me {
		err = me
	} else if be := m.binds(ctx, base, refresher, local, certificate, domains); nil != be {
		err = be
	}

	return
}

func (m *Manufacturer) cleanup(ctx context.Context, cleaner feature.Cleaner) (err error) {
	if certificates, ce := cleaner.Invalidates(ctx); nil != ce {
		err = ce
	} else {
		err = m.deletes(ctx, cleaner, certificates)
	}

	return
}

func (m *Manufacturer) upload(
	ctx context.Context,
	base drone.Base,
	refresher feature.Refresher,
	local *core.Certificate,
) (cert *core.ServerCertificate, err error) {
	fields := gox.Fields[any]{
		field.New("title", local.Title),
		field.New("domains", local.Domains),
	}
	base.Debug("证书上传开始", fields...)
	if certificate, ue := refresher.Upload(ctx, local); nil != ue {
		err = ue
		base.Warn("证书上传出错", fields.Add(field.Error(ue))...)
	} else {
		cert = certificate
		base.Info("证书上传成功", fields...)
	}

	return
}

func (m *Manufacturer) deletes(
	ctx context.Context, cleaner feature.Cleaner,
	certificates []*core.ServerCertificate,
) (err error) {
	for _, certificate := range certificates {
		err = m.delete(ctx, cleaner, certificate)
	}

	return
}

func (m *Manufacturer) binds(
	ctx context.Context,
	base drone.Base,
	refresher feature.Refresher,
	local *Certificate,
	certificate *core.ServerCertificate, domains []*core.Domain,
) (err error) {
	wg := new(sync.WaitGroup)
	for _, domain := range domains {
		if local.Match(domain) {
			wg.Add(1)
			go m.bind(ctx, base, wg, refresher, certificate, domain)
		}
	}
	wg.Wait()

	return
}

func (m *Manufacturer) bind(
	ctx context.Context,
	base drone.Base,
	wg *sync.WaitGroup,
	refresher feature.Refresher,
	cert *core.ServerCertificate, domain *core.Domain,
) {
	defer wg.Done()

	fields := gox.Fields[any]{
		field.New("domain.id", domain.Id),
		field.New("domain.name", domain.Name),

		field.New("certificate.id", cert.Id),
		field.New("certificate.title", cert.Title),
	}
	base.Info("绑定证书开始", fields...)
	if record, be := refresher.Bind(ctx, cert, domain); nil != be {
		base.Warn("绑定证书失败", fields...)
	} else {
		// 检查部署是否完成
		m.wait(ctx, base, refresher, cert, record)
		base.Info("绑定证书成功", fields...)
	}

	return
}

func (m *Manufacturer) wait(
	ctx context.Context, base drone.Base,
	refresher feature.Refresher,
	cert *core.ServerCertificate,
	record *core.Record,
) {
	checked := false

	for times := 0; times < math.MaxInt; times++ {
		fields := gox.Fields[any]{
			field.New("certificate.id", cert.Id),
			field.New("certificate.title", cert.Title),
			field.New("times", times+1),
		}
		base.Info("检查证书部署开始", fields...)
		if check, ce := refresher.Check(ctx, record); nil != ce || !check {
			time.Sleep(10 * time.Second)
		} else if check {
			checked = true
			base.Info("检查证书部署成功", fields...)
		}

		if checked {
			break
		}
	}
}
func (m *Manufacturer) delete(ctx context.Context, cleaner feature.Cleaner, cert *core.ServerCertificate) (err error) {
	total := 15
	for times := 0; times < total; times++ {
		if deleted, de := cleaner.Delete(ctx, cert); nil != de && times < total-1 {
			time.Sleep(3 * time.Second)
		} else if nil != de {
			err = de
		} else if deleted {
			break
		}
	}

	return
}
