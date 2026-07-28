# 导入云硬盘备份策略绑定存储库
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ebs_backup_policy_bind_repo.[导入配置名称] [policy_id],[repository_id],<region_id>
# 示例
terraform import ctyun_ebs_backup_policy_bind_repo.example b4d9a692-cd51-4a95-9769-492e237f148c,c4d9a692-cd51-4a95-9769-492e237f148c,bb9fdb42056f11eda1610242ac110002