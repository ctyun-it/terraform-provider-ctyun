# 导入Redis实例的账户
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_redis_account.[导入配置名称] [instance_id],[name],<region_id>
# 示例
terraform import ctyun_redis_account.test b9cc9df4e96144acb7c051262112bab4,myaccount,<region_id>