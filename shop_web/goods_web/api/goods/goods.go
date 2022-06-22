package goods

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"shop_web/goods_web/api"
	"shop_web/goods_web/forms"
	"shop_web/goods_web/global"
	"shop_web/goods_web/proto"
	"strconv"
)


func List(c *gin.Context) {
	//获取商品列表
	request := &proto.GoodsFilterRequest{}
	priceMin, _ := strconv.Atoi(c.DefaultQuery("pmin", "0"))
	request.PriceMin = int32(priceMin)

	priceMax, _ := strconv.Atoi(c.DefaultQuery("pmax", "0"))
	request.PriceMax = int32(priceMax)

	isHot, _ := strconv.Atoi(c.DefaultQuery("hot", "0"))
	if isHot == 1 {
		request.IsHot = true
	}

	isNew, _ := strconv.Atoi(c.DefaultQuery("new", "0"))
	if isNew == 1 {
		request.IsNew = true
	}

	isTab, _ := strconv.Atoi(c.DefaultQuery("tab", "0"))
	if isTab == 1 {
		request.IsTab = true
	}

	categoryId, _ := strconv.Atoi(c.DefaultQuery("cd", "0"))
	request.TopCategory = int32(categoryId)

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	request.Pages = int32(offset)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))
	request.PagePerNums = int32(limit)

	keywords := c.DefaultQuery("keywords", "")
	request.KeyWords = keywords

	brandId, _ := strconv.Atoi(c.DefaultQuery("bd", "0"))
	request.Brand = int32(brandId)

	//请求商品的service服务
	r, err := global.GoodsSrvClient.GoodsList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	reMap := map[string]interface{}{
		"total": r.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, value := range r.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":             value.Id,
			"name":           value.Name,
			"goods_brief":    value.GoodsBrief,
			"desc":           value.GoodsDesc,
			"transport_free": value.TransportFree,
			"images":         value.Images,
			"desc_images":    value.DescImages,
			"front_image":    value.GoodsFrontImage,
			"shop_price":     value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	reMap["data"] = goodsList

	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}
	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		TransportFree:   *goodsForm.TransportFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, rsp)
}

func Details(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	r, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	rsp := map[string]interface{}{
		"id":             r.Id,
		"name":           r.Name,
		"goods_brief":    r.GoodsBrief,
		"desc":           r.GoodsDesc,
		"transport_free": r.TransportFree,
		"images":         r.Images,
		"desc_images":    r.DescImages,
		"front_image":    r.GoodsFrontImage,
		"shop_price":     r.ShopPrice,
		"category": map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
			"logo": r.Brand.Logo,
		},
		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"on_sale": r.OnSale,
	}
	c.JSON(http.StatusOK, rsp)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
	return
}

func Stocks(c *gin.Context) {
	id := c.Param("id")
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	//TODO 商品库存
	return
}

func UpdateStatus(c *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := c.ShouldBindJSON(&goodsStatusForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})

}

func Update(c *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		TransportFree:   *goodsForm.TransportFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
