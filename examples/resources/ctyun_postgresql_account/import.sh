# 导入PostgreSQL账号
#[] 标记的参数为必填参数
#<> 标记的参数为可选参数,不填则取值环境变量值
terraform import ctyun_postgresql_account.[导入配置名称] [instance_id],[name],<region_id>
# 示例
terraform import ctyun_postgresql_account.account_example 20d1e9aa262246bf86b2915c6364715e,myaccount,<region_id>