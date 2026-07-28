# 导入Redis的备份
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_rabbitmq_queue.[导入配置名称] [instance_id],[name],<region_id>
# 示例
terraform import ctyun_rabbitmq_queue.queue_example b9cc9df4e96144acb7c051262112bab4,YYYYMMDDHHMMSS,<region_id>