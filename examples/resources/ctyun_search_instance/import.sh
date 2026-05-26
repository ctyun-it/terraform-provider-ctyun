# 导入 OpenSearch 实例资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数，不填则取值环境变量值
terraform import ctyun_search_instance.[导入配置名称] [id][,<region_id>]
# 示例 1: 只指定实例 ID，region_id 从 provider 或环境变量获取
terraform import ctyun_search_instance.opensearch_basic opensearch-xxxxxxxxxxxxxxxxx
# 示例 2: 同时指定实例 ID 和 region_id
terraform import ctyun_search_instance.opensearch_basic opensearch-xxxxxxxxxxxxxxxxx,200000002401
