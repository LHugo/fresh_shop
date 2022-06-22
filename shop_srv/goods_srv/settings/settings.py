import json
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
    "NameSpace": "c8d78d50-ae0f-4e5e-8918-824a5cd88f9d",
    "User": "nacos",
    "Password": "nacos",
    "DataId": "goods-srv.json",
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

# 服务相关配置
SERVICE_NAME = data["name"]
SERVICE_TAGS = data["tags"]

# 数据库相关配置
DB = ReconnectMysqlDatabase(data["mysql"]["db"], host=data["mysql"]["host"], port=data["mysql"]["port"],
                            user=data["mysql"]["user"], password=data["mysql"]["password"])
