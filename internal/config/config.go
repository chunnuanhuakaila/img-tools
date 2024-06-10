package config

import "github.com/zeromicro/go-zero/rest"

type Ask struct {
	AK string `json:"ak"`
	SK string `json:"sk"`
}

type Config struct {
	rest.RestConf
	ZijieAI []Ask
	BaiduAI []Ask
}
