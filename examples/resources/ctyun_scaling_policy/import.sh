# 导入伸缩策略
# [] 标记的参数为必填参数
# <> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_scaling_policy.[导入配置名称] [id],[group_id],[policy_type],<region_id>
# 示例
terraform import ctyun_scaling_policy.scaling_policy_example 123456789,987654321,alert,<region-123456>