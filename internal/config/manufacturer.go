package config

import (
	"context"

	"github.com/dronestock/drone"
	"github.com/dronestock/ssl/internal"
	"github.com/dronestock/ssl/internal/refresher"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type Manufacturer struct {
	Chuangcache *Chuangcache `default:"${CHUANGCACHE}" json:"chuangcache,omitempty"`
	Tencent     *Tencent     `default:"${TENCENT}" json:"tencent,omitempty"`

	titles map[string]bool
}

func (m *Manufacturer) Refresh(ctx context.Context, base drone.Base, certificates []*internal.Certificate) (err error) {
	if nil == m.titles {
		m.titles = make(map[string]bool)
	}

	refreshers := make([]internal.Refresher, 0, 1)
	if nil != m.Chuangcache {
		refreshers = append(refreshers, refresher.NewChuangcache())
	}
	if nil != m.Tencent {
		refreshers = append(refreshers, refresher.NewTencent(m.Tencent))
	}
	for _, certificate := range certificates {
		for _, _refresher := range refreshers {
			_ = m.refresh(ctx, base, _refresher, certificate)
		}
	}

	return
}

func (m *Manufacturer) refresh(
	ctx context.Context,
	base drone.Base,
	refresher internal.Refresher,
	certificate *internal.Certificate,
) (err error) {
	fields := gox.Fields[any]{
		field.New("certificate", certificate),
	}
	if id, ue := refresher.Upload(ctx, certificate); nil != ue {
		err = ue
		base.Warn("上传证书出错", fields.Add(field.Error(err))...)
	} else if domains, me := refresher.Domains(ctx); nil != me {
		err = me
		base.Warn("获取域名列表出错", fields.Add(field.Error(err))...)
	} else if be := m.binds(ctx, refresher, certificate, id, domains); nil != be {
		err = be
		base.Warn("绑定域名出错", fields.Add(field.Error(err))...)
	} else if ce := m.cleanup(ctx, refresher); nil != ce {
		err = ce
		base.Warn("清理证书出错", fields.Add(field.Error(err))...)
	}

	return
}

func (m *Manufacturer) cleanup(ctx context.Context, refresher internal.Refresher) (err error) {
	if certificates, ce := refresher.Certificates(ctx); nil != ce {
		err = ce
	} else {
		err = m.deletes(ctx, refresher, certificates)
	}

	return
}

func (m *Manufacturer) deletes(
	ctx context.Context, refresher internal.Refresher,
	certificates []*internal.ServerCertificate,
) (err error) {
	for _, certificate := range certificates {
		if _, ok := m.titles[certificate.Title]; !ok && internal.CertificateStatusInuse != certificate.Status {
			err = refresher.Delete(ctx, certificate)
		}
	}

	return
}

func (m *Manufacturer) binds(
	ctx context.Context,
	refresher internal.Refresher,
	certificate *internal.Certificate,
	id string, domains []*internal.Domain,
) (err error) {
	for _, cd := range domains {
		_domain := new(internal.Domain)
		_domain.Id = cd.Id
		_domain.Name = cd.Name
		if certificate.Match(_domain) {
			err = refresher.Bind(ctx, id, _domain)
		}
	}

	return
}
