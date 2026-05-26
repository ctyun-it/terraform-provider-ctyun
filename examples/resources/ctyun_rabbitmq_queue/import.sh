# 导入RabbitMQ队列
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_rabbitmq_queue.[导入配置名称] [instance_id],[vhost],[name],<region_id>
# 示例
terraform import ctyun_rabbitmq_queue.queue_example inst123456,myvhost,myqueue,<region_id>