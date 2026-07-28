# 导入ECS备份策略
# [] 标记的参数为必填参数
# <> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ecs_backup_policy.example [id],<region_id>
# 示例
terraform import ctyun_ecs_backup_policy.example 12345678-1234-1234-1234-123456789012,bb9fdb42056f11eda1610242ac110002