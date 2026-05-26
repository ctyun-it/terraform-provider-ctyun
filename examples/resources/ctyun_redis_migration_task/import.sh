# 导入Redis的迁移任务
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_redis_migration_task.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_redis_migration_task.test 0d3290ecfbfc4fd092c0ae3e4b12a9fb,<region_id>