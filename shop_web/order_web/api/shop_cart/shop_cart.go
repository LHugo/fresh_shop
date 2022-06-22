package shop_cart

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"shop_web/order_web/api"
	"shop_web/order_web/forms"
	"shop_web/order_web/global"
	"shop_web/order_web/proto"
	"strconv"
)

func List(c *gin.Context){
	//获取购物车商品
	userId, _ := c.Get("userId")
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("[List] 查询 【购物车列表】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	ids := make([]int32, 0)
	for _, item := range rsp.Data{
		ids = append(ids, item.GoodsId)
	}
	if len(ids) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"total":0,
		})
		return
	}

	//请求商品服务获取商品信息
	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err!=nil{
		zap.S().Errorw("[List] 批量查询【商品列表】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		for _, good := range goodsRsp.Data{
			if good.Id == item.GoodsId {
				tmpMap := map[string]interface{}{}
				tmpMap["id"] = item.Id
				tmpMap["goods_id"] = item.GoodsId
				tmpMap["good_name"] = good.Name
				tmpMap["good_image"] = good.GoodsFrontImage
				tmpMap["good_price"] = good.ShopPrice
				tmpMap["nums"] = item.Nums
				tmpMap["checked"] = item.Checked

				goodsList = append(goodsList, tmpMap)
			}
		}
	}
	reMap["data"] = goodsList
	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context){
	//添加商品到购物车
	itemForm := forms.ShopCartItemForm{}
	if err := c.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	//检查所添加的商品是否存在
	_, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[List] 查询【商品信息】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	//检查所添加商品的库存是否足够扣减
	invRsp, err := global.InventorySrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[List] 查询【库存信息】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	if invRsp.Num < itemForm.Nums {
		c.JSON(http.StatusBadRequest, gin.H{
			"nums":"库存不足",
		})
		return
	}

	userId, _ := c.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId: int32(userId.(uint)),
		Nums: itemForm.Nums,
	})
	if err != nil {
		zap.S().Errorw("添加到购物车失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{
			"msg":"url格式出错",
		})
		return
	}

	userId, _ := c.Get("userId")
	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId: int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil{
		zap.S().Errorw("删除购物车记录失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.Status(http.StatusOK)
}

func Update(c *gin.Context) {
	//修改购物车信息
	itemForm := forms.ShopCartItemUpdateForm{}
	if err := c.ShouldBindJSON(&itemForm); err != nil{
		api.HandleValidatorError(c, err)
		return
	}

	//获取要修改的商品id
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{
			"msg":"url格式出错",
		})
		return
	}

	//获取要修改购物车信息的用户的id
	userId, _ := c.Get("userId")
	request := proto.CartItemRequest{
		UserId: int32(userId.(uint)),
		GoodsId: int32(i),
		Nums: itemForm.Nums,
		Checked: false,
	}
	if itemForm.Checked != nil{
		request.Checked = *itemForm.Checked
	}
	_, err = global.OrderSrvClient.UpdateCartItem(context.Background(), &request)
	if err != nil{
		zap.S().Errorw("修改购物车信息失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.Status(http.StatusOK)
}
