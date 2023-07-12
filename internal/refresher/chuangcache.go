package refresher

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/dronestock/drone"
	"github.com/dronestock/ssl/internal"
	"github.com/dronestock/ssl/internal/chuangcache"
	"github.com/goexl/exc"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/gox/rand"
)

var _ internal.Refresher = (*Chuangcache)(nil)

type Chuangcache struct {
	token  *internal.Token
	titles map[string]bool
}

func (c *Chuangcache) refresh(ctx context.Context, base drone.Base, certificate *internal.Certificate) (err error) {
	fields := gox.Fields[any]{
		field.New("certificate", certificate),
	}
	if id, ue := c.upload(ctx, base, certificate); nil != ue {
		err = ue
		base.Warn("上传证书出错", fields.Add(field.Error(err))...)
	} else if domains, me := c.domains(ctx, base); nil != me {
		err = me
		base.Warn("获取域名列表出错", fields.Add(field.Error(err))...)
	} else if be := c.binds(ctx, base, certificate, id, domains); nil != be {
		err = be
		base.Warn("绑定域名出错", fields.Add(field.Error(err))...)
	} else if ce := c.cleanup(ctx, base); nil != ce {
		err = ce
		base.Warn("清理证书出错", fields.Add(field.Error(err))...)
	}

	return
}

func (c *Chuangcache) upload(ctx context.Context, base drone.Base, certificate *internal.Certificate) (id string, err error) {
	if nil == c.titles {
		c.titles = make(map[string]bool)
	}

	req := new(chuangcache.UploadReq)
	// ! 为避免证书名字重复，在证书名字上加上随机字符串
	req.Title = gox.StringBuilder(certificate.Title, rand.New().String().Build().Generate()).String()
	rsp := new(chuangcache.Response[*chuangcache.UploadRsp])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "config/addCertificate")
	if le := certificate.Load(req); nil != le {
		err = le
	} else if ce := c.call(ctx, base, url, req, rsp); nil != ce {
		err = ce
	} else {
		id = rsp.Data.Id
		c.titles[req.Title] = true
	}

	return
}

func (c *Chuangcache) binds(
	ctx context.Context, base drone.Base,
	certificate *internal.Certificate,
	id string, domains []*chuangcache.Domain,
) (err error) {
	for _, cd := range domains {
		_domain := new(internal.Domain)
		_domain.Id = cd.Id
		_domain.Name = cd.Name
		if certificate.Match(_domain) {
			err = c.bind(ctx, base, id, _domain)
		}
	}

	return
}

func (c *Chuangcache) bind(
	ctx context.Context, base drone.Base,
	id string, domain *internal.Domain,
) (err error) {
	req := new(chuangcache.BindReq)
	req.Id = id
	req.Domain = domain.Id
	rsp := new(chuangcache.Response[bool])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "config/bindDomainCertificate")
	if ce := c.call(ctx, base, url, req, rsp); nil != ce {
		err = ce
	} else if !rsp.Data {
		err = exc.NewFields("绑定证书失败", field.New("req", req), field.New("rsp", rsp))
	}

	return
}

func (c *Chuangcache) domains(ctx context.Context, base drone.Base) (domains []*chuangcache.Domain, err error) {
	req := new(chuangcache.Request)
	rsp := new(chuangcache.Response[[]*chuangcache.Domain])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "Domain/domainList")
	if ce := c.call(ctx, base, url, req, rsp); nil != ce {
		err = ce
	} else {
		domains = rsp.Data
	}

	return
}

func (c *Chuangcache) cleanup(ctx context.Context, base drone.Base) (err error) {
	req := new(chuangcache.ListReq)
	req.PageSize = math.MaxInt
	req.PageNo = 1
	rsp := new(chuangcache.Response[*chuangcache.ListRsp])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "v2/certificate/list")
	if ce := c.call(ctx, base, url, req, rsp); nil != ce {
		err = ce
	} else {
		err = c.deletes(ctx, base, rsp.Data.Certificates)
	}

	return
}

func (c *Chuangcache) deletes(
	ctx context.Context, base drone.Base,
	certificates []*chuangcache.Certificate,
) (err error) {
	for _, _certificate := range certificates {
		if _, ok := c.titles[_certificate.Title]; !ok && chuangcache.CertificateStatusInuse != _certificate.Status {
			err = c.delete(ctx, base, _certificate.Key)
		}
	}

	return
}

func (c *Chuangcache) delete(ctx context.Context, base drone.Base, key string) (err error) {
	req := new(chuangcache.DeleteReq)
	req.Key = key
	rsp := new(chuangcache.Response[bool])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "config/deleteCertificate")
	if ce := c.call(ctx, base, url, req, rsp); nil != ce {
		err = ce
	} else if !rsp.Data {
		err = exc.NewFields("删除证书失败", field.New("req", req), field.New("rsp", rsp))
	}

	return
}

func (c *Chuangcache) call(
	ctx context.Context, base drone.Base,
	url string, req internal.TokenSetter, rsp internal.StatusCoder,
) (err error) {
	if _token, te := c.getToken(ctx, base); nil != te {
		err = te
	} else if ce := c.send(ctx, base, url, req.Token(_token), rsp); nil != ce {
		err = ce
	}

	return
}

func (c *Chuangcache) getToken(ctx context.Context, base drone.Base) (_token string, err error) {
	if nil != c.token && c.token.Validate() {
		_token = c.token.Token
	}
	if "" != _token {
		return
	}

	req := chuangcache.TokenReq{
		Ak: c.Ak,
		Sk: c.Sk,
	}
	rsp := new(chuangcache.Response[*chuangcache.TokenRsp])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "OAuth/authorize")
	if err = c.send(ctx, base, url, req, rsp); nil == err {
		c.token = new(internal.Token)
		c.token.Token = rsp.Data.AccessToken
		c.token.Expired = time.Now().Add(time.Duration(1000 * rsp.Data.ExpiresIn))
		_token = rsp.Data.AccessToken
	}

	return
}

func (c *Chuangcache) send(ctx context.Context, base drone.Base, url string, req any, rsp internal.StatusCoder) (err error) {
	if hr, pe := base.Http().SetContext(ctx).SetBody(req).SetResult(rsp).Post(url); nil != pe {
		err = pe
	} else if hr.IsError() {
		base.Warn("创世云返回错误", field.New("status.code", hr.StatusCode()))
	} else if chuangcache.StatusOk != rsp.Code() {
		status := field.New("status", rsp.Code())
		info := field.New("info", rsp.Message())
		err = exc.NewException(rsp.Code(), "创世云操作失败", status, info)
	}

	return
}
