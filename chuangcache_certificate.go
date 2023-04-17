package main

type (
	chuangcacheUploadReq struct {
		*chuangcacheReq

		Title       string `json:"ssl_title,omitempty"`
		Key         string `json:"private_key,omitempty"`
		Certificate string `json:"certificate,omitempty"`
	}

	chuangcacheUploadRsp struct {
		Id string `json:"ssl_key,omitempty"`
	}
)
