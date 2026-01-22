# 导入云硬盘快照策略
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ebs_snapshot_policy.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_ebs_snapshot_policy.example b4d9a692-cd51-4a95-9769-492e237f148c,bb9fdb42056f11eda1610242ac110002