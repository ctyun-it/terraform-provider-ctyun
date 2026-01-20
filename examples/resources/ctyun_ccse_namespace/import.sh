# 导入CCSE命名空间
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_ccse_namespace.[导入配置名称] [namespace],[cluster_id],<region_id>
# 示例
terraform import ctyun_ccse_namespace.example test-namespace,cl-12345678,region-123