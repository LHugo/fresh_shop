from user_operation_srv.proto import favourite_pb2, favourite_pb2_grpc, address_pb2_grpc, address_pb2, message_pb2_grpc, message_pb2
from user_operation_srv.settings import settings
import consul
import grpc


class UserOpTest:
    def __init__(self):
        # 连接grpc服务器
        c = consul.Consul(settings.CONSUL_HOST, settings.CONSUL_PORT)
        services = c.agent.services()
        ip = ""
        port = ""
        for key, value in services.items():
            if value["Service"] == settings.SERVICE_NAME:
                ip = value["Address"]
                port = value["Port"]
                break
        if not ip:
            raise Exception("订单服务不可用")
        channel = grpc.insecure_channel(f"{ip}:{port}")
        self.stub = message_pb2_grpc.MessageStub(channel)

    def create_message(self):
        rsp = self.stub.CreateMessage(message_pb2.MessageRequest(userId=11, messageType=1, subject="test", text="test"))
        print(rsp)

    def delete_message(self):
        rsp = self.stub.DeleteMessage(message_pb2.MessageRequest(userId=11, id=1))
        print(rsp)


if __name__ == "__main__":
    userop = UserOpTest()
    # order.create_cart_item()
    userop.delete_message()
    # order.order_list()
    # order.cart_item_list()
