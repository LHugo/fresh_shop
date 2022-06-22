import json
import redis
from loguru import logger
import nacos
from playhouse.pool import PooledMySQLDatabase
from playhouse.shortcuts import ReconnectMixin


class ReconnectMysqlDatabase(ReconnectMixin, PooledMySQLDatabase):
    pass


def config_info_update(args):
    print(f"配置文件产生变化:{args}")
    # updated_data = json.loads(args["content"])
    # return updated_data


# nacos配置
NACOS = {
    "Host": "192.168.141.128",
    "Port": 8848,
    "NameSpace": "d8dbecac-b88c-47d9-81df-a95e68306016",
    "User": "nacos",
    "Password": "nacos",
    "DataId": "order-srv.json",
    "Group": "dev"
}

client = nacos.NacosClient(f'{NACOS["Host"]}:{NACOS["Port"]}', namespace=NACOS["NameSpace"], username=NACOS["User"],
                           password=NACOS["Password"])
raw_data = client.get_config(NACOS["DataId"], NACOS["Group"])
data = json.loads(raw_data)
logger.info(data)


# consul配置
CONSUL_HOST = data["consul"]["host"]
CONSUL_PORT = data["consul"]["port"]

ROCKETMQ_HOST = data["rocketmq"]["host"]
ROCKETMQ_PORT = data["rocketmq"]["port"]

# 服务相关配置
SERVICE_NAME = data["name"]
SERVICE_TAGS = data["tags"]

# REDIS_HOST = data["redis"]["host"]
# REDIS_PORT = data["redis"]["port"]
# REDIS_DB = data["redis"]["db"]

GOODS_SRV_NAME = data["goods_srv"]["name"]
INVENTORY_SRV_NAME = data["inventory_srv"]["name"]

# # 配置redis连接池
# pool = redis.ConnectionPool(host=REDIS_HOST, port=REDIS_PORT, db=REDIS_DB)
# REDIS_CLIENT = redis.StrictRedis(connection_pool=pool)

# 数据库相关配置
DB = ReconnectMysqlDatabase(data["mysql"]["db"], host=data["mysql"]["host"], port=data["mysql"]["port"],
                            user=data["mysql"]["user"], password=data["mysql"]["password"])
