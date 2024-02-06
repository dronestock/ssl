package manufacturer

import (
	"context"
	"reflect"

	"github.com/dronestock/ssl/internal/internel/internal/core"
	"github.com/dronestock/ssl/internal/internel/internal/feature"
	"github.com/dronestock/ssl/internal/internel/internal/manufacturer/internal"
	"github.com/dronestock/ssl/internal/internel/internal/manufacturer/internal/apisix"
	"github.com/go-resty/resty/v2"
	"github.com/goexl/exception"
	"github.com/goexl/gox/field"
	"github.com/goexl/gox/http"
	"github.com/goexl/log"
)

var (
	_ feature.Refresher = (*Apisix)(nil)
	_ feature.Cleaner   = (*Apisix)(nil)
)

type Apisix struct {
	http   *resty.Client
	config *core.Apisix
	logger log.Logger
}

func NewApisix(http *resty.Client, config *core.Apisix, logger log.Logger) *Apisix {
	return &Apisix{
		http:   http,
		config: config,
		logger: logger,
	}
}

func (a *Apisix) Upload(
	ctx *context.Context, local *core.Certificate,
) (certificate *core.ServerCertificate, err error) {
	req := new(apisix.UploadReq)
	req.Protocols = []string{"TLSv1.1", "TLSv1.2", "TLSv1.3"}

	*ctx = context.WithValue(*ctx, local.SniKey(), []string{"*.itcoursee.com"})
	domains := (*ctx).Value(local.SniKey())
	if nil != domains {
		req.SNIs = domains.([]string)
	}

	rsp := new(apisix.Response[apisix.UploadRsp])
	url := a.config.Url()
	if le := local.Load(req); nil != le {
		err = le
	} else if ce := a.call(ctx, url, req, rsp, http.MethodPost); nil != ce {
		err = ce
	} else if apisix.StatusOk != rsp.Value.Code() {
		status := field.New("status", rsp.Value.Code())
		info := field.New("info", rsp.Value.Message())
		err = exception.New().Code(rsp.Value.Code()).Message("网关操作失败").Field(status, info).Build()
	} else {
		certificate = new(core.ServerCertificate)
		certificate.Id = rsp.Value.Id
	}

	return
}

func (a *Apisix) Bind(_ *context.Context, _ *core.ServerCertificate, _ *core.Domain) (record *core.Record, err error) {
	return
}

func (a *Apisix) Check(_ *context.Context, _ *core.Record) (checked bool, err error) {
	checked = true

	return
}

func (a *Apisix) Domains(_ *context.Context) (domains []*core.Domain, err error) {
	return
}

func (a *Apisix) Invalidates(ctx *context.Context, certificate *core.Certificate) (certificates []*core.ServerCertificate, err error) {
	req := new(internal.Empty)
	page := new(apisix.Page[*apisix.Certificate])
	if ce := a.call(ctx, a.config.Url(), req, page, http.MethodGet); nil != ce {
		err = ce
	} else {
		certificates = make([]*core.ServerCertificate, 0)
		a.invalidates(ctx, page.List, certificate, &certificates)
	}

	return
}

func (a *Apisix) Delete(ctx *context.Context, certificate *core.ServerCertificate) (err error) {
	req := new(internal.Empty)
	rsp := new(apisix.UploadRsp)
	url := a.config.Id(certificate.Id)
	err = a.call(ctx, url, req, rsp, http.MethodDelete)

	return
}

func (a *Apisix) call(ctx *context.Context, url string, req any, rsp any, method http.Method) (err error) {
	request := a.http.R().SetContext(*ctx).SetBody(req).SetResult(rsp)
	request.SetHeader(apisix.HeaderKey, a.config.Key)
	if hr, pe := request.Execute(method.Uppercase(), url); nil != pe {
		err = pe
	} else if hr.IsError() {
		err = exception.New().Code(hr.StatusCode()).Message(hr.String()).Build()
		a.logger.Warn("网关返回错误", field.New("status.code", hr.StatusCode()))
	}

	return
}

func (a *Apisix) invalidates(
	ctx *context.Context,
	certificates []*apisix.Certificate, certificate *core.Certificate,
	invalidates *[]*core.ServerCertificate,
) {
	for _, cert := range certificates {
		a.invalidate(ctx, cert, certificate, invalidates)
	}

	return
}

func (a *Apisix) invalidate(
	ctx *context.Context,
	check *apisix.Certificate,
	certificate *core.Certificate,
	invalidates *[]*core.ServerCertificate,
) {
	value := (*ctx).Value(certificate.SniKey())
	domains := make([]string, 0)
	if nil != value {
		domains = value.([]string)
	}
	if !a.check(domains, check.Value.SNIs) {
		return
	}

	if invalidate, err := certificate.Invalidate(check.Value.Cert); nil == err && invalidate {
		*invalidates = append(*invalidates, &core.ServerCertificate{
			Id: check.Value.Id,
		})
	}
}

func (a *Apisix) check(original []string, check []string) bool {
	return reflect.DeepEqual(original, check)
}
