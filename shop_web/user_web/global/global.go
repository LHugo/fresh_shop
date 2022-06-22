package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_web/user_web/config"
	"shop_web/user_web/proto"
)

var (
	ServerConfig  = &config.ServerConfig{}
	NacosConfig  = &config.NacosConfig{}
	Trans ut.Translator
	UserSrvClient proto.UserClient
)
