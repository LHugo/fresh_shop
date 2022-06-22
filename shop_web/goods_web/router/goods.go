package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop_web/goods_web/middlewares"

	"shop_web/goods_web/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	zap.S().Info("正在配置商品相关的url")
	{
		GoodsRouter.GET("", goods.List)//商品列表
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)//添加商品
		GoodsRouter.GET("/:id", goods.Details)//获取商品详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)//删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)//获取商品库存信息
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)//更新商品详细信息
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)//更新商品状态
	}
}
