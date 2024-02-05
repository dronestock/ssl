package chuangcache

type DeleteReq struct {
	Request

	Key string `json:"ssl_key,omitempty"`
}
