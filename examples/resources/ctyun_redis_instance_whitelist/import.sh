# 导入Redis实例白名单
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_redis_instance_whitelist.[导入配置名称] [instance_id],[name],<region_id>
# 示例
terraform import ctyun_redis_instance_whitelist.test b9cc9df4e96144acb7c051262112bab4,white,<region_id>