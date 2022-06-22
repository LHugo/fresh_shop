from user_operation_srv.proto import favourite_pb2, favourite_pb2_grpc
from user_operation_srv.model.models import Favourite
from loguru import logger
from google.protobuf import empty_pb2
from peewee import DoesNotExist
import grpc


class FavouriteServicer(favourite_pb2_grpc.FavouriteServicer):
    @logger.catch
    def GetFavList(self, request: favourite_pb2.FavRequest, context):
        rsp = favourite_pb2.FavListResponse()
        user_favs = Favourite.select()
        if request.userId:
            user_favs = user_favs.where(Favourite.user == request.userId)
        if request.goodsId:
            user_favs = user_favs.where(Favourite.goods == request.goodsId)

        rsp.total = user_favs.count()
        for user_fav in user_favs:
            user_fav_rsp = favourite_pb2.FavResponse()
            user_fav_rsp.userId = user_fav.user
            user_fav_rsp.goodsId = user_fav.goods

            rsp.data.append(user_fav_rsp)

        return rsp

    @logger.catch
    def AddFav(self, request: favourite_pb2.FavRequest, context):
        # 添加收藏
        # TODO 查询商品是否存在
        user_fav = Favourite()
        user_fav.user = request.userId
        user_fav.goods = request.goodsId
        user_fav.save(force_insert=True)

        return empty_pb2.Empty()

    @logger.catch
    def DeleteFav(self, request: favourite_pb2.FavRequest, context):
        try:
            user_fav = Favourite.get(Favourite.user == request.userId, Favourite.goods == request.goodsId)
            user_fav.delete_instance(permanently=True)

            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('记录不存在')
            return empty_pb2.Empty()

    @logger.catch
    def GetFavDetail(self, request: favourite_pb2.FavRequest, context):
        try:
            Favourite.get(Favourite.user == request.userId, Favourite.goods == request.goodsId)
            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('记录不存在')
            return empty_pb2.Empty()