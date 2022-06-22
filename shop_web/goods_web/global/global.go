package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_web/goods_web/config"
	"shop_web/goods_web/proto"
)

var (
	ServerConfig  = &config.ServerConfig{}
	NacosConfig  = &config.NacosConfig{}
	Trans ut.Translator
	GoodsSrvClient proto.GoodsClient
)
