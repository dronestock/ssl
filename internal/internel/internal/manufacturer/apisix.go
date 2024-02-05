package manufacturer

import (
	"context"

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
	ctx context.Context, local *core.Certificate,
) (certificate *core.ServerCertificate, err error) {
	req := new(apisix.UploadReq)
	req.SNIs = []string{"*.itcoursee.com"}
	req.Protocols = []string{"TLSv1.1", "TLSv1.2", "TLSv1.3"}

	rsp := new(apisix.UploadRsp)
	url := a.config.Url()
	if le := local.Load(req); nil != le {
		err = le
	} else if ce := a.call(ctx, url, req, rsp, http.MethodPost); nil != ce {
		err = ce
	} else {
		certificate = new(core.ServerCertificate)
		certificate.Id = rsp.Id
	}

	return
}

func (a *Apisix) Bind(_ context.Context, _ *core.ServerCertificate, _ *core.Domain) (record *core.Record, err error) {
	return
}

func (a *Apisix) Check(_ context.Context, _ *core.Record) (checked bool, err error) {
	checked = true

	return
}

func (a *Apisix) Domains(_ context.Context) (domains []*core.Domain, err error) {
	return
}

func (a *Apisix) Invalidates(_ context.Context) (certificates []*core.ServerCertificate, err error) {

	return
}

func (a *Apisix) Delete(ctx context.Context, certificate *core.ServerCertificate) (deleted bool, err error) {
	req := new(internal.Empty)
	rsp := new(apisix.UploadRsp)
	url := a.config.Id(certificate.Id)
	err = a.call(ctx, url, req, rsp, http.MethodDelete)

	return
}

func (a *Apisix) call(ctx context.Context, url string, req any, rsp core.StatusCoder, method http.Method) (err error) {
	response := new(apisix.Response)
	response.Value = rsp

	request := a.http.R().SetContext(ctx).SetBody(req).SetResult(response)
	request.SetHeader(apisix.HeaderKey, a.config.Key)
	if hr, pe := request.Execute(method.Uppercase(), url); nil != pe {
		err = pe
	} else if hr.IsError() {
		a.logger.Warn("网关返回错误", field.New("status.code", hr.StatusCode()))
	} else if apisix.StatusOk != rsp.Code() {
		status := field.New("status", rsp.Code())
		info := field.New("info", rsp.Message())
		err = exception.New().Code(rsp.Code()).Message("创世云操作失败").Field(status, info).Build()
	}

	return
}
