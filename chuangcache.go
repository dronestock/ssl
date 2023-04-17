package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type chuangcache struct {
	Ak string `json:"ak,omitempty"`
	Sk string `json:"sk,omitempty"`

	token *token
}

func (c *chuangcache) refresh(ctx context.Context, base *drone.Base, certificate *certificate) (err error) {
	return
}

func (c *chuangcache) upload(ctx context.Context, base *drone.Base, certificate *certificate) (id string, err error) {
	req := new(chuangcacheUploadReq)
	rsp := new(chuangcacheRsp[*chuangcacheUploadRsp])
	url := fmt.Sprintf("%s/%s", chuangcacheApiEndpoint, "domain/domainList")
	if _token, te := c.getToken(ctx, base); nil != te {
		err = te
	} else if ce := c.call(ctx, base, url, req.token(_token), rsp); nil != ce {
		err = ce
	} else {
		id = rsp.Data.Id
	}

	return
}

func (c *chuangcache) bind(ctx context.Context, base *drone.Base, certificate *certificate) (domains []*domain, err error) {
	return
}

func (c *chuangcache) unbind(ctx context.Context, base *drone.Base, certificate *certificate) (domains []*domain, err error) {
	return
}

func (c *chuangcache) match(ctx context.Context, base *drone.Base, certificate *certificate) (domains []*domain, err error) {
	req := new(chuangcacheReq)
	rsp := new(chuangcacheRsp[[]*chuangcacheDomain])
	url := fmt.Sprintf("%s/%s", chuangcacheApiEndpoint, "domain/domainList")
	if _token, te := c.getToken(ctx, base); nil != te {
		err = te
	} else if ce := c.call(ctx, base, url, req.token(_token), rsp); nil != ce {
		err = ce
	} else {
		domains = make([]*domain, 0, 1)
		for _, cd := range rsp.Data {
			_domain := new(domain)
			_domain.id = cd.Id
			_domain.name = cd.Name
			domains = append(domains, gox.If(certificate.match(_domain), _domain))
		}
	}

	return
}

func (c *chuangcache) getToken(ctx context.Context, base *drone.Base) (_token string, err error) {
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
	rsp := new(chuangcacheTokenRsp)
	url := fmt.Sprintf("%s/%s", chuangcacheApiEndpoint, "OAuth/authorize")
	if err = c.call(ctx, base, url, req, rsp); nil == err {
		c.token = new(token)
		c.token.token = rsp.Data.AccessToken
		c.token.expiresIn = time.Now().Add(time.Duration(1000 * rsp.Data.ExpiresIn))
		_token = rsp.Data.AccessToken
	}

	return
}

func (c *chuangcache) call(ctx context.Context, base *drone.Base, url string, req any, rsp any) (err error) {
	if hr, pe := base.Http().SetContext(ctx).SetBody(req).SetResult(rsp).Post(url); nil != pe {
		err = pe
	} else if hr.IsError() {
		base.Warn("创世云返回错误", field.New("status.code", hr.StatusCode()))
	} else {
		// TODO
	}

	return
}
