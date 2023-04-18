package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/dronestock/drone"
	"github.com/goexl/exc"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/gox/rand"
)

type chuangcache struct {
	Ak string `json:"ak,omitempty"`
	Sk string `json:"sk,omitempty"`

	token *token
}

func (c *chuangcache) refresh(ctx context.Context, base drone.Base, certificate *certificate) (err error) {
	if id, ue := c.upload(ctx, base, certificate); nil != ue {
		err = ue
	} else if domains, me := c.domains(ctx, base); nil != me {
		err = me
	} else if be := c.binds(ctx, base, certificate, id, domains); nil != be {
		err = be
	} else if ce := c.cleanup(ctx, base); nil != ce {
		err = ce
	}

	return
}

func (c *chuangcache) upload(ctx context.Context, base drone.Base, certificate *certificate) (id string, err error) {
	req := new(chuangcacheUploadReq)
	// ! 为避免证书名字重复，在证书名字上加上随机字符串
	req.Title = gox.StringBuilder(certificate.Title, rand.New().String().Build().Generate()).String()
	rsp := new(chuangcacheRsp[*chuangcacheUploadRsp])
	url := fmt.Sprintf("%s/%s", chuangcacheApiEndpoint, "config/addCertificate")
	if le := certificate.load(req); nil != le {
		err = le
	} else if ce := c.call(ctx, base, url, req, rsp); nil != ce {
		err = ce
	} else {
		id = rsp.Data.Id
	}

	return
}

func (c *chuangcache) binds(
	ctx context.Context, base drone.Base,
	certificate *certificate,
	id string, domains []*chuangcacheDomain,
) (err error) {
	for _, cd := range domains {
		_domain := new(domain)
		_domain.id = cd.Id
		_domain.name = cd.Name
		if certificate.match(_domain) {
			err = c.bind(ctx, base, id, _domain)
		}
	}

	return
}

func (c *chuangcache) bind(
	ctx context.Context, base drone.Base,
	id string, domain *domain,
) (err error) {
	req := new(chuangcacheBindReq)
	req.Id = id
	req.Domain = domain.id
	rsp := new(chuangcacheRsp[bool])
	url := fmt.Sprintf("%s/%s", chuangcacheApiEndpoint, "config/bindDomainCertificate")
	if ce := c.call(ctx, base, url, req, rsp); nil != ce {
		err = ce
	} else if !rsp.Data {
		err = exc.NewFields("绑定证书失败", field.New("req", req), field.New("rsp", rsp))
	}

	return
}

func (c *chuangcache) domains(ctx context.Context, base drone.Base) (domains []*chuangcacheDomain, err error) {
	req := new(chuangcacheReq)
	rsp := new(chuangcacheRsp[[]*chuangcacheDomain])
	url := fmt.Sprintf("%s/%s", chuangcacheApiEndpoint, "domain/domainList")
	if ce := c.call(ctx, base, url, req, rsp); nil != ce {
		err = ce
	} else {
		domains = rsp.Data
	}

	return
}

func (c *chuangcache) cleanup(ctx context.Context, base drone.Base) (err error) {
	req := new(chuangcacheListReq)
	req.PageSize = math.MaxInt
	req.PageNo = 1
	rsp := new(chuangcacheRsp[*chuangcacheListRsp])
	url := fmt.Sprintf("%s/%s", chuangcacheApiEndpoint, "v2/certificate/list")
	if ce := c.call(ctx, base, url, req, rsp); nil != ce {
		err = ce
	} else {
		err = c.deletes(ctx, base, rsp.Data.Certificates)
	}

	return
}

func (c *chuangcache) deletes(
	ctx context.Context, base drone.Base,
	certificates []*chuangcacheCertificate,
) (err error) {
	for _, _certificate := range certificates {
		if chuangcacheCertificateStatusInuse != _certificate.Status {
			err = c.delete(ctx, base, _certificate.Key)
		}
	}

	return
}

func (c *chuangcache) delete(ctx context.Context, base drone.Base, key string) (err error) {
	req := new(chuangcacheDeleteReq)
	req.Key = key
	rsp := new(chuangcacheRsp[bool])
	url := fmt.Sprintf("%s/%s", chuangcacheApiEndpoint, "config/deleteCertificate")
	if ce := c.call(ctx, base, url, req, rsp); nil != ce {
		err = ce
	} else if !rsp.Data {
		err = exc.NewFields("删除证书失败", field.New("req", req), field.New("rsp", rsp))
	}

	return
}

func (c *chuangcache) call(
	ctx context.Context, base drone.Base,
	url string, req tokenSetter, rsp statusCoder,
) (err error) {
	if _token, te := c.getToken(ctx, base); nil != te {
		err = te
	} else if ce := c.send(ctx, base, url, req.token(_token), rsp); nil != ce {
		err = ce
	}

	return
}

func (c *chuangcache) getToken(ctx context.Context, base drone.Base) (_token string, err error) {
	if nil != c.token && c.token.validate() {
		_token = c.token.token
	}
	if "" != _token {
		return
	}

	req := chuangcacheTokenReq{
		Ak: c.Ak,
		Sk: c.Sk,
	}
	rsp := new(chuangcacheRsp[*chuangcacheTokenRsp])
	url := fmt.Sprintf("%s/%s", chuangcacheApiEndpoint, "OAuth/authorize")
	if err = c.send(ctx, base, url, req, rsp); nil == err {
		c.token = new(token)
		c.token.token = rsp.Data.AccessToken
		c.token.expiresIn = time.Now().Add(time.Duration(1000 * rsp.Data.ExpiresIn))
		_token = rsp.Data.AccessToken
	}

	return
}

func (c *chuangcache) send(ctx context.Context, base drone.Base, url string, req any, rsp statusCoder) (err error) {
	if hr, pe := base.Http().SetContext(ctx).SetBody(req).SetResult(rsp).Post(url); nil != pe {
		err = pe
	} else if hr.IsError() {
		base.Warn("创世云返回错误", field.New("status.code", hr.StatusCode()))
	} else if chuangcacheStatusOk != rsp.code() {
		status := field.New("status", rsp.code())
		info := field.New("info", rsp.message())
		err = exc.NewException(rsp.code(), "创世云操作失败", status, info)
	}

	return
}
