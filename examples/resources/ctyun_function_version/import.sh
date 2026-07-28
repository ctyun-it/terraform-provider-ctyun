# 导入函数版本资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数，不填则取值环境变量值
terraform import ctyun_function_version.[导入配置名称] [function_name],[version_id][,<region_id>]
# 示例 1: 指定函数名称和版本 ID，region_id 从 provider 或环境变量获取
terraform import ctyun_function_version.version_v1 tf-function-for-version,1
# 示例 2: 同时指定函数名称、版本 ID 和 region_id
terraform import ctyun_function_version.version_v1 tf-function-for-version,1,200000002401
