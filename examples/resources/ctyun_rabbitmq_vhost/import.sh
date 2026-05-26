# 导入RabbitMQ虚拟主机
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_rabbitmq_vhost.[导入配置名称] [instance_id],[name],<region_id>
# 示例
terraform import ctyun_rabbitmq_vhost.vhost_example inst123456,myvhost,<region_id>