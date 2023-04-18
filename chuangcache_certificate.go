package main

var _ loader = (*chuangcacheUploadReq)(nil)

type (
	chuangcacheUploadReq struct {
		chuangcacheReq

		Title       string `json:"ssl_title,omitempty"`
		Key         string `json:"private_key,omitempty"`
		Certificate string `json:"certificate,omitempty"`
	}

	chuangcacheUploadRsp struct {
		Id string `json:"ssl_key,omitempty"`
	}
)

func (cur *chuangcacheUploadReq) cert(cert string) {
	cur.Certificate = cert
	cur.Certificate = `-----BEGIN CERTIFICATE-----
MIIEYjCCA0qgAwIBAgISA65jw+FeYO+j+N5nSROsZQS3MA0GCSqGSIb3DQEBCwUA
MDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD
EwJSMzAeFw0yMzA0MTcwNjUzMDFaFw0yMzA3MTYwNjUzMDBaMB8xHTAbBgNVBAMT
FHRlc3QuZHJvbmVzdG9jay50ZWNoMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE
daMkPTiA3GteY1+Jt9GMFQNMHTA2o7x0mugEVOOJyEPMqOoV8QGSzSloRQ0pwHn6
Mll9fMCIjv1HDVvrskFpbqOCAk4wggJKMA4GA1UdDwEB/wQEAwIHgDAdBgNVHSUE
FjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQU
X68CoU40U19EKgh7tvSmk78XaxQwHwYDVR0jBBgwFoAUFC6zF7dYVsuuUAlA5h+v
nYsUwsYwVQYIKwYBBQUHAQEESTBHMCEGCCsGAQUFBzABhhVodHRwOi8vcjMuby5s
ZW5jci5vcmcwIgYIKwYBBQUHMAKGFmh0dHA6Ly9yMy5pLmxlbmNyLm9yZy8wHwYD
VR0RBBgwFoIUdGVzdC5kcm9uZXN0b2NrLnRlY2gwTAYDVR0gBEUwQzAIBgZngQwB
AgEwNwYLKwYBBAGC3xMBAQEwKDAmBggrBgEFBQcCARYaaHR0cDovL2Nwcy5sZXRz
ZW5jcnlwdC5vcmcwggEDBgorBgEEAdZ5AgQCBIH0BIHxAO8AdQC3Pvsk35xNunXy
OcW6WPRsXfxCz3qfNcSeHQmBJe20mQAAAYeOM/DQAAAEAwBGMEQCIB3oyeeSviaF
DDhYIPn6sMbg+XYx4EF72DWx5rMAK+6hAiAY8oI64J3AK8RMeih8OEM/lFZXQ+BA
TQ/hj6lSPM8gZQB2AK33vvp8/xDIi509nB4+GGq0Zyldz7EMJMqFhjTr3IKKAAAB
h44z8P0AAAQDAEcwRQIhAK3zEbu50RZQFomvjA2cxqSA/JKtIraM59tPYVAJvHIn
AiAuFUb2UegSORN2HjwdEXqCmr9/QREAHg8QSBv7ov4DIzANBgkqhkiG9w0BAQsF
AAOCAQEAZAPLw/AJjov/XVESpoO+X7sw5P4V72ps2C/Xjm9JDF+DemWKx58QLk2k
TdRAB/t2QxXTX/mlXEdtmrVG/4wAMvGplC3wYo0mbge5PItmpDz730FadjEBZlVy
+BUfPdfnDPHa/nxkC45KClmR07dF5M+38jsf/elxWTMvnSzaY4ZlImeNm2S7iRfp
3TbAdwagjZOapJSCGnNOJ4dH5IRYqhX5LG+YUYg/T+urDxgolKjqg34YphoiaMDl
NdcJMhpxSM8y4R7NgalIZJDj3WzPcDkiIIX4JQ2aovbPO4Stb2JuyEh3JK3IBvNo
iDQVTkIwJRvdmVLbANsygTA2fSI/tg==
-----END CERTIFICATE-----`
}

func (cur *chuangcacheUploadReq) key(key string) {
	cur.Key = key
	cur.Key = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIW3faaICKIQOHSrfr7H6W1vAAxIyFULq2VxDYUqnORhoAoGCCqGSM49
AwEHoUQDQgAEdaMkPTiA3GteY1+Jt9GMFQNMHTA2o7x0mugEVOOJyEPMqOoV8QGS
zSloRQ0pwHn6Mll9fMCIjv1HDVvrskFpbg==
-----END EC PRIVATE KEY-----`
}

func (cur *chuangcacheUploadReq) chain(_ string) {}
