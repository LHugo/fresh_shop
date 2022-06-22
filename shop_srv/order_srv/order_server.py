import argparse
import os.path
import signal
import socket
import sys
from concurrent import futures
import grpc
from rocketmq.client import PushConsumer
BASE_DIR = os.path.dirname(os.path.abspath(os.path.dirname(__file__)))
sys.path.insert(0, BASE_DIR)
from order_srv.proto import order_pb2_grpc
from order_srv.handler.order import OrderServicer, order_timeout
from loguru import logger
from common.grpc_health.v1 import health_pb2_grpc
from common.grpc_health.v1 import health
from common.register import consul
from order_srv.settings import settings
from functools import partial
import uuid


# 拦截器
class LogInterceptor(grpc.ServerInterceptor):
    def intercept_service(self, continuation, handler_call_details):
        # print("请求开始")
        rsp = continuation(handler_call_details)
        return rsp


# 获取空闲端口
def get_free_tcp_port():
    tcp = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    tcp.bind(("", 0))
    _, port = tcp.getsockname()
    tcp.close()
    return port


# 使用ctrl+c退出服务进程
def on_exit(signal, frame, service_id):
    register = consul.ConsulRegister(settings.CONSUL_HOST, settings.CONSUL_PORT)
    logger.info(f"注销{service_id}服务")
    register.deregister(service_id=service_id)
    logger.info("服务注销成功")
    sys.exit(0)


def server():
    parser = argparse.ArgumentParser()
    parser.add_argument('--ip', nargs="?", type=str, default="192.168.141.128", help="binding ip")
    parser.add_argument('--port', nargs="?", type=int, default=0, help="listening port")
    args = parser.parse_args()
    if args.port == 0:
        port = get_free_tcp_port()
    else:
        port = args.port
    logger.add("logs/order_srv_{time}.log")
    # 拦截器实例化
    interceptor = LogInterceptor()
    # 1.实例化server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10), interceptors=(interceptor,))
    # 2.注册库存服务到server中
    order_pb2_grpc.add_OrderServicer_to_server(OrderServicer(), server)
    # 3. 注册健康检查
    health_pb2_grpc.add_HealthServicer_to_server(health.HealthServicer(), server)
    # 4.启动server
    server.add_insecure_port(f"{args.ip}:{port}")

    service_id = str(uuid.uuid1())
    signal.signal(signal.SIGINT, partial(on_exit, service_id=service_id))
    signal.signal(signal.SIGTERM, partial(on_exit, service_id=service_id))

    logger.info(f"启动服务：{args.ip}:{port}")
    server.start()

    logger.info(f"服务注册开始")
    register = consul.ConsulRegister(settings.CONSUL_HOST, settings.CONSUL_PORT)
    if not register.register(name=settings.SERVICE_NAME, id=service_id, address=args.ip, port=port,
                             tags=settings.SERVICE_TAGS, check=None):
        logger.info(f"服务注册失败")

        sys.exit(0)
    logger.info(f"服务注册成功")

    # 监听超时订单消息
    consumer = PushConsumer("shop_order")
    consumer.set_name_server_address(f"{settings.ROCKETMQ_HOST}:{settings.ROCKETMQ_PORT}")
    consumer.subscribe("order_timeout", order_timeout)
    consumer.start()

    server.wait_for_termination()
    consumer.shutdown()


if __name__ == "__main__":
    settings.client.add_config_watcher(settings.NACOS["DataId"], settings.NACOS["Group"], settings.config_info_update)
    server()
