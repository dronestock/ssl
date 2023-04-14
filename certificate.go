package main

type certificate struct {
	// 域名
	Domain string `json:"domain,omitempty" validate:"required_without=Domains"`
	// 域名列表
	Domains []string `json:"domains,omitempty" validate:"required_without=Domain,dive,domain"`
	// 类型
	Type string `json:"type,omitempty" validate:"required,oneof=ali"`
}
