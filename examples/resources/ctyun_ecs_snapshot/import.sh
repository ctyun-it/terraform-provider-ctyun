# 导入云主机快照
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ecs_snapshot.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_ecs_snapshot.snapshot_example snap-12345678,<region_id>