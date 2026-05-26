# 导入Redis实例和弹性IP的绑定关系
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_redis_association_eip.[导入配置名称] [instance_id],[eip_id],<region_id>
# 示例
terraform import ctyun_redis_association_eip.test b9cc9df4e96144acb7c051262112bab4,eip-xxxxxxxx,<region_id>