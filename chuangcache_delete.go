package main

type chuangcacheDeleteReq struct {
	chuangcacheReq

	Id string `json:"ssl_key,omitempty"`
}
