from rocketmq.client import Producer, Message

topic = "TopicTest"


def create_message():
    msg = Message(topic)
    msg.set_keys("hugo")
    msg.set_tags("hugo")
    msg.set_delay_time_level(2)#1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
    msg.set_property("name", "micro services")
    msg.set_body("hello 微服务")

    return msg


def send_message_sync(count):
    producer = Producer("test")
    producer.set_name_server_address("192.168.141.128:9876")

    producer.start()
    for n in range(count):
        msg = create_message()
        ret = producer.send_sync(msg)
        print(f"发送状态:{ret.status}, 消息id:{ret.msg_id}")
    print("消息发送完成")
    producer.shutdown()


if __name__ == "__main__":
    send_message_sync(5)