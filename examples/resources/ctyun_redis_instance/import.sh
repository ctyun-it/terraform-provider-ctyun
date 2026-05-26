# 导入Redis实例
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_redis_instance.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_redis_instance.test b9cc9df4e96144acb7c051262112bab4,<region_id>