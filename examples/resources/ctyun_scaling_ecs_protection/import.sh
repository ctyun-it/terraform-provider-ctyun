# 导入伸缩组 ECS 保护配置
# [] 标记的参数为必填参数
# <> 标记的参数为可选参数,不填则取值环境变量值
# 注意：当前版本的伸缩组 ECS 保护资源配置不支持导入功能
#terraform import ctyun_scaling_ecs_protection.[导入配置名称] [id],<region_id>
# 示例
#terraform import ctyun_scaling_ecs_protection.scaling_ecs_protection_example 123456789,<region-123456>