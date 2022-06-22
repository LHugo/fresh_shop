from user_operation_srv.proto import message_pb2, message_pb2_grpc
from loguru import logger
from user_operation_srv.model.models import Messages
from peewee import DoesNotExist
from google.protobuf import empty_pb2
import grpc


class MessageServicer(message_pb2_grpc.MessageServicer):
    @logger.catch
    def MessageList(self, request: message_pb2.MessageRequest, context):
        # 获取分类列表
        rsp = message_pb2.MessageListResponse()
        messages = Messages.select()
        if request.userId:
            messages = messages.where(Messages.user == request.userId)

        rsp.total = messages.count()
        for message in messages:
            brand_rsp = message_pb2.MessageResponse()

            brand_rsp.id = message.id
            brand_rsp.userId = message.user
            brand_rsp.messageType = message.message_type
            brand_rsp.subject = message.subject
            brand_rsp.text = message.text
            brand_rsp.file = message.file

            rsp.data.append(brand_rsp)

        return rsp

    @logger.catch
    def CreateMessage(self, request: message_pb2.MessageRequest, context):
        message = Messages()

        message.user = request.userId
        message.message_type = request.messageType
        message.subject = request.subject
        message.text = request.text
        message.file = request.file

        message.save()

        rsp = message_pb2.MessageResponse()
        rsp.id = message.id
        rsp.messageType = message.message_type
        rsp.subject = message.subject
        rsp.text = message.text
        rsp.file = message.file

        return rsp

    @logger.catch
    def DeleteMessage(self, request: message_pb2.MessageRequest, context):
        try:
            user_message = Messages.get(Messages.user == request.userId, Messages.id == request.id)
            user_message.delete_instance(permanently=True)

            return empty_pb2.Empty()
        except DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('记录不存在')
            return empty_pb2.Empty()