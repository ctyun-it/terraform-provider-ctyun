# 导入函数资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数，不填则取值环境变量值
terraform import ctyun_function.[导入配置名称] [name][,<region_id>]
# 示例 1: 只指定函数名称，region_id 从 provider 或环境变量获取
terraform import ctyun_function.function_from_zos tf-function-from-zos
# 示例 2: 同时指定函数名称和 region_id
terraform import ctyun_function.function_from_zos tf-function-from-zos,200000002401
