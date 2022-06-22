import time

from rocketmq.client import TransactionMQProducer, Message, TransactionStatus

topic = "TopicTest"


def create_message():
    msg = Message(topic)
    msg.set_keys("hugo")
    msg.set_tags("hugo")
    msg.set_property("name", "micro services")
    msg.set_body("hello 微服务")

    return msg


def check_callback(msg):
    # 消息回查
    # TransactionStatus.COMMIT, TransactionStatus.ROLLBACK, TransactionStatus.UNKNOWN
    print(f"事务消息回查:{msg.body.decode('utf-8')}")
    return TransactionStatus.COMMIT


def local_execute(msg, user_args):
    # 基于事务消息发送成功后回调再执行的事务逻辑
    print("执行本地事务逻辑")
    return TransactionStatus.UNKNOWN


def send_transaction_message(count):
    producer = TransactionMQProducer("test", check_callback)
    producer.set_name_server_address("192.168.141.128:9876")

    producer.start()
    for n in range(count):
        msg = create_message()
        ret = producer.send_message_in_transaction(msg, local_execute, None)
        print(f"发送状态:{ret.status}, 消息id:{ret.msg_id}")
    print("消息发送完成")
    while True:
        time.sleep(3600)


if __name__ == "__main__":
    send_transaction_message(1)