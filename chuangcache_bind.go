package main

type chuangcacheBindReq struct {
	chuangcacheReq

	Id     string `json:"ssl_key,omitempty"`
	Domain string `json:"domain_id,omitempty"`
}
