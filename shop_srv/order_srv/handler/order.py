import json
import time
from datetime import datetime
from rocketmq.client import Producer, TransactionMQProducer, TransactionStatus, Message, SendStatus
from rocketmq.client import ConsumeStatus
from google.protobuf import empty_pb2
import grpc
from order_srv.proto import order_pb2_grpc, order_pb2, goods_pb2, goods_pb2_grpc, inventory_pb2, inventory_pb2_grpc
from loguru import logger
from order_srv.model.models import ShoppingCart, OrderInfo, OrderGoods
from peewee import DoesNotExist
from common.register import consul
from order_srv.settings import settings
from random import Random

local_execute_dict = {}


def generate_order_sn(user_id):
    # 当前时间+user_id+随机数
    return f'{time.strftime("%Y%m%d%H%M%S")}{user_id}{Random().randint(10, 99)}'


def order_timeout(msg):
    msg_body_str = msg.body.decode("utf-8")
    print(f"超时消息接收时间：{datetime.now()},内容：{msg_body_str}")
    msg_body = json.loads(msg_body_str)
    order_sn = msg_body["orderSn"]

    # 1.先查询订单支付状态
    with settings.DB.atomic() as txn:
        try:
            order = OrderInfo.get(OrderInfo.order_sn == order_sn)
            if order.status != "TRADE_SUCCESS":
                order.status = "TRADE_CLOSED"
                order.save()

                # 要给库存服务发送一个归还库存的消息
                msg = Message("order_reback")
                msg.set_keys("hugo")
                msg.set_tags("reback")
                msg.set_body(json.dumps({"orderSn": order_sn}))

                sync_producer = Producer("order_sender")
                sync_producer.set_name_server_address(f"{settings.ROCKETMQ_HOST}:{settings.ROCKETMQ_PORT}")
                sync_producer.start()

                ret = sync_producer.send_sync(msg)
                if ret.status != SendStatus.OK:
                    raise Exception("发送失败")
                sync_producer.shutdown()
        except Exception as e:
            print(e)
            txn.rollback()
            return ConsumeStatus.RECONSUME_LATER

    return ConsumeStatus.CONSUME_SUCCESS


