package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_web/user_operation_web/config"
	"shop_web/user_operation_web/proto"
)

var (
	ServerConfig  = &config.ServerConfig{}
	NacosConfig  = &config.NacosConfig{}
	Trans ut.Translator
	GoodsSrvClient proto.GoodsClient
	MessageClient proto.MessageClient
	AddressClient proto.AddressClient
	FavouriteClient proto.FavouriteClient
)
