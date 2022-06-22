package pay

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"net/http"
	"shop_web/order_web/global"
	"shop_web/order_web/proto"
)

func Notify(c *gin.Context)  {
	//支付宝回调通知
	client, err := alipay.New(global.ServerConfig.AlipayInfo.AppID, global.ServerConfig.AlipayInfo.PrivateKey, false)
	if err != nil{
		zap.S().Errorw("实例化支付宝失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(global.ServerConfig.AlipayInfo.AliPublicKey)
	if err != nil{
		zap.S().Errorw("加载支付宝公钥失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":err.Error(),
		})
		return
	}

	noti, err := client.GetTradeNotification(c.Request)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status: string(noti.TradeStatus),
	})
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.String(http.StatusOK, "success")
}

