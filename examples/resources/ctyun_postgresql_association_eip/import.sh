# 导入PostgreSQL关联EIP
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_postgresql_association_eip.[导入配置名称] [instance_id],[eip_id],<region_id>
# 示例
terraform import ctyun_postgresql_association_eip.association_example 20d1e9aa262246bf86b2915c6364715e,eip-qjefnklk36,<region_id>