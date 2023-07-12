package refresher

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/dronestock/ssl/internal"
	"github.com/dronestock/ssl/internal/chuangcache"
	"github.com/dronestock/ssl/internal/config"
	"github.com/go-resty/resty/v2"
	"github.com/goexl/exc"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/gox/rand"
	"github.com/goexl/simaqian"
)

var _ internal.Refresher = (*Chuangcache)(nil)

type Chuangcache struct {
	simaqian.Logger

	http   *resty.Client
	config *config.Chuangcache
	token  *internal.Token
}

func NewChuangcache() *Chuangcache {
	return &Chuangcache{}
}

func (c *Chuangcache) Init(_ context.Context) (err error) {
	return
}

func (c *Chuangcache) Upload(ctx context.Context, certificate *internal.Certificate) (id string, err error) {
	req := new(chuangcache.UploadReq)
	// ! 为避免证书名字重复，在证书名字上加上随机字符串
	req.Title = gox.StringBuilder(certificate.Title, rand.New().String().Build().Generate()).String()
	rsp := new(chuangcache.Response[*chuangcache.UploadRsp])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "config/addCertificate")
	if le := certificate.Load(req); nil != le {
		err = le
	} else if ce := c.call(ctx, url, req, rsp); nil != ce {
		err = ce
	} else {
		id = rsp.Data.Id
	}

	return
}

func (c *Chuangcache) Bind(ctx context.Context, id string, domain *internal.Domain) (err error) {
	req := new(chuangcache.BindReq)
	req.Id = id
	req.Domain = domain.Id
	rsp := new(chuangcache.Response[bool])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "config/bindDomainCertificate")
	if ce := c.call(ctx, url, req, rsp); nil != ce {
		err = ce
	} else if !rsp.Data {
		err = exc.NewFields("绑定证书失败", field.New("req", req), field.New("rsp", rsp))
	}

	return
}

func (c *Chuangcache) Domains(ctx context.Context) (domains []*internal.Domain, err error) {
	req := new(chuangcache.Request)
	rsp := new(chuangcache.Response[[]*chuangcache.Domain])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "Domain/domainList")
	if ce := c.call(ctx, url, req, rsp); nil != ce {
		err = ce
	} else if 0 != len(rsp.Data) {
		for _, cd := range rsp.Data {
			domain := new(internal.Domain)
			domain.Id = cd.Id
			domain.Name = cd.Name
			domain.Type = internal.DomainTypeCdn

			domains = append(domains, domain)
		}
	}

	return
}

func (c *Chuangcache) Certificates(ctx context.Context) (certificates []*internal.ServerCertificate, err error) {
	certificates = make([]*internal.ServerCertificate, 0, 1)
	err = c.certificates(ctx, &certificates)

	return
}

func (c *Chuangcache) Delete(ctx context.Context, certificate *internal.ServerCertificate) (err error) {
	req := new(chuangcache.DeleteReq)
	req.Key = certificate.Id
	rsp := new(chuangcache.Response[bool])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "config/deleteCertificate")
	if ce := c.call(ctx, url, req, rsp); nil != ce {
		err = ce
	} else if !rsp.Data {
		err = exc.NewFields("删除证书失败", field.New("req", req), field.New("rsp", rsp))
	}

	return
}

func (c *Chuangcache) certificates(ctx context.Context, certificates *[]*internal.ServerCertificate) (err error) {
	req := new(chuangcache.ListReq)
	req.PageSize = math.MaxInt
	req.PageNo = 1
	rsp := new(chuangcache.Response[*chuangcache.ListRsp])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "v2/certificate/list")
	if ce := c.call(ctx, url, req, rsp); nil != ce {
		err = ce
	} else if 0 != len(rsp.Data.Certificates) {
		for _, certificate := range rsp.Data.Certificates {
			sc := new(internal.ServerCertificate)
			sc.Id = certificate.Key
			sc.Title = certificate.Title
			sc.Status = certificate.InternalStatus()

			*certificates = append(*certificates, sc)
		}
	}

	return
}

func (c *Chuangcache) call(ctx context.Context, url string, req internal.TokenSetter, rsp internal.StatusCoder) (err error) {
	if _token, te := c.getToken(ctx); nil != te {
		err = te
	} else if ce := c.send(ctx, url, req.Token(_token), rsp); nil != ce {
		err = ce
	}

	return
}

func (c *Chuangcache) getToken(ctx context.Context) (_token string, err error) {
	if nil != c.token && c.token.Validate() {
		_token = c.token.Token
	}
	if "" != _token {
		return
	}

	req := chuangcache.TokenReq{
		Ak: c.config.Ak,
		Sk: c.config.Sk,
	}
	rsp := new(chuangcache.Response[*chuangcache.TokenRsp])
	url := fmt.Sprintf("%s/%s", chuangcache.ApiEndpoint, "OAuth/authorize")
	if err = c.send(ctx, url, req, rsp); nil == err {
		c.token = new(internal.Token)
		c.token.Token = rsp.Data.AccessToken
		c.token.Expired = time.Now().Add(time.Duration(1000 * rsp.Data.ExpiresIn))
		_token = rsp.Data.AccessToken
	}

	return
}

func (c *Chuangcache) send(ctx context.Context, url string, req any, rsp internal.StatusCoder) (err error) {
	if hr, pe := c.http.R().SetContext(ctx).SetBody(req).SetResult(rsp).Post(url); nil != pe {
		err = pe
	} else if hr.IsError() {
		c.Warn("创世云返回错误", field.New("status.code", hr.StatusCode()))
	} else if chuangcache.StatusOk != rsp.Code() {
		status := field.New("status", rsp.Code())
		info := field.New("info", rsp.Message())
		err = exc.NewException(rsp.Code(), "创世云操作失败", status, info)
	}

	return
}
