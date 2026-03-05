# 导入数据盘
#[ ] 标记的参数为必填参数
#< > 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ecs_data_volume.[导入配置名称] [instance_id],<region_id>
# 示例
terraform import ctyun_ecs_data_volume.data_volume_example xxxx-ssss-12345678,<region_id>