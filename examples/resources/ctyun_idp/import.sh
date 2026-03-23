# 导入身份提供商资源
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_idp.[导入配置名称] [id]
# 示例
terraform import ctyun_idp.idp_example 123