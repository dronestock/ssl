package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/dronestock/drone"
)

func TestChuangcache(t *testing.T) {
	_chuangcache := new(chuangcache)
	_chuangcache.Ak = "NdK8SrkrwOf8RuuN"
	_chuangcache.Sk = "1hqMT3rgmvAudrlwbrIFe9Z8UAdjLsIY"
	_certificate := new(certificate)
	_certificate.id = "8Zx6PWE9"
	_certificate.Domain = "test.dronestock.tech"

	ctx := context.Background()
	base := new(drone.Base)
	fmt.Println(_chuangcache.refresh(ctx, base, _certificate))
}
