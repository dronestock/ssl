package core

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Certificate struct {
	// 标题
	Title string `json:"title,omitempty" validate:"required"`
	// 域名
	Domain string `json:"domain,omitempty" validate:"required_without=Domains"`
	// 域名列表
	Domains []string `json:"domains,omitempty" validate:"required_without=Domain"`
	// 类型
	Type string `default:"dp" json:"type,omitempty"`

	// 用于内部使用，确定一个证书的后续操作标识
	Id string
}

func (c *Certificate) Match(check *Domain) (matched bool) {
	if "" != c.Domain {
		c.Domains = append(c.Domains, c.Domain)
	}
	for _, domain := range c.Domains {
		if check.Name == domain {
			matched = true
		} else if check.Name == strings.ReplaceAll(domain, "*.", "") {
			matched = true
		} else if match, me := path.Match(domain, check.Name); nil == me {
			matched = match
		}

		if matched {
			break
		}
	}

	return
}

func (c *Certificate) Load(loader Loader) (err error) {
	/*if ce := c.set(c.Cert(), loader.Cert); nil != ce {
		err = ce
	} else if ke := c.set(c.Key(), loader.Key); nil != ke {
		err = ke
	} else if fe := c.set(c.Chain(), loader.Chain); nil != fe {
		err = fe
	}*/
	loader.Cert(`-----BEGIN CERTIFICATE-----
MIIEMDCCAxigAwIBAgISBLDiU7O/tJ+LqmGUo+mcY+HzMA0GCSqGSIb3DQEBCwUA
MDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD
EwJSMzAeFw0yNDAxMzExNTA5MTdaFw0yNDA0MzAxNTA5MTZaMBoxGDAWBgNVBAMM
DyouaXRjb3Vyc2VlLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABOoqwEMj
zO/JOQte6jPimBeF2rmalfiLgc0FdogbyH80TKU7RZqsTYz1G1irgDPJmYENbAUh
39TaHmqm278dzY2jggIhMIICHTAOBgNVHQ8BAf8EBAMCB4AwHQYDVR0lBBYwFAYI
KwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFCDx0pag
vUU7YpK+zeT2bX/haK2YMB8GA1UdIwQYMBaAFBQusxe3WFbLrlAJQOYfr52LFMLG
MFUGCCsGAQUFBwEBBEkwRzAhBggrBgEFBQcwAYYVaHR0cDovL3IzLm8ubGVuY3Iu
b3JnMCIGCCsGAQUFBzAChhZodHRwOi8vcjMuaS5sZW5jci5vcmcvMCkGA1UdEQQi
MCCCDyouaXRjb3Vyc2VlLmNvbYINaXRjb3Vyc2VlLmNvbTATBgNVHSAEDDAKMAgG
BmeBDAECATCCAQUGCisGAQQB1nkCBAIEgfYEgfMA8QB3AEiw42vapkc0D+VqAvqd
MOscUgHLVt0sgdm7v6s52IRzAAABjWBIJusAAAQDAEgwRgIhAKr0cSSdUQMSCM+M
7KAc2d3lk7cVUPFoXorSQUsbbLwbAiEAty+bxGtxv5h7bPNIYUrAkcS+B1ea1aAh
eNrpUc+QVwgAdgA7U3d1Pi25gE6LMFsG/kA7Z9hPw/THvQANLXJv4frUFwAAAY1g
SCboAAAEAwBHMEUCIG1GvTVDAnV3qZSRFysSPFx+NZzhEDMJjFOnRJc/gB3pAiEA
1v2Tz6lX3ub1mXU7NYUY5S1xMuGmdD2RtRFIeDFNmlowDQYJKoZIhvcNAQELBQAD
ggEBAHeuARNdbwj8e4nIEwWrIM4XrGUDKxPe4kYMqYE019BF20sn8gNfpdus/ShG
IbeqHOAXYbqTdsWrRMyEpod1vfx5h//gOemOySH8v2kaiohWqYpKaktSqhGltEKM
UQ8oO9U31drMVRYLD7kMmsM/nkbY2VSwiWUryIe+F5rQrj7vrg09Obztm7fDER9h
uclLX126YDnRmEhDa8rxjyIaMPf+FoCxIU5MkUeE5yogNYhdLZxgLbqb4DErWqPS
kWDX80Y0TUcbELIDEup9+BX1NiJ9Pna/ZzYB2j/KU7Dq7B2mDaJbYKtR5TuOhIul
uPFfzxrJJdGnTeRJ72Rt3apUFPY=
-----END CERTIFICATE-----`)
	loader.Key(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIDj4TQ0rDAKGA/sRbus/lXRS01LtKM1BPbB+LZlMW4zZoAoGCCqGSM49
AwEHoUQDQgAE6irAQyPM78k5C17qM+KYF4XauZqV+IuBzQV2iBvIfzRMpTtFmqxN
jPUbWKuAM8mZgQ1sBSHf1Noeaqbbvx3NjQ==
-----END EC PRIVATE KEY-----
`)

	return
}

func (c *Certificate) set(path string, setter Setter) (err error) {
	if bytes, re := os.ReadFile(path); nil != re {
		err = re
	} else {
		setter(string(bytes))
	}

	return
}

func (c *Certificate) Cert() string {
	return filepath.Join(c.Id, "cert.pem")
}

func (c *Certificate) Key() string {
	return filepath.Join(c.Id, "key.pem")
}

func (c *Certificate) Chain() string {
	return filepath.Join(c.Id, "chain.pem")
}
