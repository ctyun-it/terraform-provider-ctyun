# 导入ECS备份策略绑定存储库
# [] 标记的参数为必填参数
# <> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ecs_backup_policy_bind_repo.example [policy_id],[repository_id],<region_id>
# 示例
terraform import ctyun_ecs_backup_policy_bind_repo.example 12345678-1234-1234-1234-123456789012,87654321-4321-4321-4321-210987654321,bb9fdb42056f11eda1610242ac110002