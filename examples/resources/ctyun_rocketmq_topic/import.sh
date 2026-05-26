# 导入 RocketMQ 主题资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数，不填则取值环境变量值
terraform import ctyun_rocketmq_topic.[导入配置名称] [instance_id],[name][,<region_id>]
# 示例 1: 指定实例 ID 和主题名称，region_id 从 provider 或环境变量获取
terraform import ctyun_rocketmq_topic.normal_topic d78f56a4-8c3e-4b9d-a123-456789abcdef,tf-normal-topic
# 示例 2: 同时指定实例 ID、主题名称和 region_id
terraform import ctyun_rocketmq_topic.normal_topic d78f56a4-8c3e-4b9d-a123-456789abcdef,tf-normal-topic,200000002401
