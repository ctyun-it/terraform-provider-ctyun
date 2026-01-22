# 导入PostgreSQL只读实例
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_postgresql_readonly_instance.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_postgresql_readonly_instance.readonly_example ro-inst123456,<region_id>