package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_web/user_operation_web/global"
	"shop_web/user_operation_web/proto"
)


func InitSrvConn(){
	//集成consul服务中心注册功能和负载均衡功能的grpc拨号连接srv
	consulInfo := global.ServerConfig.ConsulInfo

	//商品服务
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		)
	if err != nil{
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务失败】")
	}
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	//用户操作服务
	userOperationConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserOperationSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil{
		zap.S().Fatal("[InitSrvConn] 连接 【用户操作服务失败】")
	}
	global.FavouriteClient = proto.NewFavouriteClient(userOperationConn)
	global.AddressClient = proto.NewAddressClient(userOperationConn)
	global.MessageClient = proto.NewMessageClient(userOperationConn)
}
