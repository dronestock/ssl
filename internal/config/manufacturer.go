package config

import (
	"context"

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

	titles map[string]bool
}

func (m *Manufacturer) Refresh(ctx context.Context, base drone.Base, certificate *Certificate) (err error) {
	if nil == m.titles {
		m.titles = make(map[string]bool)
	}

	refreshers := make([]feature.Refresher, 0, 1)
	if nil != m.Chuangcache {
		refreshers = append(refreshers, manufacturer.NewChuangcache())
	}
	if nil != m.Tencent {
		refreshers = append(refreshers, manufacturer.NewTencent(m.Tencent))
	}
	for _, _refresher := range refreshers {
		_ = m.refresh(ctx, base, _refresher, certificate)
	}

	return
}

func (m *Manufacturer) Clean(ctx context.Context) (err error) {
	cleaners := make([]feature.Cleaner, 0, 1)
	if nil != m.Chuangcache {
		cleaners = append(cleaners, manufacturer.NewChuangcache())
	}
	if nil != m.Tencent {
		cleaners = append(cleaners, manufacturer.NewTencent(m.Tencent))
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
	certificate *Certificate,
) (err error) {
	fields := gox.Fields[any]{
		field.New("certificate", certificate),
	}
	if id, ue := refresher.Upload(ctx, &certificate.Certificate); nil != ue {
		err = ue
		base.Warn("上传证书出错", fields.Add(field.Error(err))...)
	} else if domains, me := refresher.Domains(ctx); nil != me {
		err = me
		base.Warn("获取域名列表出错", fields.Add(field.Error(err))...)
	} else if be := m.binds(ctx, refresher, certificate, id, domains); nil != be {
		err = be
		base.Warn("绑定域名出错", fields.Add(field.Error(err))...)
	} else {
		m.titles[certificate.Title] = true
	}

	return
}

func (m *Manufacturer) cleanup(ctx context.Context, cleaner feature.Cleaner) (err error) {
	if certificates, ce := cleaner.Certificates(ctx); nil != ce {
		err = ce
	} else {
		err = m.deletes(ctx, cleaner, certificates)
	}

	return
}

func (m *Manufacturer) deletes(
	ctx context.Context, cleaner feature.Cleaner,
	certificates []*core.ServerCertificate,
) (err error) {
	for _, certificate := range certificates {
		if _, ok := m.titles[certificate.Title]; !ok && core.CertificateStatusInuse != certificate.Status {
			err = cleaner.Delete(ctx, certificate)
		}
	}

	return
}

func (m *Manufacturer) binds(
	ctx context.Context,
	refresher feature.Refresher,
	certificate *Certificate,
	id string, domains []*core.Domain,
) (err error) {
	for _, cd := range domains {
		_domain := new(core.Domain)
		_domain.Id = cd.Id
		_domain.Name = cd.Name
		if certificate.Match(_domain) {
			err = refresher.Bind(ctx, id, _domain)
		}
	}

	return
}
