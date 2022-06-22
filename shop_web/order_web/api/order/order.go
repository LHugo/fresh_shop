package order

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"net/http"
	"shop_web/order_web/api"
	"shop_web/order_web/forms"
	"shop_web/order_web/global"
	"shop_web/order_web/models"
	"shop_web/order_web/proto"
	"strconv"
)

func List(c *gin.Context) {
	//订单列表
	userId, _ := c.Get("userId")
	claims, _ := c.Get("claims")

	request := proto.OrderFilterRequest{
		UserId:      0,
		Pages:       0,
		PagePerNums: 0,
	}
	//若是管理员则返回所有订单
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1{
		request.UserId = int32(userId.(uint))
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	request.Pages = int32(offset)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))
	request.PagePerNums = int32(limit)

	rsp, err := global.OrderSrvClient.OrderList(context.Background(), &request)
	if err != nil{
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	reMap := gin.H{
		"total": rsp.Total,
	}
	orderList := make([]interface{}, 0)

	for _, item := range rsp.Data {
		tmpMap := map[string]interface{}{}

		tmpMap["id"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["pay_type"] = item.PayType
		tmpMap["user"] = item.UserId
		tmpMap["post"] = item.Post
		tmpMap["total"] = item.Total
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		tmpMap["order_sn"] = item.OrderSn
		tmpMap["id"] = item.Id
		tmpMap["add_time"] = item.AddTime

		orderList = append(orderList, tmpMap)
	}
	reMap["data"] = orderList
	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context) {
	orderForm := forms.CreateOrderForm{}
	if err := c.ShouldBindJSON(&orderForm); err != nil{
		api.HandleValidatorError(c, err)
	}
	userId, _ := c.Get("userId")
	rsp, err := global.OrderSrvClient.CreateOrder(context.WithValue(context.Background(), "ginContext", c), &proto.OrderRequest{
		UserId: int32(userId.(uint)),
		Name: orderForm.Name,
		Mobile: orderForm.Mobile,
		Address: orderForm.Address,
		Post: orderForm.Post,
	})
	if err != nil{
		zap.S().Errorw("新建订单失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	//生成支付宝的支付url
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
	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AlipayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AlipayInfo.ReturnURL
	p.Subject = "生鲜商城订单-"+rsp.OrderSn
	p.OutTradeNo = rsp.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil{
		zap.S().Errorw("生成支付宝url失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
		"alipay_url":url.String(),
	})
}

func Detail(c *gin.Context) {
	id := c.Param("id")
	userId, _ := c.Get("userId")
	i, err := strconv.Atoi(id)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{
			"msg":"url格式出错",
		})
		return
	}
	claims, _ := c.Get("claims")
	request := proto.OrderRequest{
		Id: int32(i),
	}

	//若是管理员则返回所有订单
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1{
		request.UserId = int32(userId.(uint))
	}
	rsp, err := global.OrderSrvClient.OrderDetail(context.Background(), &request)
	if err != nil{
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	reMap := gin.H{}
	reMap["id"] = rsp.OrderInfo.Id
	reMap["status"] = rsp.OrderInfo.Status
	reMap["user"] = rsp.OrderInfo.UserId
	reMap["post"] = rsp.OrderInfo.Post
	reMap["total"] = rsp.OrderInfo.Total
	reMap["address"] = rsp.OrderInfo.Address
	reMap["name"] = rsp.OrderInfo.Name
	reMap["mobile"] = rsp.OrderInfo.Mobile
	reMap["pay_type"] = rsp.OrderInfo.PayType
	reMap["order_sn"] = rsp.OrderInfo.OrderSn

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Data{
		tmpMap := gin.H{
			"id": item.GoodsId,
			"name": item.GoodsName,
			"price": item.GoodsPrice,
			"nums": item.Nums,
			"image": item.GoodsImage,
		}
		goodsList = append(goodsList, tmpMap)
	}
	reMap["goods"] = goodsList

	//生成支付宝链接
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
	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AlipayInfo.NotifyURL
	p.ReturnURL = global.ServerConfig.AlipayInfo.ReturnURL
	p.Subject = "生鲜商城订单-"+rsp.OrderInfo.OrderSn
	p.OutTradeNo = rsp.OrderInfo.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.OrderInfo.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil{
		zap.S().Errorw("生成支付宝url失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":err.Error(),
		})
		return
	}
	reMap["alipay_url"] = url.String()

	c.JSON(http.StatusOK, reMap)
}
