# 导入 RocketMQ 消费者组资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数，不填则取值环境变量值
terraform import ctyun_rocketmq_group.[导入配置名称] [instance_id],[name][,<region_id>]
# 示例 1: 指定实例 ID 和组名称，region_id 从 provider 或环境变量获取
terraform import ctyun_rocketmq_group.consumer_group_basic d78f56a4-8c3e-4b9d-a123-456789abcdef,tf-consumer-group-basic
# 示例 2: 同时指定实例 ID、组名称和 region_id
terraform import ctyun_rocketmq_group.consumer_group_basic d78f56a4-8c3e-4b9d-a123-456789abcdef,tf-consumer-group-basic,200000002401
