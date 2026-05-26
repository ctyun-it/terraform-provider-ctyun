# 导入自定义域名资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数，不填则取值环境变量值
terraform import ctyun_function_domain.[导入配置名称] [domain_name][,<region_id>]
# 示例 1: 只指定域名，region_id 从 provider 或环境变量获取
terraform import ctyun_function_domain.domain_http api.example.com
# 示例 2: 同时指定域名和 region_id
terraform import ctyun_function_domain.domain_http api.example.com,200000002401
