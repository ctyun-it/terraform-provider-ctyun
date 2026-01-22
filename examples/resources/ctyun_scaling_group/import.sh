# 导入伸缩组
# [] 标记的参数为必填参数
# <> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_scaling_group.[导入配置名称] [id],[vpc_id],<region_id>
# 示例
terraform import ctyun_scaling_group.scaling_group_example 123456789,subnet-123456,<region-123456>