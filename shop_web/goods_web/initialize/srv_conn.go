package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_web/goods_web/global"
	"shop_web/goods_web/proto"
)


func InitSrvConn(){
	//集成consul服务中心注册功能和负载均衡功能的grpc拨号连接srv
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		)
	if err != nil{
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}
	global.GoodsSrvClient = proto.NewGoodsClient(userConn)
}
