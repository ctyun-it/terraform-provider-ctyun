# 导入前缀列表
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_prefix_list.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_prefix_list.prefix_list_example 376f2f85-ff34-c4e0-4f5b-320dd427a271,<region_id>