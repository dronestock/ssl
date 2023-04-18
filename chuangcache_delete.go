package main

type chuangcacheDeleteReq struct {
	chuangcacheReq

	Key string `json:"ssl_key,omitempty"`
}
