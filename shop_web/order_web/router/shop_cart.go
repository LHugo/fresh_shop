package router

import (
	"github.com/gin-gonic/gin"
	"shop_web/order_web/api/shop_cart"
	"shop_web/order_web/middlewares"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	ShopCartRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())
	{
		ShopCartRouter.GET("", shop_cart.List)          //查询购物车列表
		ShopCartRouter.POST("", shop_cart.New)          //添加商品到购物车
		ShopCartRouter.DELETE("/:id", shop_cart.Delete) //删除购物车条目
		ShopCartRouter.PATCH("/:id", shop_cart.Update)  //修改购物车条目
	}
}