class OrderServicer(order_pb2_grpc.OrderServicer):
    @logger.catch
    def CartItemList(self, request, context):
        # 获取用户的购物车信息
        items = ShoppingCart.select().where(ShoppingCart.user == request.id)
        rsp = order_pb2.CartItemListResponse(total=items.count())
        for item in items:
            item_rsp = order_pb2.ShopCartInfoResponse()
            item_rsp.id = item.id
            item_rsp.userId = item.user
            item_rsp.goodsId = item.goods
            item_rsp.nums = item.nums
            item_rsp.checked = item.checked

            rsp.data.append(item_rsp)

        return rsp

    @logger.catch
    def CreateCartItem(self, request, context):
        # 添加商品到购物车
        existed_items = ShoppingCart.select().where(ShoppingCart.goods == request.goodsId,
                                                    ShoppingCart.user == request.userId)

        # 如果记录已经存在则合并购物车
        if existed_items:
            item = existed_items[0]
            item.nums += request.nums
        else:
            item = ShoppingCart()
            item.user = request.userId
            item.goods = request.goodsId
            item.nums = request.nums
        item.save()

        return order_pb2.ShopCartInfoResponse(id=item.id)

    @logger.catch
    def UpdateCartItem(self, request, context):
        # 更新购物车条目-数量和选中状态
        try:
            item = ShoppingCart.get(ShoppingCart.user == request.userId, ShoppingCart.goods == request.goodsId)
            item.checked = request.checked
            if request.nums:
                item.nums = request.nums
            item.save()
            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("购物车记录不存在")
            return empty_pb2.Empty()

    @logger.catch
    def DeleteCartItem(self, request, context):
        # 删除购物车条目信息
        try:
            item = ShoppingCart.get(ShoppingCart.user == request.userId, ShoppingCart.goods == request.goodsId)
            item.delete_instance()

            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("购物车记录不存在")
            return empty_pb2.Empty()

    @logger.catch
    def check_callback(self, msg):
        msg_body = json.loads(msg.body.decode("utf-8"))
        order_sn = msg_body["orderSn"]

        # 查询本地数据库order_sn的订单是否已经生成
        orders = OrderInfo.select().where(OrderInfo.order_sn == order_sn)
        if orders:
            return TransactionStatus.ROLLBACK
        else:
            return TransactionStatus.COMMIT

    @logger.catch
    def local_execute(self, msg, user_args):
        msg_body = json.loads(msg.body.decode("utf-8"))
        order_sn = msg_body["orderSn"]
        local_execute_dict[order_sn] = {}

        with settings.DB.atomic() as txn:
            goods_ids = []
            goods_nums = {}
            order_amount = 0
            order_goods_list = []
            goods_sell_info = []

            # 获取购物车记录
            for cart_item in ShoppingCart.select().where(ShoppingCart.user == msg_body["userId"],
                                                         ShoppingCart.checked == True):
                goods_ids.append(cart_item.goods)
                goods_nums[cart_item.goods] = cart_item.nums

            if not goods_ids:
                local_execute_dict[order_sn]["code"] = grpc.StatusCode.NOT_FOUND
                local_execute_dict[order_sn]["detail"] = "没有选中结算的商品"
                # 回滚事务消息
                return TransactionStatus.ROLLBACK

            # 连接商品服务
            register = consul.ConsulRegister(settings.CONSUL_HOST, settings.CONSUL_PORT)
            goods_srv_host, goods_srv_port = register.get_host_port(f'Service=="{settings.GOODS_SRV_NAME}"')
            if not goods_srv_host:
                local_execute_dict[order_sn]["code"] = grpc.StatusCode.INTERNAL
                local_execute_dict[order_sn]["detail"] = "商品服务不可用"

                return TransactionStatus.ROLLBACK

            goods_channel = grpc.insecure_channel(f"{goods_srv_host}:{goods_srv_port}")
            goods_stub = goods_pb2_grpc.GoodsStub(goods_channel)

            # 批量获取商品的信息
            try:
                goods_info_rsp = goods_stub.BatchGetGoods(goods_pb2.BatchGoodsIdInfo(id=goods_ids))
            except grpc.RpcError as e:
                local_execute_dict[order_sn]["code"] = grpc.StatusCode.INTERNAL
                local_execute_dict[order_sn]["detail"] = f"商品服务不可用:{str(e)}"

                return TransactionStatus.ROLLBACK
            for goods_info in goods_info_rsp.data:
                order_amount += goods_info.shopPrice * goods_nums[goods_info.id]
                order_goods = OrderGoods(goods=goods_info.id, goods_name=goods_info.name,
                                         goods_image=goods_info.goodsFrontImage, goods_price=goods_info.shopPrice,
                                         nums=goods_nums[goods_info.id])
                order_goods_list.append(order_goods)
                goods_sell_info.append(inventory_pb2.GoodsInvInfo(goodsId=goods_info.id, num=goods_nums[goods_info.id]))

            # 扣减库存
            # 连接库存服务
            inventory_srv_host, inventory_srv_port = register.get_host_port(f'Service=="{settings.INVENTORY_SRV_NAME}"')
            if not inventory_srv_host:
                local_execute_dict[order_sn]["code"] = grpc.StatusCode.INTERNAL
                local_execute_dict[order_sn]["detail"] = "库存服务不可用"

                return TransactionStatus.ROLLBACK

            inv_channel = grpc.insecure_channel(f"{inventory_srv_host}:{inventory_srv_port}")
            inv_stub = inventory_pb2_grpc.InventoryStub(inv_channel)

            try:
                inv_stub.Sell(inventory_pb2.SellInfo(orderSn=order_sn, goodsInfo=goods_sell_info))
            except grpc.RpcError as e:
                local_execute_dict[order_sn]["code"] = grpc.StatusCode.INTERNAL
                local_execute_dict[order_sn]["detail"] = f"扣减库存失败：{str(e)}"
                err_code = e.code()
                # 库存不足，库存服务没有扣减库存，不需要确认事务消息来归还库存，其他状态也没有执行库存扣减，所以也是ROLLBACK操作;
                # 但是如果是未知状态和超时状态需要确认事务消息来进行归还库存，不代表库存一定扣减成功，要进行业务端逻辑确认是否已经扣减库存
                if err_code == grpc.StatusCode.UNKNOWN or err_code == grpc.StatusCode.DEADLINE_EXCEEDED:
                    return TransactionStatus.COMMIT
                else:
                    return TransactionStatus.ROLLBACK

            # 创建订单（如果在这一步突然出现意外情况断掉，要执行事务消息回查逻辑，确认订单是否生成）
            try:
                order = OrderInfo()
                order.user = msg_body["userId"]
                order.order_sn = order_sn
                order.order_amount = order_amount
                order.address = msg_body["userId"]
                order.signer_name = msg_body["name"]
                order.signer_mobile = msg_body["mobile"]
                order.post = msg_body["post"]
                order.save()

                # 批量插入订单商品表
                for order_good in order_goods_list:
                    order_good.order = order.id
                OrderGoods.bulk_create(order_goods_list)

                # 删除购物车记录
                ShoppingCart.delete().where(ShoppingCart.user == msg_body["userId"],
                                            ShoppingCart.checked == True).execute()
                local_execute_dict[order_sn] = {
                    "code": grpc.StatusCode.OK,
                    "detail": "下单成功",
                    "order": {
                        "id": order.id,
                        "orderSn": order_sn,
                        "total": order.order_amount
                    }
                }

                # 发送延时消息
                msg = Message("order_timeout")
                msg.set_delay_time_level(5)  # 设置超时时间为1min
                msg.set_keys("hugo")
                msg.set_tags("cancel")
                msg.set_body(json.dumps({"orderSn": order_sn}))
                sync_producer = Producer("cancel")
                sync_producer.set_name_server_address(f"{settings.ROCKETMQ_HOST}:{settings.ROCKETMQ_PORT}")
                sync_producer.start()

                ret = sync_producer.send_sync(msg)
                if ret.status != SendStatus.OK:
                    raise Exception("发送延时消息失败")
                print(f"发送时间:{datetime.now()}")
                sync_producer.shutdown()

            except Exception as e:
                txn.rollback()
                local_execute_dict[order_sn]["code"] = grpc.StatusCode.INTERNAL
                local_execute_dict[order_sn]["detail"] = f"订单创建失败：{str(e)}"

                return TransactionStatus.COMMIT

        return TransactionStatus.ROLLBACK

    @logger.catch
    def CreateOrder(self, request, context):
        """
        新建订单(逻辑放在local_execute函数里面做)
        1.价格-访问商品服务
        2.库存-访问库存服务
        3.新建订单基本信息表和订单商品信息表
        4.从购物车获取选中商品
        5.从购物车删除已购买的商品
        """

        # 准备一个半事务消息
        producer = TransactionMQProducer("shop_srv", self.check_callback)
        producer.set_name_server_address(f"{settings.ROCKETMQ_HOST}:{settings.ROCKETMQ_PORT}")
        producer.start()
        msg = Message("order_reback")
        msg.set_keys("shop_srv")
        msg.set_tags("order")

        order_sn = generate_order_sn(request.userId)
        msg_body = {
            "orderSn": order_sn,
            "userId": request.userId,
            "address": request.address,
            "name": request.name,
            "mobile": request.mobile,
            "post": request.post
        }
        msg.set_body(json.dumps(msg_body))

        ret = producer.send_message_in_transaction(msg, self.local_execute, user_args=None)
        logger.info(f"发送状态：{ret.status}, 消息id:{ret.msg_id}")
        if ret.status != SendStatus.OK:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details("新建订单失败")
            return order_pb2.OrderInfoResponse()

        while True:
            if order_sn in local_execute_dict:
                context.set_code(local_execute_dict[order_sn]["code"])
                context.set_details(local_execute_dict[order_sn]["detail"])
                producer.shutdown()
                if local_execute_dict[order_sn]["code"] == grpc.StatusCode.OK:
                    return order_pb2.OrderInfoResponse(id=local_execute_dict[order_sn]["order"]["id"],
                                                       orderSn=local_execute_dict[order_sn]["order"]["orderSn"],
                                                       total=local_execute_dict[order_sn]["order"]["total"])
                else:
                    return order_pb2.OrderInfoResponse()
            time.sleep(0.1)

    @logger.catch
    def OrderList(self, request, context):
        # 订单列表
        rsp = order_pb2.OrderListResponse()
        orders = OrderInfo.select()
        if request.userId:
            orders = orders.where(OrderInfo.user == request.userId)
        rsp.total = orders.count()

        # 分页
        per_page_nums = request.pagePerNums if request.pagePerNums else 10
        start = per_page_nums * (request.pages - 1) if request.pages else 0
        orders = orders.limit(per_page_nums).offset(start)

        for order in orders:
            tmp_rsp = order_pb2.OrderInfoResponse()

            tmp_rsp.id = order.id
            tmp_rsp.userId = order.user
            tmp_rsp.orderSn = order.order_sn
            tmp_rsp.payType = order.pay_type
            tmp_rsp.status = order.status
            tmp_rsp.post = order.post
            tmp_rsp.total = order.order_amount
            tmp_rsp.address = order.address
            tmp_rsp.name = order.signer_name
            tmp_rsp.mobile = order.signer_mobile
            tmp_rsp.add_time = order.add_time.strftime('%Y-%m-%d %H:%M:%S')

            rsp.data.append(tmp_rsp)
        return rsp

    @logger.catch
    def OrderDetail(self, request, context):
        # 订单详情
        rsp = order_pb2.OrderInfoDetailResponse()
        try:
            if request.userId:
                order = OrderInfo.get(OrderInfo.id == request.id, OrderInfo.user == request.userId)
            else:
                order = OrderInfo.get(OrderInfo.id == request.id)

            rsp.orderInfo.id = order.id
            rsp.orderInfo.userId = order.user
            rsp.orderInfo.orderSn = order.order_sn
            rsp.orderInfo.payType = order.pay_type
            rsp.orderInfo.status = order.status
            rsp.orderInfo.post = order.post
            rsp.orderInfo.address = order.address
            rsp.orderInfo.total = order.order_amount
            rsp.orderInfo.name = order.signer_name
            rsp.orderInfo.mobile = order.signer_mobile

            order_goods = OrderGoods.select().where(OrderGoods.order == order.id)
            for order_good in order_goods:
                order_goods_rsp = order_pb2.OrderItemResponse()

                order_goods_rsp.goodsId = order_good.goods
                order_goods_rsp.goodsName = order_good.goods_name
                order_goods_rsp.goodsImage = order_good.goods_image
                order_goods_rsp.goodsPrice = float(order_good.goods_price)
                order_goods_rsp.nums = order_good.nums

                rsp.data.append(order_goods_rsp)

            return rsp
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("订单记录不存在")
            return rsp

    @logger.catch
    def UpdateOrderStatus(self, request, context):
        # 更新订单的支付状态
        OrderInfo.update(status=request.status).where(OrderInfo.order_sn == request.orderSn).execute()
        return empty_pb2.Empty()
