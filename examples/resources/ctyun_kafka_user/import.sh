# 导入kafka用户
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_kafka_user.[导入配置名称] [instance_id],[name],<region_id>
# 示例
terraform import ctyun_kafka_user.kafka_user_example instance-123456,user_test,region-11111111