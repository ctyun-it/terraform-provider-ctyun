# 导入PostgreSQL实例
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_postgresql_instance.[导入配置名称] [id],<region_id>
# 示例
terraform import ctyun_postgresql_instance.instance_example 20d1e9aa262246bf86b2915c6364715e,<region_id>