package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_web/order_web/global"
	"shop_web/order_web/proto"
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

	//订单服务
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil{
		zap.S().Fatal("[InitSrvConn] 连接 【订单服务失败】")
	}
	global.OrderSrvClient = proto.NewOrderClient(orderConn)

	//库存服务
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil{
		zap.S().Fatal("[InitSrvConn] 连接 【库存服务失败】")
	}
	global.InventorySrvClient = proto.NewInventoryClient(inventoryConn)
}
