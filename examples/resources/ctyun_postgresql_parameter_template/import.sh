# 导入PostgreSQL参数模板
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_postgresql_parameter_template.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_postgresql_parameter_template.template_example 12345,<region_id>