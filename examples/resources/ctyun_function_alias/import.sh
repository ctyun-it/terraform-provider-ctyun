# 导入函数别名资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数，不填则取值环境变量值
terraform import ctyun_function_alias.[导入配置名称] [function_name],[alias_name][,<region_id>]
# 示例 1: 指定函数名称和别名，region_id 从 provider 或环境变量获取
terraform import ctyun_function_alias.alias_prod tf-function-for-alias,prod
# 示例 2: 同时指定函数名称、别名和 region_id
terraform import ctyun_function_alias.alias_prod tf-function-for-alias,prod,200000002401
