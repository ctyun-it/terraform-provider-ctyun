# 导入云主机组关联
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ecs_affinity_group_association.[导入配置名称] [instance_id],[group_id],<region_id>
# 示例
terraform import ctyun_ecs_affinity_group_association.ecs_affinity_group_association_example inst-12345678,ag-87654321,region-11111111